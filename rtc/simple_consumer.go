package rtc

import (
	"time"

	"github.com/pion/rtcp"
	"github.com/rs/zerolog"
	"go.uber.org/zap"

	FBS__Consumer "github.com/byyam/mediasoup-go-worker/fbs/FBS/Consumer"
	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
	FBS__RtpStream "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpStream"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/zaplog"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type SimpleConsumer struct {
	IConsumer
	logger     zerolog.Logger
	rtpStream  *RtpStreamSend
	rtpStreams []*RtpStreamSend

	// handler
	onConsumerSendRtpPacketHandler     func(consumer IConsumer, packet *rtpparser.Packet)
	onConsumerKeyFrameRequestedHandler func(consumer IConsumer, mappedSsrc uint32)
}

type simpleConsumerParam struct {
	*consumerParam
	OnConsumerSendRtpPacket       func(consumer IConsumer, packet *rtpparser.Packet)
	OnConsumerKeyFrameRequested   func(consumer IConsumer, mappedSsrc uint32)
	OnConsumerRetransmitRtpPacket func(packet *rtpparser.Packet)
}

func newSimpleConsumer(param simpleConsumerParam) (*SimpleConsumer, error) {
	var err error
	c := &SimpleConsumer{
		rtpStreams: make([]*RtpStreamSend, 0),
		logger:     zerowrapper.NewScope("simple-consumer", param.id),
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

	workerchannel.RegisterHandler(param.id, c.HandleRequest)
	return c, nil
}

func (c *SimpleConsumer) CreateRtpStream() {
	rtpParameters := c.IConsumer.GetRtpParameters()
	encoding := rtpParameters.Encodings[0]
	mediaCodec := rtpParameters.GetCodecForEncoding(encoding)
	param := &ParamRtpStream{
		EncodingIdx:    0,
		Ssrc:           *encoding.Ssrc,
		PayloadType:    mediaCodec.PayloadType,
		MimeType:       mediaCodec.RtpCodecMimeType,
		ClockRate:      int(mediaCodec.ClockRate),
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
		// todo Kind
	}
	c.rtpStream = newRtpStreamSend(&ParamRtpStreamSend{
		ParamRtpStream:                 param,
		bufferSize:                     0,
		OnRtpStreamRetransmitRtpPacket: c.OnRtpStreamRetransmitRtpPacket,
	})
	c.rtpStreams = append(c.rtpStreams, c.rtpStream)
}

func (c *SimpleConsumer) SendRtpPacket(packet *rtpparser.Packet) {
	if c.GetKind() == FBS__RtpParameters.MediaKindVIDEO {
		monitor.RtpSendCount(monitor.TraceVideo)
	} else if c.GetKind() == FBS__RtpParameters.MediaKindAUDIO {
		monitor.RtpSendCount(monitor.TraceAudio)
	}
	packet.SSRC = *c.GetRtpParameters().Encodings[0].Ssrc
	packet.PayloadType = c.GetRtpParameters().Codecs[0].PayloadType
	if c.rtpStream.ReceivePacket(packet) { // todo
		c.onConsumerSendRtpPacketHandler(c.IConsumer, packet)
	}
	monitor.MediasoupCount(monitor.SimpleConsumer, monitor.EventSendRtp)
	zaplog.NewLogger().Info("SimpleConsumer: SendRtpPacket", zap.String("kind", string(c.IConsumer.GetKind())), zap.String("packet", packet.String()))
}

func (c *SimpleConsumer) Close() {
	c.logger.Info().Msg("closed")
	workerchannel.UnregisterHandler(c.GetId())
}

func (c *SimpleConsumer) ReceiveKeyFrameRequest(feedbackFormat uint8, ssrc uint32) {
	// todo: trace emit
	c.RequestKeyFrame()
}

func (c *SimpleConsumer) RequestKeyFrame() {
	if c.GetKind() != FBS__RtpParameters.MediaKindVIDEO {
		return
	}
	mappedSsrc := c.GetConsumableRtpEncodings()[0].Ssrc
	c.onConsumerKeyFrameRequestedHandler(c.IConsumer, *mappedSsrc)
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
	if now.Sub(c.GetLastRtcpSentTime()) < c.GetMaxRtcpInterval() {
		return nil
	}
	c.SetLastRtcpSentTime(now)
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

func (c *SimpleConsumer) FillJsonStats() *FBS__Consumer.GetStatsResponseT {
	pStat := &FBS__Consumer.GetStatsResponseT{
		Stats: make([]*FBS__RtpStream.StatsT, 0),
	}
	if c.rtpStream != nil {
		stat := &FBS__RtpStream.StatsT{}
		c.rtpStream.FillJsonStats(stat)
	} else {
		c.logger.Warn().Msgf("rtpStream empty")
	}
	return pStat
}

func (c *SimpleConsumer) NeedWorstRemoteFractionLost(worstRemoteFractionLost *uint8) {
	// If our fraction lost is worse than the given one, update it.
	if c.rtpStream.fractionLost > *worstRemoteFractionLost {
		*worstRemoteFractionLost = c.rtpStream.fractionLost
	}
}
