package rtc

import (
	"strconv"
	"time"

	"github.com/rs/zerolog"

	mediasoupdata2 "github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/byyam/mediasoup-go-worker/pkg/seqmgr"

	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
)

const (
	MaxDropout           = 3000
	MaxMisOrder          = 1500
	RtpSeqMod            = 1 << 16
	ScoreHistogramLength = 24
)

type ParamRtpStream struct {
	EncodingIdx    int
	Ssrc           uint32
	PayloadType    uint8
	MimeType       mediasoupdata2.RtpCodecMimeType
	ClockRate      int
	Rid            string
	Cname          string
	RtxSsrc        uint32
	RtxPayloadType uint8
	UseNack        bool
	UsePli         bool
	UseFir         bool
	UseInBandFec   bool
	UseDtx         bool
	SpatialLayers  uint8
	TemporalLayers uint8
}

type RtpStream struct {
	id                   string
	score                uint8
	rtxStream            *RtxStream
	params               *ParamRtpStream
	rtt                  float64
	hasRtt               bool
	packetsLost          uint32
	fractionLost         uint8
	nackCount            uint32
	nackPacketCount      uint32
	maxPacketTs          uint32
	maxPacketMS          int64
	packetsRetransmitted uint32
	packetsRepaired      uint32
	packetsDiscarded     uint32
	pliCount             uint32
	firCount             uint32
	// Others.
	//   https://tools.ietf.org/html/rfc3550#appendix-A.1 stuff.
	maxSeq  uint16 // Highest seq. number seen.
	cycles  uint32 // Shifted count of seq. number cycles.
	baseSeq uint32 // Base seq number.
	badSeq  uint32 // Last 'bad' seq number + 1.
	started bool

	reportedPacketLost uint32

	logger zerolog.Logger
}

func newRtpStream(param *ParamRtpStream, initialScore uint8) *RtpStream {
	id := strconv.FormatInt(int64(param.Ssrc), 10)
	return &RtpStream{
		id:     id,
		score:  initialScore,
		params: param,
		logger: zerowrapper.NewScope("RtpStream", id),
	}
}

func (r *RtpStream) HasRtx() bool {
	if r.rtxStream != nil {
		return true
	}
	return false
}

func (r *RtpStream) SetRtx(payloadType uint8, ssrc uint32) {
	r.params.RtxPayloadType = payloadType
	r.params.RtxSsrc = ssrc

	if r.HasRtx() {
		r.logger.Warn().Msgf("replace RTX stream:%d", ssrc)
	}
	// Set RTX stream params.
	params := &ParamRtxStream{
		Ssrc:        ssrc,
		PayloadType: payloadType,
		MimeType:    r.params.MimeType,
		ClockRate:   r.params.ClockRate,
		RRid:        r.params.Rid,
		Cname:       r.params.Cname,
	}
	params.MimeType.SubType = mediasoupdata2.MimeSubTypeRTX
	// Tell the RtpCodecMimeType to update its string based on current type and subtype.
	params.MimeType.UpdateMimeType()
	var err error
	r.rtxStream, err = newRtxStream(params)
	if err != nil {
		r.logger.Error().Err(err).Msg("set rtx failed")
		return
	}
	r.logger.Info().Uint32("ssrc", ssrc).Msg("set RTX stream")
}

func (r *RtpStream) GetId() string {
	return r.id
}

func (r *RtpStream) GetSsrc() uint32 {
	return r.params.Ssrc
}

func (r *RtpStream) GetCname() string {
	return r.params.Cname
}

func (r *RtpStream) GetRtxSsrc() uint32 {
	return r.params.RtxSsrc
}

func (r *RtpStream) ReceivePacket(packet *rtpparser.Packet) bool {
	// If this is the first packet seen, initialize stuff.
	if !r.started {
		r.InitSeq(packet.SequenceNumber)
		r.started = true
		r.maxSeq = packet.SequenceNumber - 1
		r.maxPacketTs = packet.Timestamp
		r.maxPacketMS = utils.GetTimeMs()
	}
	// If not a valid packet ignore it.
	if !r.UpdateSeq(packet) {
		r.logger.Warn().Msgf("invalid packet [ssrc:%d,seq:%d]", packet.SSRC, packet.SequenceNumber)
		return false
	}
	// Update highest seen RTP timestamp.
	if seqmgr.IsSeqHigherThanUint32(packet.Timestamp, r.maxPacketTs) {
		r.maxPacketTs = packet.Timestamp
		r.maxPacketMS = utils.GetTimeMs()
	}
	return true
}

