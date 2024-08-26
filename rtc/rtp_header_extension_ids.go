package rtc

import (
	"fmt"

	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
)

type RtpHeaderExtensionIds struct {
	Mid               uint8
	Rid               uint8
	RRid              uint8
	AbsSendTime       uint8
	TransportWideCc01 uint8
	FrameMarking07    uint8
	FrameMarking      uint8
	SsrcAudioLevel    uint8
	VideoOrientation  uint8
	TOffset           uint8
	AbsCaptureTime    uint8
}

func (r *RtpHeaderExtensionIds) set(headerExtensions []*mediasoupdata.RtpHeaderExtensionParameters, isProducer bool) error {
	fn := func(ext *mediasoupdata.RtpHeaderExtensionParameters) {
		uri := mediasoupdata.GetFbsUri(ext.Uri)
		if r.Mid == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriMid {
			r.Mid = uint8(ext.Id)
		}
		if r.Rid == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriRtpStreamId {
			r.Rid = uint8(ext.Id)
		}
		if r.RRid == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriRepairRtpStreamId {
			r.RRid = uint8(ext.Id)
		}
		if r.AbsSendTime == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriAbsSendTime {
			r.AbsSendTime = uint8(ext.Id)
		}
		if r.TransportWideCc01 == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriTransportWideCcDraft01 {
			r.TransportWideCc01 = uint8(ext.Id)
		}
		// NOTE: Remove this once framemarking draft becomes RFC.
		if r.FrameMarking07 == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriFrameMarkingDraft07 && isProducer {
			r.FrameMarking07 = uint8(ext.Id)
		}
		if r.FrameMarking == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriFrameMarking && isProducer {
			r.FrameMarking = uint8(ext.Id)
		}
		if r.SsrcAudioLevel == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriAudioLevel {
			r.SsrcAudioLevel = uint8(ext.Id)
		}
		if r.VideoOrientation == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriVideoOrientation {
			r.VideoOrientation = uint8(ext.Id)
		}
		if r.TOffset == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriTimeOffset && isProducer {
			r.TOffset = uint8(ext.Id)
		}
		if r.AbsCaptureTime == 0 && uri == FBS__RtpParameters.RtpHeaderExtensionUriAbsCaptureTime && isProducer {
			r.AbsCaptureTime = uint8(ext.Id)
		}
	}

	for _, ext := range headerExtensions {
		if ext.Id == 0 {
			return mserror.ErrInvalidParam
		}
		fn(ext)
	}
	return nil
}

func (r *RtpHeaderExtensionIds) String() string {
	out := "RtpHeaderExtensionIds:\n"
	out += fmt.Sprintf("\tMID:%d\n", r.Mid)
	out += fmt.Sprintf("\tRID:%d\n", r.Rid)
	return out
}
