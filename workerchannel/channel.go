package workerchannel

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/fbs/FBS/Message"
	"github.com/byyam/mediasoup-go-worker/fbs/FBS/Notification"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/tidwall/gjson"

	"github.com/byyam/mediasoup-go-worker/pkg/netparser"
)

const (
	// netstring length for a 4194304 bytes payload.
	NS_MESSAGE_MAX_LEN = 4194308
	NS_PAYLOAD_MAX_LEN = 4194304
)

const (
	UNDEFINED = "undefined"
)

const (
	NativeJsonFormat = iota + 1
	NativeFormat
	FlatBufferFormat
)

type Channel struct {
	logger       zerolog.Logger
	netParser    netparser.INetParser
	bufferFormat int

	// handle request message
	OnRequestHandler atomic.Value // func(request RequestData) ResponseData
}

func NewChannel(netParser netparser.INetParser, id string, bufferFormat int) *Channel {
	c := &Channel{
		netParser:    netParser,
		bufferFormat: bufferFormat,
		logger:       zerowrapper.NewScope(fmt.Sprintf("channel-json[%v]", bufferFormat), id),
	}
	c.logger.Info().Msgf("channel start, bufferFormat=%d", c.bufferFormat)
	go c.runReadLoop()
	return c
}

func (c *Channel) runReadLoop() {
	defer c.Close()
	payload := make([]byte, NS_PAYLOAD_MAX_LEN)
	for {
		n, err := c.netParser.ReadBuffer(payload)
		if err != nil {
			c.logger.Error().Err(err).Msg("[runReadLoop]Channel read buffer error")
			break
		}
		c.logger.Debug().Str("payload", string(payload[:n])).Send()
		switch c.bufferFormat {
		case NativeJsonFormat:
			c.processPayloadJsonFormat(payload[:n])
		case NativeFormat:
			c.processPayloadNative(payload[:n])
		case FlatBufferFormat:
			c.processPayloadFB(payload[:n])
		default:
			c.logger.Error().Int("[runReadLoop]bufferFormat", c.bufferFormat).Send()
		}
	}
}

func (c *Channel) processPayloadJsonFormat(nsPayload []byte) {
	switch nsPayload[0] {
	case '{':
		if err := c.processJsonMessage(nsPayload); err != nil {
			c.logger.Error().Err(err).Msg("[processPayloadJsonFormat] processJsonMessage failed")
		}
	case 'X':
		c.logger.Debug().Str("payload", string(nsPayload[1:])).Send()
	default:
		c.logger.Error().Str("payload", string(nsPayload)).Msg("unexpected data")
	}
}

func (c *Channel) processPayloadNative(nsPayload []byte) {
	// https://github.com/versatica/mediasoup/commit/ed15a863a5ed095f58a16d972c8e25bf24f17933
	// const request = `${id}:${method}:${handlerId}:${JSON.stringify(data)}`;
	messages := strings.SplitN(string(nsPayload), ":", 4)
	c.logger.Info().Strs("messages", messages).Send()
	if len(messages) != 4 {
		c.logger.Error().Strs("messages", messages).Msg("messages length invalid")
		return
	}
	if err := c.processNativeMessage(messages); err != nil {
		c.logger.Error().Err(err).Msg("[processPayloadNative] processNativeMessage failed")
	}
}

func (c *Channel) processPayloadFB(nsPayload []byte) {
	message := Message.GetRootAsMessage(nsPayload, 0)
	messageOffset := message.UnPack()
	bodyOffset := messageOffset.Data.Type.UnPack(message.Table())
	c.logger.Info().Msgf("[processPayloadFB]msg offset:%+v", messageOffset)
	switch bodyOffset.Type {
	case Message.BodyRequest:
		requestOffset := bodyOffset.Value.(*FBS__Request.RequestT)
		c.logger.Info().Msgf("[processPayloadFB]request method:%+v", requestOffset)
		if err := c.processFBMessage(requestOffset); err != nil {
			c.logger.Error().Err(err).Msg("[processPayloadFB] processFBMessage failed")
		}
	default:
		c.logger.Error().Int("DataType", int(message.DataType())).Msg("[processPayloadFB]unexpected data type")
	}
}

