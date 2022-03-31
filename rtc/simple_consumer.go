package rtc

import (
	"sync/atomic"

	"github.com/pion/rtcp"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/pion/rtp"
)

type SimpleConsumer struct {
	IConsumer
	logger    utils.Logger
	rtpStream *RtpStreamSend

	// handler
	onConsumerSendRtpPacketHandler     atomic.Value
	onConsumerKeyFrameRequestedHandler func(consumer IConsumer, mappedSsrc uint32)
}

type simpleConsumerParam struct {
	consumerParam
	OnConsumerSendRtpPacket     func(consumer IConsumer, packet *rtp.Packet)
	OnConsumerKeyFrameRequested func(consumer IConsumer, mappedSsrc uint32)
}

func newSimpleConsumer(param simpleConsumerParam) (*SimpleConsumer, error) {
	var err error
	c := &SimpleConsumer{
		logger: utils.NewLogger("simple-consumer", param.id),
	}
	c.IConsumer, err = newConsumer(mediasoupdata.ConsumerType_Simple, param.consumerParam)
	c.onConsumerSendRtpPacketHandler.Store(param.OnConsumerSendRtpPacket)
	c.onConsumerKeyFrameRequestedHandler = param.OnConsumerKeyFrameRequested
	if err != nil {
		return nil, err
	}
	// Create RtpStreamSend instance for sending a single stream to the remote.
	c.CreateRtpStream()
	return c, nil
}

func (c *SimpleConsumer) CreateRtpStream() {
	c.rtpStream = newRtpStreamSend(&ParamRtpStreamSend{
		ParamRtpStream: nil,
		bufferSize:     0,
	})
}

func (c *SimpleConsumer) SendRtpPacket(packet *rtp.Packet) {
	if c.GetKind() == mediasoupdata.MediaKind_Video {
		monitor.RtpSendCount(monitor.TraceVideo)
	} else if c.GetKind() == mediasoupdata.MediaKind_Audio {
		monitor.RtpSendCount(monitor.TraceAudio)
	}
	packet.SSRC = c.GetRtpParameters().Encodings[0].Ssrc
	packet.PayloadType = c.GetRtpParameters().Codecs[0].PayloadType
	if handler, ok := c.onConsumerSendRtpPacketHandler.Load().(func(consumer IConsumer, packet *rtp.Packet)); ok && handler != nil {
		handler(c.IConsumer, packet)
	}
	monitor.MediasoupCount(monitor.SimpleConsumer, monitor.EventSendRtp)
	c.logger.Trace("SendRtpPacket:%+v", packet.Header)
}

func (c *SimpleConsumer) Close() {
	c.logger.Info("%s closed", c.GetId())
}

func (c *SimpleConsumer) ReceiveKeyFrameRequest(feedbackFormat uint8, ssrc uint32) {
	// todo: trace emit
	c.RequestKeyFrame()
}

func (c *SimpleConsumer) RequestKeyFrame() {
	if c.GetKind() != mediasoupdata.MediaKind_Video {
		return
	}
	mappedSsrc := c.GetConsumableRtpEncodings()[0].Ssrc
	c.onConsumerKeyFrameRequestedHandler(c.IConsumer, mappedSsrc)
}

func (c *SimpleConsumer) ReceiveRtcpReceiverReport(report *rtcp.ReceptionReport) {
	c.rtpStream.ReceiveRtcpReceiverReport(report)
}
