package rtc

import (
	"encoding/json"
	"errors"

	"github.com/kr/pretty"
	"github.com/rs/zerolog"

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
	consumerParam
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

	c.logger.Info().Msgf("param: %# v", pretty.Formatter(param.consumerParam))
	if err := c.initParam(param.consumerParam); err != nil {
		return nil, err
	}
	// Create RtpStreamSend instance for sending a single stream to the remote.
	c.CreateRtpStream()

	workerchannel.RegisterHandler(param.id, c.HandleRequest)
	return c, nil
}

func (c *SimulcastConsumer) initParam(param consumerParam) error {
	if len(param.consumableRtpEncodings) <= 1 {
		return errors.New("invalid consumableRtpEncodings with size <= 1")
	}
	encodings := param.consumableRtpEncodings[0]
	// Ensure there are as many spatial layers as encodings.
	if int(encodings.SpatialLayers) != len(param.consumableRtpEncodings) {
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

func (c *SimulcastConsumer) FillJsonStats() json.RawMessage {
	var jsonData []mediasoupdata.ConsumerStat
	if c.rtpStream != nil {
		var stat mediasoupdata.ConsumerStat
		c.rtpStream.FillJsonStats(&stat)
		jsonData = append(jsonData, stat)
	}
	data, _ := json.Marshal(&jsonData)
	c.logger.Debug().Msgf("getStats:%+v", jsonData)
	return data
}

func (c *SimulcastConsumer) OnRtpStreamRetransmitRtpPacket(packet *rtpparser.Packet) {

}