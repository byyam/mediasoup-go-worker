package rtc

import (
	"sync/atomic"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/pion/rtp"
)

type SimpleConsumer struct {
	IConsumer
	logger utils.Logger

	// handler
	onConsumerSendRtpPacketHandler atomic.Value
}

type simpleConsumerParam struct {
	consumerParam
	OnConsumerSendRtpPacket func(consumer IConsumer, packet *rtp.Packet)
}

func newSimpleConsumer(param simpleConsumerParam) (*SimpleConsumer, error) {
	var err error
	c := &SimpleConsumer{
		logger: utils.NewLogger("simple-consumer"),
	}
	c.IConsumer, err = newConsumer(mediasoupdata.ConsumerType_Simple, param.consumerParam)
	c.onConsumerSendRtpPacketHandler.Store(param.OnConsumerSendRtpPacket)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *SimpleConsumer) SendRtpPacket(packet *rtp.Packet) {
	packet.SSRC = c.GetRtpParameters().Encodings[0].Ssrc
	packet.PayloadType = c.GetRtpParameters().Codecs[0].PayloadType
	if handler, ok := c.onConsumerSendRtpPacketHandler.Load().(func(consumer IConsumer, packet *rtp.Packet)); ok && handler != nil {
		handler(c.IConsumer, packet)
	}
	c.logger.Trace("SendRtpPacket:%+v", packet.Header)
}

func (c *SimpleConsumer) Close() {
	c.logger.Info("%s closed", c.GetId())
}
