package rtc

import (
	"github.com/byyam/mediasoup-go-worker/common"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
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

func (r *RtpHeaderExtensionIds) fill(headerExtensions []mediasoupdata.RtpHeaderExtensionParameters) error {
	fn := func(ext mediasoupdata.RtpHeaderExtensionParameters) {
		if r.Mid == 0 && ext.Type == mediasoupdata.MID {
			r.Mid = uint8(ext.Id)
		}
		// todo
	}

	for _, ext := range headerExtensions {
		if ext.Id == 0 {
			return common.ErrInvalidParam
		}
		fn(ext)
	}
	return nil
}