func (r *RtpStream) FillJsonStats(stat *mediasoupdata2.ProducerStat) {
	stat.Ssrc = r.GetSsrc()
	stat.RtxSsrc = r.GetRtxSsrc()
	stat.Rid = r.params.Rid
	stat.Kind = r.params.MimeType.Type2String()
	stat.MimeType = r.params.MimeType.MimeType
	stat.PacketsLost = r.packetsLost
	stat.FractionLost = r.fractionLost
	stat.PacketsRepaired = r.packetsRepaired
	stat.NackCount = r.nackCount
	stat.NackPacketCount = r.nackPacketCount
	stat.PliCount = r.pliCount
	stat.FirCount = r.firCount
	stat.Score = uint32(r.score)
	stat.Rid = r.params.Rid
	stat.RtxSsrc = r.params.RtxSsrc
	if r.HasRtx() {
		stat.RtxPacketsDiscarded = r.rtxStream.GetPacketsDiscarded()
	}
	stat.RoundTripTime = float32(r.rtt)
}

func (r *RtpStream) GetRtpTimestamp(now time.Time) uint32 {
	// Calculate TS difference between now and maxPacketMs.
	diffMs := uint32(now.UnixNano()/1000) - r.maxPacketTs
	diffTs := diffMs * uint32(r.params.ClockRate) / 1000
	return diffTs + r.maxPacketTs
}

func (r *RtpStream) GetExpectedPackets() uint32 {
	return r.cycles + uint32(r.maxSeq) - r.baseSeq + 1
}

func (r *RtpStream) InitSeq(seq uint16) {
	r.logger.Trace().Msg("init seq")
	// Initialize/reset RTP counters.
	r.baseSeq = uint32(seq)
	r.maxSeq = seq
	r.badSeq = RtpSeqMod + 1
}

func (r *RtpStream) UpdateSeq(packet *rtpparser.Packet) bool {
	udelta := packet.SequenceNumber - r.maxSeq
	// If the new packet sequence number is greater than the max seen but not
	// "so much bigger", accept it.
	// NOTE: udelta also handles the case of a new cycle, this is:
	//    maxSeq:65536, seq:0 => udelta:1
	if udelta < MaxDropout {
		// In order, with permissible gap.
		if packet.SequenceNumber < r.maxSeq {
			// Sequence number wrapped: count another 64K cycle.
			r.cycles += RtpSeqMod
		}
		r.maxSeq = packet.SequenceNumber
	} else if udelta <= RtpSeqMod-MaxMisOrder {
		// Too old packet received (older than the allowed misorder).
		// Or to new packet (more than acceptable dropout).

		// The sequence number made a very large jump. If two sequential packets
		// arrive, accept the latter.
		if uint32(packet.SequenceNumber) == r.badSeq {
			// Two sequential packets. Assume that the other side restarted without
			// telling us so just re-sync (i.e., pretend this was the first packet).
			r.logger.Warn().Msgf("too bad sequence number, re-syncing RTP [ssrc:%d,seq:%d]", packet.SSRC, packet.SequenceNumber)
			r.InitSeq(packet.SequenceNumber)
			r.maxPacketTs = packet.Timestamp
			r.maxPacketMS = utils.GetTimeMs()
		} else {
			r.logger.Warn().Msgf("bad sequence number, ignoring packet [ssrc:%d,seq:%d]", packet.SSRC, packet.SequenceNumber)
			r.badSeq = (uint32(packet.SequenceNumber) + 1) & (RtpSeqMod - 1)
			// Packet discarded due to late or early arriving.
			r.packetsDiscarded++
			return false
		}
	} else {
		// Acceptable misorder.
		// Do nothing.
	}
	return true
}
