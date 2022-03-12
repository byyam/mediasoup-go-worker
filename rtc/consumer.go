package rtc

import (
	"encoding/json"

	"github.com/byyam/mediasoup-go-worker/common"

	"github.com/byyam/mediasoup-go-worker/utils"

	"github.com/pion/rtp"
)

type IConsumer interface {
	GetId() string
	Close()
	FillJson() json.RawMessage
	SendRtpPacket(packet *rtp.Packet)
}

type Consumer struct {
	Id         string
	ProducerId string

	logger utils.Logger
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
	id         string
	producerId string
}

func (c consumerParam) valid() bool {
	return true
}

func newConsumer(param consumerParam) (IConsumer, error) {
	if !param.valid() {
		return nil, common.ErrInvalidParam
	}

	c := &Consumer{
		Id:     param.id,
		logger: utils.NewLogger("consumer"),
	}

	return c, nil
}

func (c *Consumer) SendRtpPacket(packet *rtp.Packet) {

}
