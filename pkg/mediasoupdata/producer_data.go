package mediasoupdata

import (
	"strings"

	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
	FBS__RtpStream "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpStream"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
)

type ProducerOptions struct {
	/**
	 * Producer id (just for Router.pipeToRouter() method).
	 */
	Id string `json:"id,omitempty"`

	/**
	 * Media kind ('audio' or 'video').
	 */
	Kind MediaKind `json:"kind,omitempty"`

	/**
	 * RTP parameters defining what the endpoint is sending.
	 */
	RtpParameters *RtpParameters `json:"rtpParameters,omitempty"`

	/**
	 * Whether the producer must start in paused mode. Default false.
	 */
	Paused bool `json:"paused,omitempty"`

	/**
	 * Just for video. Time (in ms) before asking the sender for a new key frame
	 * after having asked a previous one. Default 0.
	 */
	KeyFrameRequestDelay uint32 `json:"keyFrameRequestDelay,omitempty"`

	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`

	RtpMapping *RtpMapping `json:"rtpMapping"`
}

func (o ProducerOptions) Convert() *FBS__Transport.ProduceRequestT {
	p := &FBS__Transport.ProduceRequestT{
		ProducerId:           o.Id,
		Kind:                 FBS__RtpParameters.EnumValuesMediaKind[strings.ToUpper(string(o.Kind))],
		KeyFrameRequestDelay: o.KeyFrameRequestDelay,
		Paused:               o.Paused,
		RtpParameters:        o.RtpParameters.Convert(),
	}
	return p
}

func (o ProducerOptions) Valid() bool {
	if o.Kind != MediaKind_Audio && o.Kind != MediaKind_Video {
		return false
	}
	if !o.RtpMapping.Valid() || !o.RtpParameters.Valid() {
		logger.Error().Msg("RtpMapping or RtpParameters invalid")
		return false
	}
	if len(o.RtpMapping.Encodings) != len(o.RtpParameters.Encodings) {
		logger.Error().Msg("rtpParameters.encodings size does not match rtpMapping.encodings size")
		return false
	}
	return true
}

/**
 * Producer type.
 */
type ProducerType = RtpParametersType

const (
	ProducerType_Simple    ProducerType = RtpParametersType_Simple
	ProducerType_Simulcast              = RtpParametersType_Simulcast
	ProducerType_Svc                    = RtpParametersType_Svc
)

type ProducerData struct {
	Kind                    MediaKind     `json:"kind,omitempty"`
	Type                    ProducerType  `json:"type,omitempty"`
	RtpParameters           RtpParameters `json:"rtpParameters,omitempty"`
	ConsumableRtpParameters RtpParameters `json:"consumableRtpParameters,omitempty"`
}

type ProducerStat struct {
	// Common to all RtpStreams.
	Type                 string  `json:"type,omitempty"`
	Timestamp            uint64  `json:"timestamp,omitempty"`
	Ssrc                 uint32  `json:"ssrc,omitempty"`
	RtxSsrc              uint32  `json:"rtxSsrc,omitempty"`
	Rid                  string  `json:"rid,omitempty"`
	Kind                 string  `json:"kind,omitempty"`
	MimeType             string  `json:"mimeType,omitempty"`
	PacketsLost          uint64  `json:"packetsLost,omitempty"`
	FractionLost         uint8   `json:"fractionLost,omitempty"`
	PacketsDiscarded     uint64  `json:"packetsDiscarded,omitempty"`
	PacketsRetransmitted uint64  `json:"packetsRetransmitted,omitempty"`
	PacketsRepaired      uint64  `json:"packetsRepaired,omitempty"`
	NackCount            uint64  `json:"nackCount,omitempty"`
	NackPacketCount      uint64  `json:"nackPacketCount,omitempty"`
	PliCount             uint64  `json:"pliCount,omitempty"`
	FirCount             uint64  `json:"firCount,omitempty"`
	Score                uint32  `json:"score,omitempty"`
	PacketCount          uint64  `json:"packetCount,omitempty"`
	ByteCount            uint64  `json:"byteCount,omitempty"`
	Bitrate              uint32  `json:"bitrate,omitempty"`
	RoundTripTime        float32 `json:"roundTripTime,omitempty"`
	RtxPacketsDiscarded  uint64  `json:"rtxPacketsDiscarded,omitempty"`

	// RtpStreamRecv specific.
	Jitter         uint32 `json:"jitter,omitempty"`
	BitrateByLayer H      `json:"bitrateByLayer,omitempty"`
}

func (p *ProducerStat) Set(fbs *FBS__RtpStream.StatsDataT) {
	switch fbs.Type {
	case FBS__RtpStream.StatsDataBaseStats:
		stat := &FBS__RtpStream.BaseStatsT{}
		_ = Clone(fbs.Value, stat)
		p.SetBase(stat)
	case FBS__RtpStream.StatsDataRecvStats:
		stat := &FBS__RtpStream.RecvStatsT{}
		_ = Clone(fbs.Value, stat)
		if stat.Base.Data.Type == FBS__RtpStream.StatsDataBaseStats {
			baseStat := &FBS__RtpStream.BaseStatsT{}
			_ = Clone(stat.Base.Data.Value, baseStat)
			p.SetBase(baseStat)
		}
		p.Type = "inbound-rtp"
		p.Jitter = stat.Jitter
		p.Bitrate = stat.Bitrate
		p.PacketCount = stat.PacketCount
		p.ByteCount = stat.ByteCount
		p.BitrateByLayer = make(H) // todo
	case FBS__RtpStream.StatsDataSendStats:
		stat := &FBS__RtpStream.SendStatsT{}
		_ = Clone(fbs.Value, stat)
		if stat.Base.Data.Type == FBS__RtpStream.StatsDataBaseStats {
			baseStat := &FBS__RtpStream.BaseStatsT{}
			_ = Clone(stat.Base.Data.Value, baseStat)
			p.SetBase(baseStat)
		}
		p.Type = "outbound-rtp"
		p.PacketCount = stat.PacketCount
		p.ByteCount = stat.ByteCount
		p.Bitrate = stat.Bitrate
	}
}

func (p *ProducerStat) SetBase(stat *FBS__RtpStream.BaseStatsT) {
	p.Rid = stat.Rid
	p.Timestamp = stat.Timestamp
	p.Ssrc = stat.Ssrc
	if stat.RtxSsrc != nil {
		p.RtxSsrc = *stat.RtxSsrc
	}
	p.MimeType = stat.MimeType
	p.Kind = strings.ToLower(FBS__RtpParameters.EnumNamesMediaKind[stat.Kind])
	p.PacketsLost = stat.PacketsLost
	p.FractionLost = stat.FractionLost
	p.PacketsDiscarded = stat.PacketsDiscarded
	p.PacketsRetransmitted = stat.PacketsRetransmitted
	p.PacketsRepaired = stat.PacketsRepaired
	p.NackPacketCount = stat.NackPacketCount
	p.NackCount = stat.NackCount
	p.PliCount = stat.PliCount
	p.FirCount = stat.FirCount
	p.Score = uint32(stat.Score)
	p.RtxPacketsDiscarded = stat.RtxPacketsDiscarded
	p.RoundTripTime = stat.RoundTripTime
}
