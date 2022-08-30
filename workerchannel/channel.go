package workerchannel

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/pkg/netparser"
)

const (
	// netstring length for a 4194304 bytes payload.
	NS_MESSAGE_MAX_LEN = 4194308
	NS_PAYLOAD_MAX_LEN = 4194304
)

type Channel struct {
	logger     utils.Logger
	netParser  netparser.INetParser
	jsonFormat bool

	// handle request message
	OnRequestHandler atomic.Value // func(request RequestData) ResponseData
}

func NewChannel(netParser netparser.INetParser, id string, jsonFormat bool) *Channel {
	c := &Channel{
		netParser:  netParser,
		jsonFormat: jsonFormat,
		logger:     utils.NewLogger(fmt.Sprintf("channel-json[%v]", jsonFormat), id),
	}
	c.logger.Info("channel start")
	go c.runReadLoop()
	return c
}

func (c *Channel) runReadLoop() {
	defer c.Close()
	payload := make([]byte, NS_PAYLOAD_MAX_LEN)
	for {
		n, err := c.netParser.ReadBuffer(payload)
		if err != nil {
			c.logger.Error("Channel error:%v", err)
			break
		}
		c.logger.Debug("payload:%s", string(payload[:n]))
		if !c.jsonFormat {
			c.processPayload(payload[:n])
		} else {
			c.processPayloadJsonFormat(payload[:n])
		}
	}
}

func (c *Channel) processPayloadJsonFormat(nsPayload []byte) {
	switch nsPayload[0] {
	case '{':
		if err := c.processJsonMessage(nsPayload); err != nil {
			c.logger.Error("process message failed:%s", err)
		}
	case 'X':
		c.logger.Debug("%s\n", nsPayload[1:])
	default:
		c.logger.Error("unexpected data:[%s]", nsPayload)
	}
}

func (c *Channel) processPayload(nsPayload []byte) {
	// https://github.com/versatica/mediasoup/commit/ed15a863a5ed095f58a16d972c8e25bf24f17933
	// const request = `${id}:${method}:${handlerId}:${JSON.stringify(data)}`;
	messages := strings.Split(string(nsPayload), ":")
	if len(messages) != 4 {
		c.logger.Error("messages length invalid:[%s]", nsPayload)
		return
	}

}

func (c *Channel) OnRequest(fn func(request RequestData) ResponseData) {
	c.OnRequestHandler.Store(fn)
}

func (c *Channel) processMessage(message []string) error {
	//reqData := channelData{
	//	Id:       0, // todo
	//	Method:   message[1],
	//	Internal: nil,
	//	Data:     nil,
	//}
	c.logger("message:%v", message)
	return nil
}

func (c *Channel) processJsonMessage(nsPayload []byte) error {
	var reqData channelData
	if err := json.Unmarshal(nsPayload, &reqData); err != nil {
		return err
	}
	var internal InternalData
	_ = internal.Unmarshal(reqData.Internal)
	c.logger.Info("request Id=%d, Method=%s", reqData.Id, reqData.Method)

	rspData, _ := c.handleMessage(&reqData, &internal)

	jsonByte, _ := json.Marshal(&rspData)

	if len(jsonByte) > NS_MESSAGE_MAX_LEN {
		return errors.New("channel response too big")
	}
	c.logger.Trace("WriteBuffer:[%s],rspData:%+v", string(jsonByte), rspData)
	if err := c.netParser.WriteBuffer(jsonByte); err != nil {
		return err
	}
	c.logger.Info("response Id=%d,err=[%v]", rspData.Id, rspData.Error)
	return nil
}

func (c *Channel) handleMessage(reqData *channelData, internal *InternalData) (*channelData, error) {
	var ret ResponseData
	rspData := new(channelData)
	rspData.Id = reqData.Id
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
		c.logger.Error("response error:%v, Id=%d, Method=%s", ret.Err, reqData.Id, reqData.Method)
		rspData.Error = ret.Err.Error()
		rspData.Reason = ret.Err.Error()
	} else {
		rspData.Accepted = true
	}
	rspData.Data = ret.Data

	return rspData, nil
}

func (c *Channel) Event(targetId int, event string) {
	msg := channelData{
		TargetId: strconv.Itoa(targetId),
		Event:    event,
	}
	jsonByte, _ := json.Marshal(&msg)
	err := c.netParser.WriteBuffer(jsonByte)
	c.logger.Info("send Event msg:%+v,err=%v", msg, err)
}

func (c *Channel) Close() {
	c.logger.Info("closed")
}
