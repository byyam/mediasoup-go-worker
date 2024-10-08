package rtc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	FBS__Producer "github.com/byyam/mediasoup-go-worker/fbs/FBS/Producer"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Response "github.com/byyam/mediasoup-go-worker/fbs/FBS/Response"
	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
	FBS__RtpStream "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpStream"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	"github.com/byyam/mediasoup-go-worker/internal/ms_rtcp"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/kr/pretty"

	"github.com/pion/rtcp"

	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type Producer struct {
	id     string
	logger zerolog.Logger

	Kind                  FBS__RtpParameters.MediaKind
	RtpParametersFBS      *FBS__RtpParameters.RtpParametersT
	RtpParameters         *mediasoupdata.RtpParameters
	Type                  FBS__RtpParameters.Type
	RtpHeaderExtensionIds RtpHeaderExtensionIds
	Paused                bool

	rtpStreamByEncodingIdx []*RtpStreamRecv
	rtpStreamScores        []uint8
	rtpMapping             *FBS__RtpParameters.RtpMappingT
	mapMappedSsrcSsrc      sync.Map
	mapSsrcRtpStream       sync.Map
	mapRtxSsrcRtpStream    sync.Map
	mapRtpStreamMappedSsrc sync.Map

	keyFrameRequestManager *KeyFrameRequestManager
	maxRtcpInterval        time.Duration
	lastRtcpSentTime       time.Time

	// handler
	onProducerRtpPacketReceivedHandler           atomic.Value
	onProducerSendRtcpPacketHandler              func(packet rtcp.Packet)
	onTransportProducerNewRtpStreamHandler       func(producerId string, rtpStream *RtpStreamRecv, mappedSsrc uint32)
	onProducerNeedWorstRemoteFractionLostHandler func(producerId string, worstRemoteFractionLost *uint8)
}

type ReceiveRtpPacketResult int

const (
	ReceiveRtpPacketResultDISCARDED ReceiveRtpPacketResult = iota
	ReceiveRtpPacketResultMEDIA
	ReceiveRtpPacketResultRETRANSMISSION
)

type producerParam struct {
	id                                    string
	options                               mediasoupdata.ProducerOptions
	optionsFBS                            *FBS__Transport.ProduceRequestT
	OnProducerRtpPacketReceived           func(*Producer, *rtpparser.Packet)
	OnProducerSendRtcpPacket              func(packet rtcp.Packet)
	OnProducerNeedWorstRemoteFractionLost func(producerId string, worstRemoteFractionLost *uint8)
}

func newProducer(param producerParam) (*Producer, error) {
	p := &Producer{
		id:     param.id,
		logger: zerowrapper.NewScope("producer", param.id),

		Kind:                   param.optionsFBS.Kind,
		RtpParametersFBS:       param.optionsFBS.RtpParameters,
		Paused:                 param.optionsFBS.Paused,
		RtpParameters:          mediasoupdata.NewRtpParameters(param.optionsFBS.RtpParameters),
		rtpStreamByEncodingIdx: make([]*RtpStreamRecv, len(param.optionsFBS.RtpParameters.Encodings)),
		rtpStreamScores:        make([]uint8, len(param.optionsFBS.RtpParameters.Encodings)),
		rtpMapping:             param.optionsFBS.RtpMapping,
	}
	p.onProducerRtpPacketReceivedHandler.Store(param.OnProducerRtpPacketReceived)
	p.onProducerSendRtcpPacketHandler = param.OnProducerSendRtcpPacket
	p.onProducerNeedWorstRemoteFractionLostHandler = param.OnProducerNeedWorstRemoteFractionLost

	p.logger.Debug().Any("rtpMapping", p.rtpMapping).Msgf("init rtp mapping for producer")
	// init producer with param
	if err := p.init(param); err != nil {
		return nil, err
	}

	p.logger.Info().Any("RtpParameters", mediasoupdata.JsonFormat(p.RtpParameters)).Msgf("init rtp param for producer")
	workerchannel.RegisterHandler(p.id, p.HandleRequest)
	return p, nil
}

