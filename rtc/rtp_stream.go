package rtc

import (
	"strconv"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/pion/rtp"
)

type ParamRtpStream struct {
	EncodingIdx    int
	Ssrc           uint32
	PayloadType    uint8
	MimeType       mediasoupdata.RtpCodecMimeType
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
	id           string
	score        uint8
	rtxStream    *RtxStream
	params       *ParamRtpStream
	rtt          float64
	hasRtt       bool
	packetsLost  uint32
	fractionLost uint8
	logger       utils.Logger
}

func newRtpStream(param *ParamRtpStream, initialScore uint8) *RtpStream {
	id := strconv.FormatInt(int64(param.Ssrc), 10)
	return &RtpStream{
		id:     id,
		score:  initialScore,
		params: param,
		logger: utils.NewLogger("RtpStream", id),
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
		r.logger.Warn("replace RTX stream:%d", ssrc)
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
	params.MimeType.SubType = mediasoupdata.MimeSubTypeRTX
}

func (r *RtpStream) GetId() string {
	return r.id
}

func (r *RtpStream) GetSsrc() uint32 {
	return r.params.Ssrc
}

func (r *RtpStream) GetRtxSsrc() uint32 {
	return r.params.RtxSsrc
}

func (r *RtpStream) ReceivePacket(packet *rtp.Packet) bool {

	return true
}

func (r *RtpStream) FillJsonStats(stat *mediasoupdata.ProducerStat) {
	stat.Ssrc = r.GetSsrc()
	stat.RtxSsrc = r.GetRtxSsrc()
	stat.Rid = r.params.Rid
	stat.Kind = r.params.MimeType.Type2String()
	stat.MimeType = r.params.MimeType.MimeType
}
