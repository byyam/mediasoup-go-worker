// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpParameters

import "strconv"

type RtpHeaderExtensionUri byte

const (
	RtpHeaderExtensionUriMid                    RtpHeaderExtensionUri = 0
	RtpHeaderExtensionUriRtpStreamId            RtpHeaderExtensionUri = 1
	RtpHeaderExtensionUriRepairRtpStreamId      RtpHeaderExtensionUri = 2
	RtpHeaderExtensionUriFrameMarkingDraft07    RtpHeaderExtensionUri = 3
	RtpHeaderExtensionUriFrameMarking           RtpHeaderExtensionUri = 4
	RtpHeaderExtensionUriAudioLevel             RtpHeaderExtensionUri = 5
	RtpHeaderExtensionUriVideoOrientation       RtpHeaderExtensionUri = 6
	RtpHeaderExtensionUriTimeOffset             RtpHeaderExtensionUri = 7
	RtpHeaderExtensionUriTransportWideCcDraft01 RtpHeaderExtensionUri = 8
	RtpHeaderExtensionUriAbsSendTime            RtpHeaderExtensionUri = 9
	RtpHeaderExtensionUriAbsCaptureTime         RtpHeaderExtensionUri = 10
)

var EnumNamesRtpHeaderExtensionUri = map[RtpHeaderExtensionUri]string{
	RtpHeaderExtensionUriMid:                    "Mid",
	RtpHeaderExtensionUriRtpStreamId:            "RtpStreamId",
	RtpHeaderExtensionUriRepairRtpStreamId:      "RepairRtpStreamId",
	RtpHeaderExtensionUriFrameMarkingDraft07:    "FrameMarkingDraft07",
	RtpHeaderExtensionUriFrameMarking:           "FrameMarking",
	RtpHeaderExtensionUriAudioLevel:             "AudioLevel",
	RtpHeaderExtensionUriVideoOrientation:       "VideoOrientation",
	RtpHeaderExtensionUriTimeOffset:             "TimeOffset",
	RtpHeaderExtensionUriTransportWideCcDraft01: "TransportWideCcDraft01",
	RtpHeaderExtensionUriAbsSendTime:            "AbsSendTime",
	RtpHeaderExtensionUriAbsCaptureTime:         "AbsCaptureTime",
}

var EnumValuesRtpHeaderExtensionUri = map[string]RtpHeaderExtensionUri{
	"Mid":                    RtpHeaderExtensionUriMid,
	"RtpStreamId":            RtpHeaderExtensionUriRtpStreamId,
	"RepairRtpStreamId":      RtpHeaderExtensionUriRepairRtpStreamId,
	"FrameMarkingDraft07":    RtpHeaderExtensionUriFrameMarkingDraft07,
	"FrameMarking":           RtpHeaderExtensionUriFrameMarking,
	"AudioLevel":             RtpHeaderExtensionUriAudioLevel,
	"VideoOrientation":       RtpHeaderExtensionUriVideoOrientation,
	"TimeOffset":             RtpHeaderExtensionUriTimeOffset,
	"TransportWideCcDraft01": RtpHeaderExtensionUriTransportWideCcDraft01,
	"AbsSendTime":            RtpHeaderExtensionUriAbsSendTime,
	"AbsCaptureTime":         RtpHeaderExtensionUriAbsCaptureTime,
}

func (v RtpHeaderExtensionUri) String() string {
	if s, ok := EnumNamesRtpHeaderExtensionUri[v]; ok {
		return s
	}
	return "RtpHeaderExtensionUri(" + strconv.FormatInt(int64(v), 10) + ")"
}