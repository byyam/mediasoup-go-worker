package rtc

import (
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/google/uuid"
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
}

func newRtpStream(param *ParamRtpStream, initialScore uint8) *RtpStream {
	return &RtpStream{
		id:     uuid.New().String(),
		score:  initialScore,
		params: param,
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
		// todo
	}
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
	stat.Kind = r.params.MimeType.TypeStr
	stat.MimeType = r.params.MimeType.MimeType
}
