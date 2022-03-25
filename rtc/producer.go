package rtc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/kr/pretty"

	"github.com/byyam/mediasoup-go-worker/common"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/rtc/rtc_rtcp"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
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
	mapRtpStreamMappedSsrc sync.Map

	keyFrameRequestManager *KeyFrameRequestManager
	maxRtcpInterval        time.Duration

	// handler
	onProducerRtpPacketReceivedHandler     atomic.Value
	onProducerSendRtcpPacketHandler        func(packet rtcp.Packet)
	onTransportProducerNewRtpStreamHandler func(producerId string, rtpStream *RtpStreamRecv, mappedSsrc uint32)
}

type ReceiveRtpPacketResult int

const (
	ReceiveRtpPacketResultDISCARDED ReceiveRtpPacketResult = iota
	ReceiveRtpPacketResultMEDIA
	ReceiveRtpPacketResultRETRANSMISSION
)

type producerParam struct {
	id                          string
	options                     mediasoupdata.ProducerOptions
	OnProducerRtpPacketReceived func(*Producer, *rtp.Packet)
	OnProducerSendRtcpPacket    func(packet rtcp.Packet)
}

func newProducer(param producerParam) (*Producer, error) {
	if ok := param.options.Valid(); !ok {
		return nil, common.ErrInvalidParam
	}
	p := &Producer{
		id:     param.id,
		logger: utils.NewLogger("producer", param.id),

		Kind:          param.options.Kind,
		RtpParameters: param.options.RtpParameters,
		Type:          param.options.RtpParameters.GetType(),
		Paused:        param.options.Paused,

		rtpStreamByEncodingIdx: make([]*RtpStreamRecv, len(param.options.RtpParameters.Encodings)),
		rtpStreamScores:        make([]uint8, len(param.options.RtpParameters.Encodings)),
	}
	p.onProducerRtpPacketReceivedHandler.Store(param.OnProducerRtpPacketReceived)
	p.onProducerSendRtcpPacketHandler = param.OnProducerSendRtcpPacket

	p.logger.Info("input param for producer: %# v", pretty.Formatter(param.options))

	// init producer with param
	if err := p.init(param); err != nil {
		return nil, err
	}

	p.logger.Info("init param for producer: %# v", pretty.Formatter(p.RtpParameters))
	p.logger.Info("init param for producer: %# v", pretty.Formatter(p.rtpMapping))
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
		p.keyFrameRequestManager = NewKeyFrameRequestManager(&KeyFrameRequestManagerParam{
			keyFrameRequestDelay: param.options.KeyFrameRequestDelay,
			onKeyFrameNeeded:     p.OnKeyFrameNeeded,
		})
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

func (p *Producer) ReceiveRtpPacket(packet *rtp.Packet) (result ReceiveRtpPacketResult) {
	if p.Kind == mediasoupdata.MediaKind_Video && isKeyFrame(packet.Payload) {
		monitor.KeyframeCount(packet.SSRC, monitor.KeyframePkgRecv)
		p.logger.Debug("isKeyFrame")
	}

	rtpStream := p.GetRtpStream(packet)
	if rtpStream == nil {
		p.logger.Warn("no stream found for received packet [ssrc:%d]", packet.SSRC)
		monitor.RtpRecvCount(monitor.TraceRtpStreamNotFound)
		return ReceiveRtpPacketResultDISCARDED
	}
	// Pre-process the packet.
	p.PreProcessRtpPacket(packet)

	if rtpStream.GetSsrc() == packet.SSRC { // Media packet.
		result = ReceiveRtpPacketResultMEDIA
		if !rtpStream.ReceivePacket(packet) { // Process the packet.
			// todo
			monitor.RtpRecvCount(monitor.TraceRtpStreamRecvFailed)
			return
		}
	} else if rtpStream.GetRtxSsrc() == packet.SSRC { // RTX packet.
		result = ReceiveRtpPacketResultRETRANSMISSION
		if !rtpStream.ReceiveRtxPacket(packet) {
			monitor.RtpRecvCount(monitor.TraceRtpRtxStreamRecvFailed)
			return
		}
	} else { // Should not happen.
		panic("found stream does not match received packet")
	}

	// If paused stop here.
	if p.Paused {
		return
	}

	// Post-process the packet.
	p.PostProcessRtpPacket(packet)

	if handler, ok := p.onProducerRtpPacketReceivedHandler.Load().(func(*Producer, *rtp.Packet)); ok && handler != nil {
		handler(p, packet)
	}
	return
}

func (p *Producer) FillJsonStats() json.RawMessage {
	jsonData := mediasoupdata.ProducerStat{
		Type:                 "",
		Timestamp:            0,
		Ssrc:                 0,
		RtxSsrc:              0,
		Rid:                  "",
		Kind:                 "",
		MimeType:             "",
		PacketsLost:          0,
		FractionLost:         0,
		PacketsDiscarded:     0,
		PacketsRetransmitted: 0,
		PacketsRepaired:      0,
		NackCount:            0,
		NackPacketCount:      0,
		PliCount:             0,
		FirCount:             0,
		Score:                0,
		PacketCount:          0,
		ByteCount:            0,
		Bitrate:              0,
		RoundTripTime:        0,
		RtxPacketsDiscarded:  0,
		Jitter:               0,
		BitrateByLayer:       nil,
	}
	data, _ := json.Marshal(&jsonData)
	p.logger.Debug("getStats:%+v", jsonData)
	return data
}

func (p *Producer) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		p.logger.Debug("method=%s,internal=%+v,response:%s", request.Method, request.Internal, response)
	}()

	switch request.Method {
	case mediasoupdata.MethodProducerDump:
		response.Data = p.FillJson()
	case mediasoupdata.MethodProducerGetStats:
		response.Data = p.FillJsonStats()
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

	// If the current RTP packet is a key frame for the given mapped SSRC do
	// nothing since we are gonna provide Consumers with the requested key frame
	// right now.
	//
	// NOTE: We know that this may only happen before calling MangleRtpPacket()
	// so the SSRC of the packet is still the original one and not the mapped one.
	//

	p.logger.Debug("RequestKeyFrame:%d,%d", mappedSsrc, ssrc)
	// todo
	p.keyFrameRequestManager.KeyFrameNeeded(ssrc)
}

