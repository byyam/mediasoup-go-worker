package rtc

import (
	"encoding/json"
	"time"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/kr/pretty"
	"github.com/pion/rtcp"
)

type IConsumer interface {
	GetId() string
	Close()
	FillJson() json.RawMessage
	HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData)
	GetType() mediasoupdata.ConsumerType
	GetRtpParameters() mediasoupdata.RtpParameters
	SendRtpPacket(packet *rtpparser.Packet)
	ReceiveKeyFrameRequest(feedbackFormat uint8, ssrc uint32)
	GetMediaSsrcs() []uint32
	GetKind() mediasoupdata.MediaKind
	GetConsumableRtpEncodings() []mediasoupdata.RtpEncodingParameters
	ReceiveRtcpReceiverReport(report *rtcp.ReceptionReport)
	ReceiveNack(nackPacket *rtcp.TransportLayerNack)
	GetRtpStreams() []*RtpStreamSend
	GetRtcp(rtpStream *RtpStreamSend, now time.Time) []rtcp.Packet
	NeedWorstRemoteFractionLost(worstRemoteFractionLost *uint8)
}

type Consumer struct {
	Id                         string
	ProducerId                 string
	Kind                       mediasoupdata.MediaKind
	RtpHeaderExtensionIds      RtpHeaderExtensionIds
	mediaSsrcs                 []uint32
	rtxSsrcs                   []uint32
	supportedCodecPayloadTypes []uint8

	consumerType           mediasoupdata.ConsumerType
	rtpParameters          mediasoupdata.RtpParameters
	consumableRtpEncodings []mediasoupdata.RtpEncodingParameters
	fillJsonStatsFunc      func() json.RawMessage

	logger utils.Logger
}

func (c *Consumer) GetRtcp(rtpStream *RtpStreamSend, now time.Time) []rtcp.Packet {
	panic("implement me")
}

func (c *Consumer) GetRtpStreams() []*RtpStreamSend {
	panic("implement me")
}

func (c *Consumer) GetKind() mediasoupdata.MediaKind {
	return c.Kind
}

func (c *Consumer) SendRtpPacket(packet *rtpparser.Packet) {
	panic("implement me")
}

func (c *Consumer) GetType() mediasoupdata.ConsumerType {
	return c.consumerType
}

func (c *Consumer) GetRtpParameters() mediasoupdata.RtpParameters {
	return c.rtpParameters
}

func (c *Consumer) GetConsumableRtpEncodings() []mediasoupdata.RtpEncodingParameters {
	return c.consumableRtpEncodings
}

func (c *Consumer) GetId() string {
	return c.Id
}

func (c *Consumer) Close() {
	c.logger.Info("%s closed", c.Id)
}

func (c *Consumer) FillJson() json.RawMessage {
	jsonData := mediasoupdata.ConsumerDump{
		Id:                         c.Id,
		ProducerId:                 c.ProducerId,
		Kind:                       string(c.Kind),
		Type:                       string(c.consumerType),
		RtpParameters:              c.rtpParameters,
		ConsumableRtpEncodings:     nil,
		SupportedCodecPayloadTypes: nil,
		Paused:                     false,
		ProducerPaused:             false,
		Priority:                   0,
		TraceEventTypes:            "",
		RtpStreams:                 nil,
		RtpStream:                  nil,
		SimulcastConsumerDump:      nil,
	}
	data, _ := json.Marshal(&jsonData)
	c.logger.Debug("dump:%+v", jsonData)
	return data
}

func (c *Consumer) FillJsonStats() json.RawMessage {
	return c.fillJsonStatsFunc()
}

type consumerParam struct {
	id                     string
	producerId             string
	kind                   mediasoupdata.MediaKind
	rtpParameters          mediasoupdata.RtpParameters
	consumableRtpEncodings []mediasoupdata.RtpEncodingParameters
	fillJsonStatsFunc      func() json.RawMessage
}

func (c consumerParam) valid() bool {
	if len(c.consumableRtpEncodings) == 0 {
		return false
	}
	if !c.rtpParameters.Valid() {
		return false
	}
	if c.fillJsonStatsFunc == nil {
		return false
	}
	return true
}

func newConsumer(typ mediasoupdata.ConsumerType, param consumerParam) (IConsumer, error) {
	if !param.valid() {
		return nil, mserror.ErrInvalidParam
	}

	c := &Consumer{
		Id:                     param.id,
		logger:                 utils.NewLogger("consumer", param.id),
		consumerType:           typ,
		Kind:                   param.kind,
		rtpParameters:          param.rtpParameters,
		consumableRtpEncodings: param.consumableRtpEncodings,
		fillJsonStatsFunc:      param.fillJsonStatsFunc,
	}
	c.logger.Info("input param for consumer: %# v", pretty.Formatter(param))
	// init consumer with param
	if err := c.init(param); err != nil {
		return nil, err
	}
	c.logger.Info("new consumer:%# v", pretty.Formatter(c.rtpParameters))
	c.logger.Info("new consumer:%# v", pretty.Formatter(c.consumableRtpEncodings))
	c.logger.Info("new consumer:%# v", pretty.Formatter(c.mediaSsrcs))

	return c, nil
}

func (c *Consumer) init(param consumerParam) error {
	if err := c.rtpParameters.Init(); err != nil {
		return err
	}
	if err := c.RtpHeaderExtensionIds.set(param.rtpParameters.HeaderExtensions, false); err != nil {
		c.logger.Error("set RtpHeaderExtensionIds failed:%v", err)
		return err
	}
	c.logger.Info("set RtpHeaderExtensionIds:%# v", pretty.Formatter(c.RtpHeaderExtensionIds))
	// Fill supported codec payload types.
	for _, codec := range c.rtpParameters.Codecs {
		if codec.RtpCodecMimeType.IsMediaCodec() {
			c.supportedCodecPayloadTypes = append(c.supportedCodecPayloadTypes, codec.PayloadType)
		}
	}
	// Fill media SSRCs vector.
	for _, encoding := range c.rtpParameters.Encodings {
		c.mediaSsrcs = append(c.mediaSsrcs, encoding.Ssrc)
	}
	// todo: Fill RTX SSRCs vector.
	//for _, encoding := range c.rtpParameters.Encodings {
	//}
	return nil
}

func (c *Consumer) ReceiveKeyFrameRequest(feedbackFormat uint8, ssrc uint32) {
	panic("implement me")
}

func (c *Consumer) GetMediaSsrcs() []uint32 {
	return c.mediaSsrcs
}

func (c *Consumer) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		c.logger.Debug("method=%s,internal=%+v,response:%s", request.Method, request.Internal, response)
	}()

	switch request.Method {

	case mediasoupdata.MethodConsumerDump:
		response.Data = c.FillJson()

	case mediasoupdata.MethodConsumerGetStats:
		response.Data = c.FillJsonStats()
	}
}

func (c *Consumer) ReceiveRtcpReceiverReport(report *rtcp.ReceptionReport) {
	panic("implement me")
}

func (c *Consumer) ReceiveNack(nackPacket *rtcp.TransportLayerNack) {
	panic("implement me")
}

func (c *Consumer) NeedWorstRemoteFractionLost(worstRemoteFractionLost *uint8) {
	panic("implement me")
}
