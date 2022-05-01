package rtc

import (
	"strconv"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
)

type RtxStream struct {
	id     string
	params *ParamRtxStream
	logger utils.Logger
}

type ParamRtxStream struct {
	Ssrc        uint32
	PayloadType uint8
	MimeType    mediasoupdata.RtpCodecMimeType
	ClockRate   int
	RRid        string
	Cname       string
}

func newRtxStream(params *ParamRtxStream) (*RtxStream, error) {
	if params == nil {
		return nil, mserror.ErrInvalidParam
	}
	id := strconv.FormatInt(int64(params.Ssrc), 10)
	r := &RtxStream{
		params: params,
		logger: utils.NewLogger("RtxStream", id),
	}
	if params.MimeType.SubType != mediasoupdata.MimeSubTypeRTX {
		return nil, mserror.ErrSubTypeNotRtx
	}
	return r, nil
}

func (r *RtxStream) ReceivePacket(packet *rtpparser.Packet) bool {
	// seq := packet.SequenceNumber

	return true
}

func (r RtxStream) GetPacketsDiscarded() uint32 {
	return 0
}
