package rtc

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/kr/pretty"
	"github.com/pion/rtcp"
	"github.com/rs/zerolog"

	FBS__Consumer "github.com/byyam/mediasoup-go-worker/fbs/FBS/Consumer"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Response "github.com/byyam/mediasoup-go-worker/fbs/FBS/Response"
	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
	"github.com/byyam/mediasoup-go-worker/internal/ms_rtcp"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type IConsumer interface {
	GetId() string
	Close()
	FillJson() json.RawMessage
	HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData)
	GetType() mediasoupdata.ConsumerType
	GetRtpParameters() *mediasoupdata.RtpParameters
	SendRtpPacket(packet *rtpparser.Packet)
	ReceiveKeyFrameRequest(feedbackFormat uint8, ssrc uint32)
	GetMediaSsrcs() []uint32
	GetKind() FBS__RtpParameters.MediaKind
	GetConsumableRtpEncodings() []*FBS__RtpParameters.RtpEncodingParametersT
	ReceiveRtcpReceiverReport(report *rtcp.ReceptionReport)
	ReceiveNack(nackPacket *rtcp.TransportLayerNack)
	GetRtpStreams() []*RtpStreamSend
	GetRtcp(rtpStream *RtpStreamSend, now time.Time) []rtcp.Packet
	NeedWorstRemoteFractionLost(worstRemoteFractionLost *uint8)

	// maxRtcpInterval
	GetMaxRtcpInterval() time.Duration
	// lastRtcpSentTime
	GetLastRtcpSentTime() time.Time
	SetLastRtcpSentTime(v time.Time)
}

type Consumer struct {
	Id                         string
	ProducerId                 string
	Kind                       FBS__RtpParameters.MediaKind
	RtpHeaderExtensionIds      RtpHeaderExtensionIds
	mediaSsrcs                 []uint32
	rtxSsrcs                   []uint32
	supportedCodecPayloadTypes []uint8

	maxRtcpInterval  time.Duration
	lastRtcpSentTime time.Time

	consumerType           mediasoupdata.ConsumerType
	rtpParameters          *mediasoupdata.RtpParameters
	consumableRtpEncodings []*FBS__RtpParameters.RtpEncodingParametersT
	fillJsonStatsFunc      func() *FBS__Consumer.GetStatsResponseT

	logger zerolog.Logger
}

// getter & setter
func (c *Consumer) GetMaxRtcpInterval() time.Duration {
	return c.maxRtcpInterval
}

func (c *Consumer) GetLastRtcpSentTime() time.Time {
	return c.lastRtcpSentTime
}

func (c *Consumer) SetLastRtcpSentTime(v time.Time) {
	c.lastRtcpSentTime = v
}

func (c *Consumer) GetRtcp(rtpStream *RtpStreamSend, now time.Time) []rtcp.Packet {
	panic("implement me")
}

func (c *Consumer) GetRtpStreams() []*RtpStreamSend {
	panic("implement me")
}

func (c *Consumer) GetKind() FBS__RtpParameters.MediaKind {
	return c.Kind
}

func (c *Consumer) SendRtpPacket(packet *rtpparser.Packet) {
	panic("implement me")
}

func (c *Consumer) GetType() mediasoupdata.ConsumerType {
	return c.consumerType
}

func (c *Consumer) GetRtpParameters() *mediasoupdata.RtpParameters {
	return c.rtpParameters
}

func (c *Consumer) GetConsumableRtpEncodings() []*FBS__RtpParameters.RtpEncodingParametersT {
	return c.consumableRtpEncodings
}

func (c *Consumer) GetId() string {
	return c.Id
}

func (c *Consumer) Close() {
	c.logger.Info().Msg("closed")
}

func (c *Consumer) FillJson() json.RawMessage {
	jsonData := mediasoupdata.ConsumerDump{
		Id:         c.Id,
		ProducerId: c.ProducerId,
		Kind:       string(c.Kind),
		Type:       string(c.consumerType),
		// RtpParameters:              c.rtpParameters,
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
	c.logger.Debug().Msgf("dump:%+v", jsonData)
	return data
}

func (c *Consumer) FillJsonStats() *FBS__Consumer.GetStatsResponseT {
	return c.fillJsonStatsFunc()
}

type consumerParam struct {
	id                     string
	producerId             string
	kind                   FBS__RtpParameters.MediaKind
	rtpParameters          *mediasoupdata.RtpParameters
	consumableRtpEncodings []*FBS__RtpParameters.RtpEncodingParametersT
	fillJsonStatsFunc      func() *FBS__Consumer.GetStatsResponseT
}

func (c consumerParam) valid() error {
	if len(c.consumableRtpEncodings) == 0 {
		return errors.New("consumableRtpEncodings empty")
	}
	if !c.rtpParameters.Valid() {
		return errors.New("rtpParameters invalid")
	}
	if c.fillJsonStatsFunc == nil {
		return errors.New("fillJsonStatsFunc nil")
	}
	return nil
}

func newConsumer(typ mediasoupdata.ConsumerType, param *consumerParam) (IConsumer, error) {
	if err := param.valid(); err != nil {
		return nil, err
	}

	c := &Consumer{
		Id:                     param.id,
		logger:                 zerowrapper.NewScope("consumer", param.id),
		consumerType:           typ,
		Kind:                   param.kind,
		rtpParameters:          param.rtpParameters,
		consumableRtpEncodings: param.consumableRtpEncodings,
		fillJsonStatsFunc:      param.fillJsonStatsFunc,
	}
	// init consumer with param
	if err := c.init(param); err != nil {
		return nil, err
	}
	c.logger.Info().Any("RtpParameters", mediasoupdata.JsonFormat(c.rtpParameters)).
		Any("consumableRtpEncodings", mediasoupdata.JsonFormat(c.consumableRtpEncodings)).
		Uints32("mediaSsrcs", c.mediaSsrcs).
		Msgf("init rtp param for consumer")

	// Set the RTCP report generation interval.
	if c.GetKind() == FBS__RtpParameters.MediaKindAUDIO {
		c.maxRtcpInterval = ms_rtcp.MaxAudioIntervalMs
	} else {
		c.maxRtcpInterval = ms_rtcp.MaxVideoIntervalMs
	}

	return c, nil
}

func (c *Consumer) init(param *consumerParam) error {
	if err := c.rtpParameters.Init(); err != nil {
		return err
	}
	if err := c.RtpHeaderExtensionIds.set(param.rtpParameters.HeaderExtensions, false); err != nil {
		c.logger.Error().Err(err).Msg("set RtpHeaderExtensionIds failed")
		return err
	}
	c.logger.Info().Msgf("set RtpHeaderExtensionIds:%# v", pretty.Formatter(c.RtpHeaderExtensionIds))
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
		c.logger.Debug().Msgf("method=%s,internal=%+v,response:%s", request.Method, request.Internal, response)
	}()

	switch request.MethodType {

	case FBS__Request.MethodCONSUMER_DUMP:
		response.Data = c.FillJson()

	case FBS__Request.MethodCONSUMER_GET_STATS:
		dataDump := c.FillJsonStats()
		// set rsp
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyDataConsumer_GetStatsResponse,
			Value: dataDump,
		}
		response.RspBody = rspBody
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
