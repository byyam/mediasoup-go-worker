package rtc

import (
	"encoding/json"

	"github.com/pion/rtp"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"

	"github.com/byyam/mediasoup-go-worker/common"

	"github.com/byyam/mediasoup-go-worker/utils"
)

type IConsumer interface {
	GetId() string
	Close()
	FillJson() json.RawMessage
	GetType() mediasoupdata.ConsumerType
	GetRtpParameters() mediasoupdata.RtpParameters
	SendRtpPacket(packet *rtp.Packet)
}

type Consumer struct {
	Id            string
	ProducerId    string
	consumerType  mediasoupdata.ConsumerType
	rtpParameters mediasoupdata.RtpParameters

	logger utils.Logger
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

func (c *Consumer) GetId() string {
	return c.Id
}

func (c *Consumer) Close() {
	//TODO implement me
	panic("implement me")
}

func (c *Consumer) FillJson() json.RawMessage {
	//TODO implement me
	panic("implement me")
}

type consumerParam struct {
	id            string
	producerId    string
	rtpParameters mediasoupdata.RtpParameters
}

func (c consumerParam) valid() bool {
	return true
}

func newConsumer(typ mediasoupdata.ConsumerType, param consumerParam) (IConsumer, error) {
	if !param.valid() {
		return nil, common.ErrInvalidParam
	}

	c := &Consumer{
		Id:            param.id,
		logger:        utils.NewLogger("consumer"),
		consumerType:  typ,
		rtpParameters: param.rtpParameters,
	}

	return c, nil
}
