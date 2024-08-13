package rtc

import (
	"errors"
	"time"

	"github.com/pion/rtcp"
	"github.com/rs/zerolog"

	FBS__Consumer "github.com/byyam/mediasoup-go-worker/fbs/FBS/Consumer"
	FBS__RtpStream "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpStream"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type SimulcastConsumer struct {
	IConsumer
	logger     zerolog.Logger
	rtpStream  *RtpStreamSend
	rtpStreams []*RtpStreamSend

	// handler
	onConsumerSendRtpPacketHandler     func(consumer IConsumer, packet *rtpparser.Packet)
	onConsumerKeyFrameRequestedHandler func(consumer IConsumer, mappedSsrc uint32)
}

type simulcastConsumerParam struct {
	*consumerParam
	OnConsumerSendRtpPacket       func(consumer IConsumer, packet *rtpparser.Packet)
	OnConsumerKeyFrameRequested   func(consumer IConsumer, mappedSsrc uint32)
	OnConsumerRetransmitRtpPacket func(packet *rtpparser.Packet)
}

func newSimulcastConsumer(param simulcastConsumerParam) (*SimulcastConsumer, error) {
	var err error
	c := &SimulcastConsumer{
		rtpStreams: make([]*RtpStreamSend, 0),
		logger:     zerowrapper.NewScope("simulcast-consumer", param.id),
	}
	param.fillJsonStatsFunc = c.FillJsonStats
	c.IConsumer, err = newConsumer(mediasoupdata.ConsumerType_Simulcast, param.consumerParam)
	c.onConsumerSendRtpPacketHandler = param.OnConsumerSendRtpPacket
	c.onConsumerKeyFrameRequestedHandler = param.OnConsumerKeyFrameRequested
	if err != nil {
		return nil, err
	}

	if err := c.initParam(param.consumerParam); err != nil {
		return nil, err
	}
	// Create RtpStreamSend instance for sending a single stream to the remote.
	c.CreateRtpStream()

	workerchannel.RegisterHandler(param.id, c.HandleRequest)
	return c, nil
}

func (c *SimulcastConsumer) initParam(param *consumerParam) error {
	c.logger.Info().Any("consumableRtpEncodings", mediasoupdata.JsonFormat(param.consumableRtpEncodings)).
		Any("rtpParameters", mediasoupdata.JsonFormat(param.rtpParameters)).
		Msgf("initParam")
	if len(param.consumableRtpEncodings) <= 1 || len(param.rtpParameters.Encodings) == 0 {
		return errors.New("invalid consumableRtpEncodings with size <= 1")
	}
	encodings := param.rtpParameters.Encodings[0]
	// Ensure there are as many spatial layers as encodings.
	if int(encodings.ParsedScalabilityMode.SpatialLayers) != len(param.consumableRtpEncodings) {
		return errors.New("encoding.spatialLayers does not match number of consumableRtpEncodings")
	}

	return nil
}

func (c *SimulcastConsumer) CreateRtpStream() {
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
		//todo Kind
	}
	c.rtpStream = newRtpStreamSend(&ParamRtpStreamSend{
		ParamRtpStream:                 param,
		bufferSize:                     0,
		OnRtpStreamRetransmitRtpPacket: c.OnRtpStreamRetransmitRtpPacket,
	})
	c.rtpStreams = append(c.rtpStreams, c.rtpStream)
}

func (c *SimulcastConsumer) FillJsonStats() *FBS__Consumer.GetStatsResponseT {
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

func (c *SimulcastConsumer) OnRtpStreamRetransmitRtpPacket(packet *rtpparser.Packet) {

}

func (c *SimulcastConsumer) SendRtpPacket(packet *rtpparser.Packet) {
	// todo
}

func (c *SimulcastConsumer) GetRtpStreams() []*RtpStreamSend {
	return c.rtpStreams
}

func (c *SimulcastConsumer) GetRtcp(rtpStream *RtpStreamSend, now time.Time) []rtcp.Packet {
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
