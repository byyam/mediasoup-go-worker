package mediasoupdata

import (
	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
)

type RtpHeaderExtensionUri uint8

const (
	UNKNOWN                RtpHeaderExtensionUri = 0
	MID                                          = 1
	RTP_STREAM_ID                                = 2
	REPAIRED_RTP_STREAM_ID                       = 3
	ABS_SEND_TIME                                = 4
	TRANSPORT_WIDE_CC_01                         = 5
	FRAME_MARKING_07                             = 6 // NOTE: Remove once RFC.
	FRAME_MARKING                                = 7
	SSRC_AUDIO_LEVEL                             = 10
	VIDEO_ORIENTATION                            = 11
	TOFFSET                                      = 12
	ABS_CAPTURE_TIME                             = 13
)

var extensionUri = map[string]RtpHeaderExtensionUri{
	"urn:ietf:params:rtp-hdrext:sdes:mid":                                       MID,
	"urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id":                             RTP_STREAM_ID,
	"urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id":                    REPAIRED_RTP_STREAM_ID,
	"http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time":                ABS_SEND_TIME,
	"http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01": TRANSPORT_WIDE_CC_01,
	// NOTE: Remove this once framemarking draft becomes RFC.
	"http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07":  FRAME_MARKING_07,
	"urn:ietf:params:rtp-hdrext:framemarking":                       FRAME_MARKING,
	"urn:ietf:params:rtp-hdrext:ssrc-audio-level":                   SSRC_AUDIO_LEVEL,
	"urn:3gpp:video-orientation":                                    VIDEO_ORIENTATION,
	"urn:ietf:params:rtp-hdrext:toffset":                            TOFFSET,
	"http://www.webrtc.org/experiments/rtp-hdrext/abs-capture-time": ABS_CAPTURE_TIME,
}

func GetRtpHeaderExtensionUri(uri string) RtpHeaderExtensionUri {
	v, ok := extensionUri[uri]
	if !ok {
		return UNKNOWN
	}
	return v
}

var extensionFbsUri = map[string]FBS__RtpParameters.RtpHeaderExtensionUri{
	"urn:ietf:params:rtp-hdrext:sdes:mid":                                       FBS__RtpParameters.RtpHeaderExtensionUriMid,
	"urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id":                             FBS__RtpParameters.RtpHeaderExtensionUriRtpStreamId,
	"urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id":                    FBS__RtpParameters.RtpHeaderExtensionUriRepairRtpStreamId,
	"http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time":                FBS__RtpParameters.RtpHeaderExtensionUriAbsSendTime,
	"http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01": FBS__RtpParameters.RtpHeaderExtensionUriTransportWideCcDraft01,
	// NOTE: Remove this once framemarking draft becomes RFC.
	"http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07":  FBS__RtpParameters.RtpHeaderExtensionUriFrameMarkingDraft07,
	"urn:ietf:params:rtp-hdrext:framemarking":                       FBS__RtpParameters.RtpHeaderExtensionUriFrameMarking,
	"urn:ietf:params:rtp-hdrext:ssrc-audio-level":                   FBS__RtpParameters.RtpHeaderExtensionUriAudioLevel,
	"urn:3gpp:video-orientation":                                    FBS__RtpParameters.RtpHeaderExtensionUriVideoOrientation,
	"urn:ietf:params:rtp-hdrext:toffset":                            FBS__RtpParameters.RtpHeaderExtensionUriTimeOffset,
	"http://www.webrtc.org/experiments/rtp-hdrext/abs-capture-time": FBS__RtpParameters.RtpHeaderExtensionUriAbsCaptureTime,
}

func GetFbsUri(uri string) FBS__RtpParameters.RtpHeaderExtensionUri {
	v, ok := extensionFbsUri[uri]
	if !ok {
		return FBS__RtpParameters.RtpHeaderExtensionUriMid
	}
	return v
}

var extensionUriFbs = map[FBS__RtpParameters.RtpHeaderExtensionUri]string{
	FBS__RtpParameters.RtpHeaderExtensionUriMid:                    "urn:ietf:params:rtp-hdrext:sdes:mid",
	FBS__RtpParameters.RtpHeaderExtensionUriRtpStreamId:            "urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id",
	FBS__RtpParameters.RtpHeaderExtensionUriRepairRtpStreamId:      "urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id",
	FBS__RtpParameters.RtpHeaderExtensionUriAbsSendTime:            "http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time",
	FBS__RtpParameters.RtpHeaderExtensionUriTransportWideCcDraft01: "http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01",
	// NOTE: Remove this once framemarking draft becomes RFC.
	FBS__RtpParameters.RtpHeaderExtensionUriFrameMarkingDraft07: "http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07",
	FBS__RtpParameters.RtpHeaderExtensionUriFrameMarking:        "urn:ietf:params:rtp-hdrext:framemarking",
	FBS__RtpParameters.RtpHeaderExtensionUriAudioLevel:          "urn:ietf:params:rtp-hdrext:ssrc-audio-level",
	FBS__RtpParameters.RtpHeaderExtensionUriVideoOrientation:    "urn:3gpp:video-orientation",
	FBS__RtpParameters.RtpHeaderExtensionUriTimeOffset:          "urn:ietf:params:rtp-hdrext:toffset",
	FBS__RtpParameters.RtpHeaderExtensionUriAbsCaptureTime:      "http://www.webrtc.org/experiments/rtp-hdrext/abs-capture-time",
}

func GetUriFbs(uri FBS__RtpParameters.RtpHeaderExtensionUri) string {
	v, ok := extensionUriFbs[uri]
	if !ok {
		return ""
	}
	return v
}
