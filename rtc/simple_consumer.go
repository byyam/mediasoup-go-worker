package rtc

import (
	"encoding/json"

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
	onConsumerSendRtpPacketHandler     func(consumer IConsumer, packet *rtp.Packet)
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
	param.fillJsonStatsFunc = c.FillJsonStats
	c.IConsumer, err = newConsumer(mediasoupdata.ConsumerType_Simple, param.consumerParam)
	c.onConsumerSendRtpPacketHandler = param.OnConsumerSendRtpPacket
	c.onConsumerKeyFrameRequestedHandler = param.OnConsumerKeyFrameRequested
	if err != nil {
		return nil, err
	}
	// Create RtpStreamSend instance for sending a single stream to the remote.
	c.CreateRtpStream()
	return c, nil
}

func (c *SimpleConsumer) CreateRtpStream() {
	rtpParameters := c.IConsumer.GetRtpParameters()
	encoding := rtpParameters.Encodings[0]
	mediaCodec := rtpParameters.GetCodecForEncoding(encoding)
	param := &ParamRtpStream{
		EncodingIdx:    0,
		Ssrc:           encoding.Ssrc,
		PayloadType:    mediaCodec.PayloadType,
		MimeType:       mediaCodec.RtpCodecMimeType,
		ClockRate:      mediaCodec.ClockRate,
		Rid:            "",
		Cname:          rtpParameters.Rtcp.Cname,
		RtxSsrc:        0,
		RtxPayloadType: 0,
		UseNack:        false,
		UsePli:         false,
		UseFir:         false,
		UseInBandFec:   false,
		UseDtx:         false,
		SpatialLayers:  0,
		TemporalLayers: 0,
	}
	c.rtpStream = newRtpStreamSend(&ParamRtpStreamSend{
		ParamRtpStream: param,
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
	if c.rtpStream.ReceivePacket(packet) { // todo
		c.onConsumerSendRtpPacketHandler(c.IConsumer, packet)
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

func (c *SimpleConsumer) FillJsonStats() json.RawMessage {
	var jsonData []mediasoupdata.ConsumerStat
	if c.rtpStream != nil {
		var stat mediasoupdata.ConsumerStat
		c.rtpStream.FillJsonStats(&stat)
		jsonData = append(jsonData, stat)
	}
	data, _ := json.Marshal(&jsonData)
	c.logger.Debug("getStats:%+v", jsonData)
	return data
}
