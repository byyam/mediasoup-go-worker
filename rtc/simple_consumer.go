package rtc

import (
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/pion/rtp"
)

type SimpleConsumer struct {
	IConsumer
	logger utils.Logger
}

type simpleConsumerParam struct {
	consumerParam
}

func newSimpleConsumer(param simpleConsumerParam) (*SimpleConsumer, error) {
	var err error
	c := &SimpleConsumer{
		logger: utils.NewLogger("simple-consumer"),
	}
	c.IConsumer, err = newConsumer(param.consumerParam)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *SimpleConsumer) SendRtpPacket(packet *rtp.Packet) {
	c.logger.Debug("SendRtpPacket")
}