func (p *Producer) OnRtpStreamSendRtcpPacket(packet rtcp.Packet) {
	p.onProducerSendRtcpPacketHandler(packet)
}

func (p *Producer) CreateRtpStream(packet *rtp.Packet, mediaCodec *mediasoupdata.RtpCodecParameters, encodingIdx int) *RtpStreamRecv {
	ssrc := packet.SSRC
	if _, ok := p.mapSsrcRtpStream.Load(ssrc); ok {
		panic("RtpStream with given SSRC already exists")
	}
	if v := p.rtpStreamByEncodingIdx[encodingIdx]; v != nil {
		panic("RtpStream for given encoding index already exists")
	}
	encoding := p.RtpParameters.Encodings[encodingIdx]
	encodingMapping := p.rtpMapping.Encodings[encodingIdx]
	p.logger.Debug("CreateRtpStream ssrc:%d,mappedSsrc:%d,encodingIdx:%d,rid:%s,payloadType:%d",
		ssrc, encodingMapping.MappedSsrc, encodingIdx, encoding.Rid, mediaCodec.PayloadType)

	params := &ParamRtpStream{
		EncodingIdx:    encodingIdx,
		Ssrc:           ssrc,
		PayloadType:    mediaCodec.PayloadType,
		MimeType:       mediaCodec.RtpCodecMimeType,
		ClockRate:      mediaCodec.ClockRate,
		Rid:            encoding.Rid,
		Cname:          p.RtpParameters.Rtcp.Cname,
		RtxSsrc:        0,
		RtxPayloadType: 0,
		UseNack:        false,
		UsePli:         false,
		UseFir:         false,
		UseInBandFec:   false,
		UseDtx:         false,
		SpatialLayers:  encoding.SpatialLayers,
		TemporalLayers: encoding.TemporalLayers,
	}
	// Check in band FEC in codec parameters.
	if mediaCodec.Parameters.Useinbandfec == 1 {
		params.UseInBandFec = true
	}
	// Check DTX in codec parameters.
	if mediaCodec.Parameters.Usedtx == 1 {
		params.UseDtx = true
	}
	// Check DTX in the encoding.
	if encoding.Dtx {
		params.UseDtx = true
	}
	for _, fb := range mediaCodec.RtcpFeedback {
		if !params.UseNack && fb.Type == "nack" && fb.Parameter == "" {
			params.UseNack = true
		} else if !params.UsePli && fb.Type == "nack" && fb.Parameter == "pli" {
			params.UsePli = true
		} else if !params.UseFir && fb.Type == "ccm" && fb.Parameter == "fir" {
			params.UseFir = true
		}
	}

	// Create a RtpStreamRecv for receiving a media stream.
	rtpStream := newRtpStreamRecv(&ParamRtpStreamRecv{
		ParamRtpStream:            params,
		onRtpStreamSendRtcpPacket: p.OnRtpStreamSendRtcpPacket,
	})

	// Insert into the maps.
	p.mapSsrcRtpStream.Store(ssrc, rtpStream)
	p.rtpStreamByEncodingIdx[encodingIdx] = rtpStream
	p.rtpStreamScores[encodingIdx] = rtpStream.GetScore()

	// Set the mapped SSRC.
	p.mapRtpStreamMappedSsrc.Store(rtpStream.GetId(), encodingMapping.MappedSsrc)
	p.mapMappedSsrcSsrc.Store(encodingMapping.MappedSsrc, ssrc)

	// If the Producer is paused tell it to the new RtpStreamRecv.
	if p.Paused {
		rtpStream.Pause()
	}

	return rtpStream
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
	for idx, encoding := range p.RtpParameters.Encodings {
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
			rtpStream := p.CreateRtpStream(packet, mediaCodec, idx)
			return rtpStream
		} else if isRtxPacket && encoding.Rtx != nil && encoding.Rtx.Ssrc == ssrc {
			v, ok := p.mapSsrcRtpStream.Load(encoding.Ssrc)
			// Ignore if no stream has been created yet for the corresponding encoding.
			if !ok {
				p.logger.Debug("ignoring RTX packet for not yet created RtpStream (ssrc lookup)")
				return nil
			}
			rtpStream := v.(*RtpStreamRecv)

			// Ensure no RTX ssrc was previously detected.
			if rtpStream.HasRtx() {
				p.logger.Debug("ignoring RTX packet with new ssrc (ssrc lookup)")
				return nil
			}

			// Update the stream RTX data.
			rtpStream.SetRtx(payloadType, ssrc)

			// Insert the new RTX ssrc into the map.
			p.mapRtxSsrcRtpStream.Store(ssrc, rtpStream)

			return rtpStream
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

func (p *Producer) PostProcessRtpPacket(packet *rtp.Packet) {
	if p.Kind == mediasoupdata.MediaKind_Video {
	}
	// todo
}

func (p *Producer) NotifyNewRtpStream(rtpStream *RtpStreamRecv) {
	v, ok := p.mapRtpStreamMappedSsrc.Load(rtpStream.GetId())
	if !ok {
		return
	}
	mappedSsrc := v.(uint32)
	p.onTransportProducerNewRtpStreamHandler(p.id, rtpStream, mappedSsrc)
}

func (p *Producer) OnKeyFrameNeeded(ssrc uint32) {
	v, ok := p.mapSsrcRtpStream.Load(ssrc)
	if !ok {
		p.logger.Warn("no associated RtpStream found [ssrc:%d]", ssrc)
		return
	}
	rtpStream := v.(*RtpStreamRecv)
	rtpStream.RequestKeyFrame()
}
