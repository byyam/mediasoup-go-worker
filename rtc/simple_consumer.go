package rtc

import (
	"encoding/json"
	"github.com/byyam/mediasoup-go-worker/utils"
	"time"

	"github.com/byyam/mediasoup-go-worker/rtc/ms_rtcp"

	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"

	"github.com/pion/rtcp"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
)

type SimpleConsumer struct {
	IConsumer
	logger           utils.Logger
	rtpStream        *RtpStreamSend
	rtpStreams       []*RtpStreamSend
	maxRtcpInterval  time.Duration
	lastRtcpSentTime time.Time

	// handler
	onConsumerSendRtpPacketHandler     func(consumer IConsumer, packet *rtpparser.Packet)
	onConsumerKeyFrameRequestedHandler func(consumer IConsumer, mappedSsrc uint32)
}

type simpleConsumerParam struct {
	consumerParam
	OnConsumerSendRtpPacket       func(consumer IConsumer, packet *rtpparser.Packet)
	OnConsumerKeyFrameRequested   func(consumer IConsumer, mappedSsrc uint32)
	OnConsumerRetransmitRtpPacket func(packet *rtpparser.Packet)
}

func newSimpleConsumer(param simpleConsumerParam) (*SimpleConsumer, error) {
	var err error
	c := &SimpleConsumer{
		rtpStreams: make([]*RtpStreamSend, 0),
		logger:     utils.NewLogger("simple-consumer", param.id),
	}
	param.fillJsonStatsFunc = c.FillJsonStats
	c.IConsumer, err = newConsumer(mediasoupdata.ConsumerType_Simple, param.consumerParam)
	c.onConsumerSendRtpPacketHandler = param.OnConsumerSendRtpPacket
	c.onConsumerKeyFrameRequestedHandler = param.OnConsumerKeyFrameRequested
	if err != nil {
		return nil, err
	}
	// Set the RTCP report generation interval.
	if c.GetKind() == mediasoupdata.MediaKind_Audio {
		c.maxRtcpInterval = ms_rtcp.MaxAudioIntervalMs
	} else {
		c.maxRtcpInterval = ms_rtcp.MaxVideoIntervalMs
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
		ParamRtpStream:                 param,
		bufferSize:                     0,
		OnRtpStreamRetransmitRtpPacket: c.OnRtpStreamRetransmitRtpPacket,
	})
	c.rtpStreams = append(c.rtpStreams, c.rtpStream)
}

func (c *SimpleConsumer) SendRtpPacket(packet *rtpparser.Packet) {
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

func (c *SimpleConsumer) OnRtpStreamRetransmitRtpPacket(packet *rtpparser.Packet) {

}

func (c *SimpleConsumer) ReceiveNack(nackPacket *rtcp.TransportLayerNack) {
	c.rtpStream.ReceiveNack(nackPacket)
}

func (c *SimpleConsumer) GetRtpStreams() []*RtpStreamSend {
	return c.rtpStreams
}

func (c *SimpleConsumer) GetRtcp(rtpStream *RtpStreamSend, now time.Time) []rtcp.Packet {
	if now.Sub(c.lastRtcpSentTime) < c.maxRtcpInterval {
		return nil
	}
	c.lastRtcpSentTime = now
	var packets []rtcp.Packet
	report := rtpStream.GetRtcpSenderReport(now)
	if report == nil {
		return nil
	}
	packets = append(packets, report)
	// Build SDES chunk for this sender.
	sdes := rtpStream.GetRtcpSdesChunk()
	packets = append(packets, sdes)
	return packets
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

func (c *SimpleConsumer) NeedWorstRemoteFractionLost(worstRemoteFractionLost *uint8) {
	// If our fraction lost is worse than the given one, update it.
	if c.rtpStream.fractionLost > *worstRemoteFractionLost {
		*worstRemoteFractionLost = c.rtpStream.fractionLost
	}
}