func (c *Channel) OnRequest(fn func(request RequestData) ResponseData) {
	c.OnRequestHandler.Store(fn)
}

func (c *Channel) setHandlerId(method, handlerId, data string, internal *InternalData) error {
	if handlerId == UNDEFINED && data != UNDEFINED {
		if err := internal.Unmarshal(json.RawMessage(data)); err != nil {
			return err
		}
		return nil
	}
	// get prefix of method
	methodFields := strings.SplitN(method, ".", 2)
	if len(methodFields) != 2 {
		return errors.New("method is not formatted")
	}
	// set handlerId
	switch methodFields[0] {
	case mediasoupdata.MethodPrefixWorker, mediasoupdata.MethodPrefixRouter:
		internal.RouterId = handlerId
	case mediasoupdata.MethodPrefixTransport:
		internal.TransportId = handlerId
	case mediasoupdata.MethodPrefixProducer:
		internal.ProducerId = handlerId
	case mediasoupdata.MethodPrefixConsumer:
		internal.ConsumerId = handlerId
	case mediasoupdata.MethodPrefixRtpObserver:
		internal.RtpObserverId = handlerId
	default:
		return errors.New("unknown method prefix")
	}
	// set objectId
	switch methodFields[0] {
	case mediasoupdata.MethodPrefixRouter:
		// set ids if fields exist
		value := gjson.Get(data, "transportId")
		internal.TransportId = value.String()

		value = gjson.Get(data, "rtpObserverId")
		internal.RtpObserverId = value.String()
	case mediasoupdata.MethodPrefixTransport, mediasoupdata.MethodPrefixRtpObserver: // include producer and consumer
		value := gjson.Get(data, "producerId")
		internal.ProducerId = value.String()

		value = gjson.Get(data, "dataProducerId")
		internal.DataProducerId = value.String()
	}
	if method == mediasoupdata.MethodTransportConsume {
		value := gjson.Get(data, "consumerId")
		internal.ConsumerId = value.String()
	}
	return nil
}

func (c *Channel) processFBMessage(req *FBS__Request.RequestT) error {
	reqData := &channelData{
		Id:       int64(req.Id),
		Method:   FBSRequestMethod[req.Method],
		Internal: nil,
		Data:     nil,
	}
	internal := &InternalData{
		RouterId:       "",
		TransportId:    "",
		ProducerId:     "",
		ConsumerId:     "",
		DataProducerId: "",
		DataConsumerId: "",
		RtpObserverId:  "",
	}
	// handle
	rspData, _ := c.handleMessage(reqData, internal)
	c.logger.Info().Int64("id", rspData.Id).Str("method", rspData.Method).Msg("rspData")

	return nil
}

func (c *Channel) processNativeMessage(messages []string) error {
	idStr := messages[0]
	method := messages[1]
	handlerId := messages[2]
	data := messages[3]
	// decode
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	var internal InternalData

	if err := c.setHandlerId(method, handlerId, data, &internal); err != nil {
		c.logger.Error().Err(err).Msg("set handler Id failed")
		return err
	}

	reqData := channelData{
		Id:       int64(id),
		Method:   method,
		Internal: nil,
		Data:     json.RawMessage(data),
	}

	// handle
	rspData, _ := c.handleMessage(&reqData, &internal)
	c.logger.Info().Int64("id", rspData.Id).Str("method", rspData.Method).Msg("rspData")

	// encode
	if err := c.returnMessage(rspData); err != nil {
		c.logger.Error().Err(err).Msg("return message failed")
		return err
	}
	return nil
}

