package rtc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byyam/mediasoup-go-worker/common"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/rtc/rtc_rtcp"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/kr/pretty"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
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
	rtpMapping             mediasoupdata.RtpMapping
	mapMappedSsrcSsrc      sync.Map
	mapSsrcRtpStream       sync.Map
	mapRtxSsrcRtpStream    sync.Map

	keyFrameRequestManager *KeyFrameRequestManager
	maxRtcpInterval        time.Duration

	// handler
	onProducerRtpPacketReceivedHandler atomic.Value
	onProducerSendRtcpPacketHandler    func(header *rtcp.Header, packet rtcp.Packet)
}

type producerParam struct {
	id                          string
	options                     mediasoupdata.ProducerOptions
	OnProducerRtpPacketReceived func(*Producer, *rtp.Packet)
	OnProducerSendRtcpPacket    func(header *rtcp.Header, packet rtcp.Packet)
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
	p.onProducerSendRtcpPacketHandler = param.OnProducerSendRtcpPacket

	p.logger.Info("input param for producer[%s]: %# v", param.id, pretty.Formatter(param.options))

	// init producer with param
	if err := p.init(param); err != nil {
		return nil, err
	}

	p.logger.Info("init param for producer[%s]: %# v", param.id, pretty.Formatter(p.RtpParameters))
	return p, nil
}

func (p *Producer) init(param producerParam) error {
	if err := p.RtpParameters.Init(); err != nil {
		return err
	}

	p.initRtpMapping(param.options.RtpMapping)
	if err := p.RtpHeaderExtensionIds.fill(param.options.RtpParameters.HeaderExtensions); err != nil {
		p.logger.Error("fill RtpHeaderExtensionIds failed:%v", err)
		return err
	}

	if p.Kind == mediasoupdata.MediaKind_Audio {
		p.maxRtcpInterval = rtc_rtcp.MaxAudioIntervalMs
	} else {
		p.maxRtcpInterval = rtc_rtcp.MaxVideoIntervalMs
		p.keyFrameRequestManager = NewKeyFrameRequestManager(param.options.KeyFrameRequestDelay)
	}
	return nil
}

func (p *Producer) initRtpMapping(rtpMapping mediasoupdata.RtpMapping) {
	p.rtpMapping.Encodings = rtpMapping.Encodings
	p.rtpMapping.Codecs = rtpMapping.Codecs

}

// todo: other codecs
func isKeyFrame(data []byte) bool {
	const (
		typeSTAPA       = 24
		typeSPS         = 7
		naluTypeBitmask = 0x1F
	)

	var word uint32

	payload := bytes.NewReader(data)
	if err := binary.Read(payload, binary.BigEndian, &word); err != nil {
		return false
	}

	naluType := (word >> 24) & naluTypeBitmask
	if naluType == typeSTAPA && word&naluTypeBitmask == typeSPS {
		return true
	} else if naluType == typeSPS {
		return true
	}

	return false
}

func (p *Producer) ReceiveRtpPacket(packet *rtp.Packet) {
	if p.Kind == mediasoupdata.MediaKind_Video && isKeyFrame(packet.Payload) {
		p.logger.Debug("isKeyFrame")
	}

	rtpStream := p.GetRtpStream(packet)
	if rtpStream == nil {
		p.logger.Warn("no stream found for received packet [ssrc:%d]", packet.SSRC)
		return
	}
	// Pre-process the packet.
	p.PreProcessRtpPacket(packet)

	if rtpStream.GetSsrc() == packet.SSRC { // Media packet.

	} else if rtpStream.GetRtxSsrc() == packet.SSRC { // RTX packet.

	} else { // Should not happen.
		panic("found stream does not match received packet")
	}

	if handler, ok := p.onProducerRtpPacketReceivedHandler.Load().(func(*Producer, *rtp.Packet)); ok && handler != nil {
		handler(p, packet)
	}
}

func (p *Producer) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		p.logger.Debug("method=%s,internal=%+v,response:%s", request.Method, request.Internal, response)
	}()

	switch request.Method {
	case mediasoupdata.MethodProducerDump:
		response.Data = p.FillJson()
	}

}

func (p *Producer) FillJson() json.RawMessage {
	dumpData := mediasoupdata.ProducerDump{
		Id:              p.id,
		Kind:            string(p.Kind),
		Type:            string(p.Type),
		RtpParameters:   p.RtpParameters,
		RtpMapping:      p.rtpMapping,
		Encodings:       mediasoupdata.RtpMappingEncoding{},
		RtpStreams:      nil,
		Paused:          p.Paused,
		TraceEventTypes: "",
	}
	data, _ := json.Marshal(&dumpData)
	p.logger.Debug("dumpData:%+v", dumpData)
	return data
}

func (p *Producer) Close() {
	p.logger.Info("producer:%s closed", p.id)
}

func (p *Producer) RequestKeyFrame(mappedSsrc uint32) {
	v, ok := p.mapMappedSsrcSsrc.Load(mappedSsrc)
	if !ok {
		p.logger.Warn("given mappedSsrc[%d] not found, ignoring", mappedSsrc)
		return
	}
	ssrc := v.(uint32)
	p.logger.Debug("RequestKeyFrame:%d,%d", mappedSsrc, ssrc)
	// todo
}

func (p *Producer) OnRtpStreamSendRtcpPacket(header *rtcp.Header, packet rtcp.Packet) {

}

func (p *Producer) CreateRtpStream(packet *rtp.Packet, mediaCodec *mediasoupdata.RtpCodecParameters, encoding *mediasoupdata.RtpEncodingParameters) *RtpStreamRecv {
	ssrc := packet.SSRC
	if _, ok := p.mapSsrcRtpStream.Load(ssrc); ok {
		panic("RtpStream with given SSRC already exists")
	}
	p.logger.Debug("CreateRtpStream ssrc:%d", ssrc)
	// todo

	return nil
}

func (p *Producer) GetRtpStream(packet *rtp.Packet) *RtpStreamRecv {
	ssrc := packet.SSRC
	payloadType := packet.PayloadType

	// If stream found in media ssrcs map, return it.
	{
		v, ok := p.mapSsrcRtpStream.Load(ssrc)
		if ok {
			return v.(*RtpStreamRecv)
		}
	}
	// If stream found in RTX ssrcs map, return it.
	{
		v, ok := p.mapRtxSsrcRtpStream.Load(ssrc)
		if ok {
			return v.(*RtpStreamRecv)
		}
	}

	// Otherwise, check our encodings and, if appropriate, create a new stream.

	// First, look for an encoding with matching media or RTX ssrc value.
	for _, encoding := range p.RtpParameters.Encodings {
		mediaCodec := p.RtpParameters.GetCodecForEncoding(encoding)
		rtxCodec := p.RtpParameters.GetRtxCodecForEncoding(encoding)
		var isMediaPacket, isRtxPacket bool
		if mediaCodec.PayloadType == payloadType {
			isMediaPacket = true
		}
		if rtxCodec != nil && rtxCodec.PayloadType == payloadType {
			isRtxPacket = true
		}

		if isMediaPacket && encoding.Ssrc == ssrc {
			rtpStream := p.CreateRtpStream(packet, mediaCodec, encoding)
			return rtpStream
		} else if isRtxPacket && encoding.Rtx != nil && encoding.Rtx.Ssrc == ssrc {
			// todo
		}
	}

	// If not found, look for an encoding matching the packet RID value.
	// todo
	p.logger.Warn("ignoring packet with unknown RID (RID lookup)")
	return nil
}

func (p *Producer) PreProcessRtpPacket(packet *rtp.Packet) {
	// todo
}
