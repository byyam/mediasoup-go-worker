package rtc

import (
	"encoding/json"

	"github.com/pion/rtp"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"

	"github.com/byyam/mediasoup-go-worker/common"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
)

type IConsumer interface {
	GetId() string
	Close()
	FillJson() json.RawMessage
	GetType() mediasoupdata.ConsumerType
	GetRtpParameters() mediasoupdata.RtpParameters
	SendRtpPacket(packet *rtp.Packet)
	ReceiveKeyFrameRequest(feedbackFormat uint8, ssrc uint32)
	GetMediaSsrcs() []uint32
	GetKind() mediasoupdata.MediaKind
	GetConsumableRtpEncodings() []mediasoupdata.RtpEncodingParameters
}

type Consumer struct {
	Id         string
	ProducerId string
	Kind       mediasoupdata.MediaKind
	mediaSsrcs []uint32

	consumerType           mediasoupdata.ConsumerType
	rtpParameters          mediasoupdata.RtpParameters
	consumableRtpEncodings []mediasoupdata.RtpEncodingParameters

	logger utils.Logger
}

func (c *Consumer) GetKind() mediasoupdata.MediaKind {
	return c.Kind
}

func (c *Consumer) SendRtpPacket(packet *rtp.Packet) {
	//TODO implement me
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
	//TODO implement me
	panic("implement me")
}

type consumerParam struct {
	id                     string
	producerId             string
	kind                   mediasoupdata.MediaKind
	rtpParameters          mediasoupdata.RtpParameters
	consumableRtpEncodings []mediasoupdata.RtpEncodingParameters
}

func (c consumerParam) valid() bool {
	if len(c.consumableRtpEncodings) == 0 {
		return false
	}
	if !c.rtpParameters.Valid() {
		return false
	}
	return true
}

func newConsumer(typ mediasoupdata.ConsumerType, param consumerParam) (IConsumer, error) {
	if !param.valid() {
		return nil, common.ErrInvalidParam
	}

	c := &Consumer{
		Id:                     param.id,
		logger:                 utils.NewLogger("consumer"),
		consumerType:           typ,
		Kind:                   param.kind,
		rtpParameters:          param.rtpParameters,
		consumableRtpEncodings: param.consumableRtpEncodings,
	}
	for _, encoding := range c.rtpParameters.Encodings {
		c.mediaSsrcs = append(c.mediaSsrcs, encoding.Ssrc)
	}

	return c, nil
}

func (c *Consumer) ReceiveKeyFrameRequest(feedbackFormat uint8, ssrc uint32) {
	//TODO implement me
	panic("implement me")
}

func (c *Consumer) GetMediaSsrcs() []uint32 {
	return c.mediaSsrcs
}
