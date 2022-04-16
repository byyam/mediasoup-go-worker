package rtc

import (
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/mserror"
)

type RtxStream struct {
	params *ParamRtxStream
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
	r := &RtxStream{params: params}
	if params.MimeType.SubType == mediasoupdata.MimeSubTypeRTX {
		return nil, mserror.ErrSubTypeNotRtx
	}
	params.MimeType.UpdateMimeType()
	return r, nil
}