func (c *Channel) processJsonMessage(nsPayload []byte) error {
	// decode
	var reqData channelData
	if err := json.Unmarshal(nsPayload, &reqData); err != nil {
		return err
	}
	var internal InternalData
	_ = internal.Unmarshal(reqData.Internal)
	c.logger.Info().Int64("id", reqData.Id).Str("method", reqData.Method).Msg("reqData")

	// handle
	rspData, _ := c.handleMessage(&reqData, &internal)

	// encode
	if err := c.returnMessage(rspData); err != nil {
		c.logger.Error().Err(err).Msg("return message failed")
	}
	return nil
}

func (c *Channel) returnMessage(rspData *channelData) error {
	jsonByte, _ := json.Marshal(&rspData)

	if len(jsonByte) > NS_MESSAGE_MAX_LEN {
		return errors.New("channel response too big")
	}
	c.logger.Trace().Str("WriteBuffer", string(jsonByte)).Int64("id", rspData.Id).Str("method", rspData.Method).Send()
	if err := c.netParser.WriteBuffer(jsonByte); err != nil {
		return err
	}
	c.logger.Info().Str("error", rspData.Error).Int64("id", rspData.Id).Msg("response")
	return nil
}

func (c *Channel) handleMessage(reqData *channelData, internal *InternalData) (*channelData, error) {
	var ret ResponseData
	rspData := new(channelData)
	rspData.Id = reqData.Id
	rspData.Method = reqData.Method
	if handler, ok := c.OnRequestHandler.Load().(func(request RequestData) ResponseData); ok && handler != nil {
		ret = handler(RequestData{
			Method:   reqData.Method,
			Internal: *internal,
			Data:     reqData.Data,
		})
	} else {
		rspData.Error = "OnRequestHandler not register"
	}

	if ret.Err != nil {
		c.logger.Error().Err(ret.Err).Int64("id", reqData.Id).Str("method", reqData.Method).Msg("response error")
		rspData.Error = ret.Err.Error()
		rspData.Reason = ret.Err.Error()
	} else {
		rspData.Accepted = true
	}
	rspData.Data = ret.Data

	return rspData, nil
}

func (c *Channel) Event(targetId string, event Notification.Event) {
	switch c.bufferFormat {
	case NativeJsonFormat, NativeFormat:
		c.eventJson(targetId, event)
	case FlatBufferFormat:
		c.eventFB(targetId, event)
	default:
		c.logger.Error().Int("[Event]bufferFormat", c.bufferFormat).Send()
	}
}

func (c *Channel) eventJson(targetId string, event Notification.Event) {
	msg := channelData{
		TargetId: targetId,
		Event:    EventMap[event],
	}
	jsonByte, _ := json.Marshal(&msg)
	err := c.netParser.WriteBuffer(jsonByte)
	c.logger.Info().Err(err).Str("targetId", msg.TargetId).Int("event", int(event)).Msg("[eventJson]send Event msg")
}

func (c *Channel) eventFB(targetId string, event Notification.Event) {
	b := flatbuffers.NewBuilder(0)
	// targetIdOffset
	targetIdOffset := flatbuffers.UOffsetT(0)
	if targetId != "" {
		targetIdOffset = b.CreateString(targetId)
	}
	// notification offset
	Notification.NotificationStart(b)
	Notification.NotificationAddHandlerId(b, targetIdOffset)
	Notification.NotificationAddEvent(b, Notification.EventWORKER_RUNNING)
	Notification.NotificationAddBodyType(b, Notification.BodyNONE)
	Notification.NotificationAddBody(b, flatbuffers.UOffsetT(0))
	notificationOffset := Notification.NotificationEnd(b)
	// msg offset
	Message.MessageStart(b)
	Message.MessageAddDataType(b, Message.BodyNotification)
	Message.MessageAddData(b, notificationOffset)
	messageOffset := Message.MessageEnd(b)
	// send msg
	b.Finish(messageOffset)
	err := c.netParser.WriteBuffer(b.FinishedBytes())
	c.logger.Info().Err(err).Str("targetId", targetId).Int("event", int(event)).Msg("[eventFB]send Event msg")
}

func (c *Channel) Close() {
	c.logger.Info().Msg("closed")
}