func (p *Producer) init(param producerParam) error {
	if err := p.RtpParameters.Init(); err != nil {
		return err
	}
	if err := p.RtpHeaderExtensionIds.set(p.RtpParameters.HeaderExtensions, true); err != nil {
		p.logger.Error().Err(err).Msg("set RtpHeaderExtensionIds failed")
		return err
	}
	p.logger.Debug().Msgf("set RtpHeaderExtensionIds:%# v", pretty.Formatter(p.RtpHeaderExtensionIds))

	if p.Kind == FBS__RtpParameters.MediaKindAUDIO {
		p.maxRtcpInterval = ms_rtcp.MaxAudioIntervalMs
	} else {
		p.maxRtcpInterval = ms_rtcp.MaxVideoIntervalMs
		p.keyFrameRequestManager = NewKeyFrameRequestManager(&KeyFrameRequestManagerParam{
			keyFrameRequestDelay: param.optionsFBS.KeyFrameRequestDelay,
			onKeyFrameNeeded:     p.OnKeyFrameNeeded,
		})
	}
	p.Type = p.RtpParameters.GetType()
	return nil
}

func (p *Producer) ReceiveRtpPacket(packet *rtpparser.Packet) (result ReceiveRtpPacketResult) {
	if p.Kind == FBS__RtpParameters.MediaKindVIDEO && packet.IsKeyFrame() {
		monitor.RtcpCountBySSRC(packet.SSRC, monitor.KeyframePkgRecv)
		p.logger.Debug().Msg("isKeyFrame")
	}
	if p.Kind == FBS__RtpParameters.MediaKindVIDEO {
		monitor.RtpRecvCount(packet.SSRC, monitor.TraceVideo, packet.GetLen())
	} else if p.Kind == FBS__RtpParameters.MediaKindAUDIO {
		monitor.RtpRecvCount(packet.SSRC, monitor.TraceAudio, packet.GetLen())
	}

	rtpStream := p.GetRtpStream(packet)
	if rtpStream == nil {
		p.logger.Warn().Str("packet", packet.String()).Str("mid", packet.GetMid()).Str("rrid", packet.GetRrid()).
			Str("rid", packet.GetRid()).Msg("no stream found for received packet")
		monitor.RtpRecvCount(packet.SSRC, monitor.TraceRtpStreamNotFound, packet.GetLen())
		return ReceiveRtpPacketResultDISCARDED
	}
	// Pre-process the packet.
	p.PreProcessRtpPacket(packet)

	if rtpStream.GetSsrc() == packet.SSRC { // Media packet.
		result = ReceiveRtpPacketResultMEDIA
		monitor.RtpRecvCount(packet.SSRC, monitor.TraceRtpStream, packet.GetLen())
		if !rtpStream.ReceivePacket(packet) { // Process the packet.
			// todo
			monitor.RtpRecvCount(packet.SSRC, monitor.TraceRtpStreamRecvFailed, packet.GetLen())
			return
		}
	} else if rtpStream.GetRtxSsrc() == packet.SSRC { // RTX packet.
		result = ReceiveRtpPacketResultRETRANSMISSION
		monitor.RtpRecvCount(packet.SSRC, monitor.TraceRtpRtxStream, packet.GetLen())
		if !rtpStream.ReceiveRtxPacket(packet) {
			monitor.RtpRecvCount(packet.SSRC, monitor.TraceRtpRtxStreamRecvFailed, packet.GetLen())
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
		p.logger.Warn().Msgf("unknown payload type [payloadType:%d]", payloadType)
		return false
	}
	packet.PayloadType = mappedPayloadType

	// Mangle the SSRC.
	v, ok := p.mapRtpStreamMappedSsrc.Load(rtpStream.GetId())
	if !ok {
		p.logger.Warn().Msgf("unknown rtpStream type [rtpStream:%s]", rtpStream.GetId())
		return false
	}
	mappedSsrc := v.(uint32)
	packet.SSRC = mappedSsrc

	// Mangle RTP header extensions.
	// Add urn:ietf:params:rtp-hdrext:sdes:mid.
	//{
	//	payload := packet.GetExtension(p.RtpHeaderExtensionIds.Mid)
	//	if payload == nil {
	//		p.logger.Warn("set RTP header extension MID[%d] failed,mappedSsrc:%d", p.RtpHeaderExtensionIds.Mid, mappedSsrc)
	//	} else {
	//		if err := packet.SetExtension(mediasoupdata.MID, payload); err != nil {
	//			p.logger.Warn("set RTP header extension MID failed:%s,mappedSsrc:%d", err, mappedSsrc)
	//		}
	//	}
	//}
	if p.Kind == FBS__RtpParameters.MediaKindAUDIO {
		// Proxy urn:ietf:params:rtp-hdrext:ssrc-audio-level.
		payload := packet.GetExtension(p.RtpHeaderExtensionIds.SsrcAudioLevel)
		if payload != nil {
			if err := packet.SetExtension(mediasoupdata.SSRC_AUDIO_LEVEL, payload); err != nil {
				p.logger.Warn().Msgf("set RTP header extension ssrc audio level failed:%s,mappedSsrc:%d", err, mappedSsrc)
			}
		}
	} else if p.Kind == FBS__RtpParameters.MediaKindVIDEO {
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

func (p *Producer) FillJsonStats() *FBS__Producer.GetStatsResponseT {
	pStat := &FBS__Producer.GetStatsResponseT{
		Stats: make([]*FBS__RtpStream.StatsT, 0),
	}
	for idx, rtpStream := range p.rtpStreamByEncodingIdx {
		if rtpStream == nil {
			p.logger.Warn().Msgf("rtpStream empty, idx=%d", idx)
			continue
		}
		stat := &FBS__RtpStream.StatsT{}
		rtpStream.FillJsonStats(stat)
		pStat.Stats = append(pStat.Stats, stat)
		p.logger.Info().Msgf("stat:%+v", *stat)
	}
	return pStat
}

func (p *Producer) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	p.logger.Debug().Str("request", request.String()).Msg("handle channel request")

	switch request.MethodType {
	case FBS__Request.MethodPRODUCER_DUMP:
		response.Data = p.FillJson()
	case FBS__Request.MethodPRODUCER_GET_STATS:
		// request body is null
		dataDump := p.FillJsonStats()
		// set rsp
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyProducer_GetStatsResponse,
			Value: dataDump,
		}
		response.RspBody = rspBody

	default:
		response.Err = mserror.ErrInvalidMethod
		p.logger.Error().Msgf("unknown method:%s", request.Method)
	}

}

func (p *Producer) FillJson() json.RawMessage {
	dumpData := FBS__Producer.DumpResponseT{
		Id:              p.id,
		Kind:            p.Kind,
		Type:            p.Type,
		RtpParameters:   p.RtpParametersFBS,
		RtpMapping:      p.rtpMapping,
		RtpStreams:      nil,
		Paused:          p.Paused,
		TraceEventTypes: nil,
	}
	data, _ := json.Marshal(&dumpData)
	p.logger.Debug().Msgf("dumpData:%+v", dumpData)
	return data
}

func (p *Producer) Close() {
	p.logger.Info().Msgf("producer:%s closed", p.id)
	workerchannel.UnregisterHandler(p.id)
}

func (p *Producer) RequestKeyFrame(mappedSsrc uint32) {
	v, ok := p.mapMappedSsrcSsrc.Load(mappedSsrc)
	if !ok {
		p.logger.Warn().Msgf("given mappedSsrc[%d] not found, ignoring", mappedSsrc)
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

	p.logger.Debug().Msgf("RequestKeyFrame:%d,%d", mappedSsrc, ssrc)
	// todo
	p.keyFrameRequestManager.KeyFrameNeeded(ssrc)
}

func (p *Producer) OnRtpStreamSendRtcpPacket(packet rtcp.Packet) {
	p.onProducerSendRtcpPacketHandler(packet)
}

func (p *Producer) CreateRtpStream(packet *rtpparser.Packet, mediaCodec *mediasoupdata.RtpCodecParameters, encodingIdx int) *RtpStreamRecv {
	ssrc := packet.SSRC
	if _, ok := p.mapSsrcRtpStream.Load(ssrc); ok {
		panic(fmt.Sprintf("RtpStream with given SSRC=%d already exists", ssrc))
	}
	if v := p.rtpStreamByEncodingIdx[encodingIdx]; v != nil {
		panic("RtpStream for given encoding index already exists,idx=" + strconv.Itoa(encodingIdx))
	}
	encoding := p.RtpParameters.Encodings[encodingIdx]
	encodingMapping := p.rtpMapping.Encodings[encodingIdx]
	p.logger.Info().Msgf("CreateRtpStream ssrc:%d,mappedSsrc:%d,encodingIdx:%d,rid:%s,payloadType:%d",
		ssrc, encodingMapping.MappedSsrc, encodingIdx, encoding.Rid, mediaCodec.PayloadType)

	params := &ParamRtpStream{
		EncodingIdx:    encodingIdx,
		Ssrc:           ssrc,
		PayloadType:    mediaCodec.PayloadType,
		MimeType:       mediaCodec.RtpCodecMimeType,
		ClockRate:      int(mediaCodec.ClockRate),
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
		Kind:           p.Kind,
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
			p.logger.Info().Msg("NACK supported")
			params.UseNack = true
		} else if !params.UsePli && fb.Type == "nack" && fb.Parameter == "pli" {
			p.logger.Info().Msg("PLI supported")
			params.UsePli = true
		} else if !params.UseFir && fb.Type == "ccm" && fb.Parameter == "fir" {
			p.logger.Info().Msg("FIR supported")
			params.UseFir = true
		}
	}

	// Create a RtpStreamRecv for receiving a media stream.
	rtpStream := newRtpStreamRecv(&ParamRtpStreamRecv{
		ParamRtpStream:            params,
		onRtpStreamSendRtcpPacket: p.OnRtpStreamSendRtcpPacket,
		sendNackDelayMs:           10,
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

		if isMediaPacket && encoding != nil && encoding.Ssrc != 0 && encoding.Ssrc == ssrc {
			rtpStream := p.CreateRtpStream(packet, mediaCodec, idx)
			p.logger.Info().Uint32("ssrc", ssrc).Uint8("pt", payloadType).Msg("[GetRtpStream]CreateRtpStream by SSRC")
			return rtpStream
		} else if isRtxPacket && encoding.Rtx != nil && encoding.Rtx.Ssrc == ssrc {
			v, ok := p.mapSsrcRtpStream.Load(encoding.Ssrc)
			// Ignore if no stream has been created yet for the corresponding encoding.
			if !ok {
				p.logger.Warn().Uint32("ssrc", encoding.Ssrc).Msg("ignoring RTX packet for not yet created RtpStream (ssrc lookup)")
				return nil
			}
			rtpStream := v.(*RtpStreamRecv)

			// Ensure no RTX ssrc was previously detected.
			if rtpStream.HasRtx() {
				p.logger.Debug().Msg("ignoring RTX packet with new ssrc (ssrc lookup)")
				return nil
			}

			// Update the stream RTX data.
			rtpStream.SetRtx(payloadType, ssrc)

			// Insert the new RTX ssrc into the map.
			p.mapRtxSsrcRtpStream.Store(ssrc, rtpStream)
			p.logger.Info().Uint32("ssrc", ssrc).Uint8("pt", payloadType).Msg("[GetRtpStream]Find RTX RtpStream by SSRC")
			return rtpStream
		}
	}

	// If not found, look for an encoding matching the packet RID value.
	var rid string
	if packet.ReadRid(&rid) {
		for idx, encoding := range p.RtpParameters.Encodings {
			if encoding.Rid != rid {
				continue
			}
			mediaCodec := p.RtpParameters.GetCodecForEncoding(encoding)
			rtxCodec := p.RtpParameters.GetRtxCodecForEncoding(encoding)
			var isMediaPacket, isRtxPacket bool
			if mediaCodec.PayloadType == payloadType {
				isMediaPacket = true
			}
			if rtxCodec != nil && rtxCodec.PayloadType == payloadType {
				isRtxPacket = true
			}
			if isMediaPacket {
				// Ensure no other stream already exists with same RID.
				ignore := false
				p.mapSsrcRtpStream.Range(func(key, value any) bool {
					rtpStream, ok := value.(*RtpStreamRecv)
					if !ok || rtpStream == nil {
						return true
					}
					if rtpStream.GetRid() == packet.GetRid() {
						p.logger.Warn().Msg("ignoring packet with unknown ssrc but already handled RID (RID lookup)")
						ignore = true
						return false
					}
					return true
				})
				if ignore {
					return nil
				}
				rtpStream := p.CreateRtpStream(packet, mediaCodec, idx)
				p.logger.Info().Str("RID", packet.GetRid()).Uint32("ssrc", ssrc).Uint8("pt", payloadType).
					Msg("[GetRtpStream]CreateRtpStream by RID")

				return rtpStream

			} else if isRtxPacket {
				var gotRtpStream *RtpStreamRecv
				// Ensure a stream already exists with same RID.
				p.mapSsrcRtpStream.Range(func(key, value any) bool {
					rtpStream, ok := value.(*RtpStreamRecv)
					if !ok || rtpStream == nil {
						return true
					}
					if rtpStream.GetRid() == rid { // RID or RRID
						// Ensure no RTX ssrc was previously detected.
						if rtpStream.rtxStream != nil {
							p.logger.Warn().Msg("ignoring RTX packet with new SSRC (RID lookup)")
							return false
						}
						// Update the stream RTX data.
						rtpStream.SetRtx(payloadType, ssrc)
						// Insert the new RTX ssrc into the map.
						p.mapRtxSsrcRtpStream.Store(ssrc, rtpStream)
						gotRtpStream = rtpStream
					}
					return true
				})
				if gotRtpStream != nil {
					p.logger.Info().Str("RID", packet.GetRid()).Str("RRID", packet.GetRrid()).
						Uint32("ssrc", ssrc).Uint8("pt", payloadType).Msg("[GetRtpStream]Find RTX RtpStream by RID or RRID")
					return gotRtpStream
				}
			}
		}
		p.logger.Warn().Msg("[GetRtpStream]ignoring packet with unknown RID (RID lookup)")
	}
	// If not found, and there is a single encoding without ssrc and RID, this
	// may be the media or RTX stream.
	if len(p.RtpParameters.Encodings) == 1 && p.RtpParameters.Encodings[0].Ssrc == 0 && p.RtpParameters.Encodings[0].Rid == "" {
		// todo
		p.logger.Warn().Msgf("[GetRtpStream]may be the media or RTX stream")
	}

	return nil
}

func (p *Producer) PreProcessRtpPacket(packet *rtpparser.Packet) {
	if p.Kind == FBS__RtpParameters.MediaKindVIDEO {
		packet.SetFrameMarking07ExtensionId(p.RtpHeaderExtensionIds.FrameMarking07)
		packet.SetFrameMarkingExtensionId(p.RtpHeaderExtensionIds.FrameMarking)
	}
}

func (p *Producer) PostProcessRtpPacket(packet *rtpparser.Packet) {
	if p.Kind == FBS__RtpParameters.MediaKindVIDEO {
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
		p.logger.Warn().Msgf("no associated RtpStream found [ssrc:%d]", ssrc)
		return
	}
	rtpStream := v.(*RtpStreamRecv)
	rtpStream.RequestKeyFrame()
}

func (p *Producer) ReceiveRtcpSenderReport(sr *rtcp.SenderReport, report *rtcp.ReceptionReport) {
	v, ok := p.mapSsrcRtpStream.Load(report.SSRC)
	if !ok {
		p.logger.Warn().Msgf("RtpStream not found [ssrc:%d]", report.SSRC)
		return
	}
	rtpStream := v.(*RtpStreamRecv)
	rtpStream.ReceiveRtcpSenderReport(sr, report)
}

func (p *Producer) GetRtcp(now time.Time) []rtcp.Packet {
	if now.Sub(p.lastRtcpSentTime) < p.maxRtcpInterval {
		return nil
	}
	p.lastRtcpSentTime = now
	var packets []rtcp.Packet

	p.mapSsrcRtpStream.Range(func(key, value interface{}) bool {
		rtpStream, ok := value.(*RtpStreamRecv)
		if !ok || rtpStream == nil {
			return true
		}
		var worstRemoteFractionLost uint8
		if rtpStream.params.UseInBandFec {
			// Notify the listener, so we'll get the worst remote fraction lost.
			p.onProducerNeedWorstRemoteFractionLostHandler(p.id, &worstRemoteFractionLost)
			if worstRemoteFractionLost > 0 {
				p.logger.Debug().Msgf("using worst remote fraction lost:%d", worstRemoteFractionLost)
			}
		}
		report := rtpStream.GetRtcpReceiverReport(now, worstRemoteFractionLost)
		if report != nil {
			packets = append(packets, &rtcp.ReceiverReport{
				Reports: []rtcp.ReceptionReport{*report},
			})
		}
		// todo: rtx
		return true
	})

	return packets
}
