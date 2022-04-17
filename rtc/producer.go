package rtc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/kr/pretty"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/rtc/ms_rtcp"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/pion/rtcp"
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
	OnProducerRtpPacketReceived func(*Producer, *rtpparser.Packet)
	OnProducerSendRtcpPacket    func(packet rtcp.Packet)
}

func newProducer(param producerParam) (*Producer, error) {
	if ok := param.options.Valid(); !ok {
		return nil, mserror.ErrInvalidParam
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
	if err := p.RtpHeaderExtensionIds.set(param.options.RtpParameters.HeaderExtensions, true); err != nil {
		p.logger.Error("set RtpHeaderExtensionIds failed:%v", err)
		return err
	}
	p.logger.Info("set RtpHeaderExtensionIds:%# v", pretty.Formatter(p.RtpHeaderExtensionIds))

	if p.Kind == mediasoupdata.MediaKind_Audio {
		p.maxRtcpInterval = ms_rtcp.MaxAudioIntervalMs
	} else {
		p.maxRtcpInterval = ms_rtcp.MaxVideoIntervalMs
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

func (p *Producer) ReceiveRtpPacket(packet *rtpparser.Packet) (result ReceiveRtpPacketResult) {
	if p.Kind == mediasoupdata.MediaKind_Video && isKeyFrame(packet.Payload) {
		monitor.KeyframeCount(packet.SSRC, monitor.KeyframePkgRecv)
		p.logger.Debug("isKeyFrame")
	}
	if p.Kind == mediasoupdata.MediaKind_Video {
		monitor.RtpRecvCount(monitor.TraceVideo)
	} else if p.Kind == mediasoupdata.MediaKind_Audio {
		monitor.RtpRecvCount(monitor.TraceAudio)
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

	// Mangle the packet before providing the listener with it.
	if !p.MangleRtpPacket(packet, rtpStream) {
		return ReceiveRtpPacketResultDISCARDED
	}

	// Post-process the packet.
	p.PostProcessRtpPacket(packet)

	if handler, ok := p.onProducerRtpPacketReceivedHandler.Load().(func(*Producer, *rtpparser.Packet)); ok && handler != nil {
		handler(p, packet)
	}
	return
}

func (p *Producer) MangleRtpPacket(packet *rtpparser.Packet, rtpStream *RtpStreamRecv) bool {
	// Mangle the payload type.
	payloadType := packet.PayloadType
	var mappedPayloadType uint8
	for _, codec := range p.rtpMapping.Codecs { // todo: use map
		if codec.PayloadType == payloadType {
			mappedPayloadType = codec.MappedPayloadType
			break
		}
	}
	if mappedPayloadType == 0 {
		p.logger.Warn("unknown payload type [payloadType:%d]", payloadType)
		return false
	}
	packet.PayloadType = mappedPayloadType

	// Mangle the SSRC.
	v, ok := p.mapRtpStreamMappedSsrc.Load(rtpStream.GetId())
	if !ok {
		p.logger.Warn("unknown rtpStream type [rtpStream:%s]", rtpStream.GetId())
		return false
	}
	mappedSsrc := v.(uint32)
	packet.SSRC = mappedSsrc

	// Mangle RTP header extensions.
	// Add urn:ietf:params:rtp-hdrext:sdes:mid.
	{
		payload := packet.GetExtension(p.RtpHeaderExtensionIds.Mid)
		if payload == nil {
			p.logger.Warn("set RTP header extension MID failed,mappedSsrc:%d", mappedSsrc)
		} else {
			if err := packet.SetExtension(mediasoupdata.MID, payload); err != nil {
				p.logger.Warn("set RTP header extension MID failed:%s,mappedSsrc:%d", err, mappedSsrc)
			}
		}
	}
	if p.Kind == mediasoupdata.MediaKind_Audio {
		// Proxy urn:ietf:params:rtp-hdrext:ssrc-audio-level.
		payload := packet.GetExtension(p.RtpHeaderExtensionIds.SsrcAudioLevel)
		if payload != nil {
			if err := packet.SetExtension(mediasoupdata.SSRC_AUDIO_LEVEL, payload); err != nil {
				p.logger.Warn("set RTP header extension ssrc audio level failed:%s,mappedSsrc:%d", err, mappedSsrc)
			}
		}
	} else if p.Kind == mediasoupdata.MediaKind_Video {
		// todo
	}
	// Assign mediasoup RTP header extension ids (just those that mediasoup may
	// be interested in after passing it to the Router).
	packet.SetMidExtensionId(mediasoupdata.MID)
	packet.SetAbsSendTimeExtensionId(mediasoupdata.ABS_SEND_TIME)
	packet.SetTransportWideCc01ExtensionId(mediasoupdata.TRANSPORT_WIDE_CC_01)
	// NOTE: Remove this once framemarking draft becomes RFC.
	packet.SetFrameMarking07ExtensionId(mediasoupdata.FRAME_MARKING_07)
	packet.SetFrameMarkingExtensionId(mediasoupdata.FRAME_MARKING)
	packet.SetSsrcAudioLevelExtensionId(mediasoupdata.SSRC_AUDIO_LEVEL)
	packet.SetVideoOrientationExtensionId(mediasoupdata.VIDEO_ORIENTATION)
	return true
}

func (p *Producer) FillJsonStats() json.RawMessage {
	var jsonData []mediasoupdata.ProducerStat
	for idx, rtpStream := range p.rtpStreamByEncodingIdx {
		if rtpStream == nil {
			p.logger.Warn("rtpStream empty, idx=%d", idx)
			continue
		}
		stat := &mediasoupdata.ProducerStat{}
		rtpStream.FillJsonStats(stat)
		jsonData = append(jsonData, *stat)
		p.logger.Info("stat:%+v", *stat)
	}
	data, _ := json.Marshal(&jsonData)
	p.logger.Info("getStats:%+v", jsonData)
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

	default:
		response.Err = mserror.ErrInvalidMethod
		p.logger.Error("unknown method:%s", request.Method)
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

func (p *Producer) CreateRtpStream(packet *rtpparser.Packet, mediaCodec *mediasoupdata.RtpCodecParameters, encodingIdx int) *RtpStreamRecv {
	ssrc := packet.SSRC
	if _, ok := p.mapSsrcRtpStream.Load(ssrc); ok {
		panic("RtpStream with given SSRC already exists")
	}
	if v := p.rtpStreamByEncodingIdx[encodingIdx]; v != nil {
		panic("RtpStream for given encoding index already exists")
	}
	encoding := p.RtpParameters.Encodings[encodingIdx]
	encodingMapping := p.rtpMapping.Encodings[encodingIdx]
	p.logger.Info("CreateRtpStream ssrc:%d,mappedSsrc:%d,encodingIdx:%d,rid:%s,payloadType:%d",
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
			p.logger.Info("NACK supported")
			params.UseNack = true
		} else if !params.UsePli && fb.Type == "nack" && fb.Parameter == "pli" {
			p.logger.Info("PLI supported")
			params.UsePli = true
		} else if !params.UseFir && fb.Type == "ccm" && fb.Parameter == "fir" {
			p.logger.Info("FIR supported")
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

func (p *Producer) GetRtpStream(packet *rtpparser.Packet) *RtpStreamRecv {
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
	// If not found, and there is a single encoding without ssrc and RID, this
	// may be the media or RTX stream.

	return nil
}

func (p *Producer) PreProcessRtpPacket(packet *rtpparser.Packet) {
	if p.Kind == mediasoupdata.MediaKind_Video {
		packet.SetFrameMarking07ExtensionId(p.RtpHeaderExtensionIds.FrameMarking07)
		packet.SetFrameMarkingExtensionId(p.RtpHeaderExtensionIds.FrameMarking)
	}
}

func (p *Producer) PostProcessRtpPacket(packet *rtpparser.Packet) {
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

func (p *Producer) ReceiveRtcpSenderReport(report *rtcp.ReceptionReport) {
	v, ok := p.mapSsrcRtpStream.Load(report.SSRC)
	if !ok {
		p.logger.Warn("RtpStream not found [ssrc:%d]", report.SSRC)
		return
	}
	rtpStream := v.(*RtpStreamRecv)
	rtpStream.ReceiveRtcpSenderReport(report)
}
