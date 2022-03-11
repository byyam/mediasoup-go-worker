package rtc

import (
	"sync"
	"sync/atomic"

	"github.com/pion/rtp"

	"github.com/byyam/mediasoup-go-worker/common"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/utils"
)

type Producer struct {
	id     string
	logger utils.Logger

	Kind                  mediasoupdata.MediaKind
	RtpParameters         mediasoupdata.RtpParameters
	Type                  mediasoupdata.ProducerType
	RtpHeaderExtensionIds RtpHeaderExtensionIds
	Paused                bool

	rtpStreamByEncodingIdx []*RtpStreamRecv
	rtpStreamScores        []uint8
	rtpMapping             RtpMapping

	// handler
	onProducerRtpPacketReceivedHandler atomic.Value
}

type RtpEncodingMapping struct {
	rid        string
	ssrc       uint32
	mappedSsrc uint32
}

func (r *RtpEncodingMapping) copy(encoding mediasoupdata.RtpMappingEncoding) {
	r.rid = encoding.Rid
	r.ssrc = encoding.Ssrc
	r.mappedSsrc = encoding.MappedSsrc
}

type RtpMapping struct {
	codecs    sync.Map
	encodings []RtpEncodingMapping
}

type producerParam struct {
	id                          string
	options                     mediasoupdata.ProducerOptions
	OnProducerRtpPacketReceived func(*Producer, *rtp.Packet)
}

func newProducer(param producerParam) (*Producer, error) {
	if ok := param.options.Valid(); !ok {
		return nil, common.ErrInvalidParam
	}
	p := &Producer{
		id:     param.id,
		logger: utils.NewLogger("producer"),

		Kind:          param.options.Kind,
		RtpParameters: param.options.RtpParameters,
		Type:          param.options.RtpParameters.GetType(),
		Paused:        param.options.Paused,

		rtpStreamByEncodingIdx: make([]*RtpStreamRecv, 0),
		rtpStreamScores:        make([]uint8, 0),
	}
	p.onProducerRtpPacketReceivedHandler.Store(param.OnProducerRtpPacketReceived)
	// todo
	p.initRtpMapping(param.options.RtpMapping)
	if err := p.RtpHeaderExtensionIds.fill(param.options.RtpParameters.HeaderExtensions); err != nil {
		p.logger.Error("fill RtpHeaderExtensionIds failed:%v", err)
		return nil, err
	}

	return p, nil
}

func (p *Producer) initRtpMapping(rtpMapping mediasoupdata.RtpMapping) {
	for _, codec := range rtpMapping.Codecs {
		p.rtpMapping.codecs.Store(codec.PayloadType, codec.MappedPayloadType)
	}
	for _, encoding := range rtpMapping.Encodings {
		var e RtpEncodingMapping
		e.copy(encoding)
		p.rtpMapping.encodings = append(p.rtpMapping.encodings, e)
	}
}

func (p *Producer) ReceiveRtpPacket(packet *rtp.Packet) {
	if handler, ok := p.onProducerRtpPacketReceivedHandler.Load().(func(*Producer, *rtp.Packet)); ok && handler != nil {
		handler(p, packet)
	}
}
