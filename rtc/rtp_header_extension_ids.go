package rtc

import (
	"github.com/byyam/mediasoup-go-worker/mserror"
	mediasoupdata2 "github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
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

func (r *RtpHeaderExtensionIds) set(headerExtensions []mediasoupdata2.RtpHeaderExtensionParameters, isProducer bool) error {
	fn := func(ext mediasoupdata2.RtpHeaderExtensionParameters) {
		ext.Type = mediasoupdata2.GetRtpHeaderExtensionUri(ext.Uri)
		if r.Mid == 0 && ext.Type == mediasoupdata2.MID {
			r.Mid = uint8(ext.Id)
		}
		if r.Rid == 0 && ext.Type == mediasoupdata2.RTP_STREAM_ID {
			r.Rid = uint8(ext.Id)
		}
		if r.RRid == 0 && ext.Type == mediasoupdata2.REPAIRED_RTP_STREAM_ID {
			r.RRid = uint8(ext.Id)
		}
		if r.AbsSendTime == 0 && ext.Type == mediasoupdata2.ABS_SEND_TIME {
			r.AbsSendTime = uint8(ext.Id)
		}
		if r.TransportWideCc01 == 0 && ext.Type == mediasoupdata2.TRANSPORT_WIDE_CC_01 {
			r.TransportWideCc01 = uint8(ext.Id)
		}
		// NOTE: Remove this once framemarking draft becomes RFC.
		if r.FrameMarking07 == 0 && ext.Type == mediasoupdata2.FRAME_MARKING_07 && isProducer {
			r.FrameMarking07 = uint8(ext.Id)
		}
		if r.FrameMarking == 0 && ext.Type == mediasoupdata2.FRAME_MARKING && isProducer {
			r.FrameMarking = uint8(ext.Id)
		}
		if r.SsrcAudioLevel == 0 && ext.Type == mediasoupdata2.SSRC_AUDIO_LEVEL {
			r.SsrcAudioLevel = uint8(ext.Id)
		}
		if r.VideoOrientation == 0 && ext.Type == mediasoupdata2.VIDEO_ORIENTATION {
			r.VideoOrientation = uint8(ext.Id)
		}
		if r.TOffset == 0 && ext.Type == mediasoupdata2.TOFFSET && isProducer {
			r.TOffset = uint8(ext.Id)
		}
		if r.AbsCaptureTime == 0 && ext.Type == mediasoupdata2.ABS_CAPTURE_TIME && isProducer {
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
