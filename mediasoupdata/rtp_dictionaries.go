package mediasoupdata

import (
	"errors"
	"strings"
)

type MimeType int32

const (
	MimeTypeUnset MimeType = iota
	MimeTypeAudio
	MimeTypeVideo
)

type MimeSubType int32

const (
	MimeSubTypeAudioCodec         = 100
	MimeSubTypeVideoCodec         = 200
	MimeSubTypeComplementaryCodec = 300
	MimeSubTypeFeatureCodec       = 400
)

const (
	MimeSubTypeUNSET MimeSubType = iota
)

// Audio codecs:
const (
	MimeSubTypeOPUS MimeSubType = MimeSubTypeAudioCodec + iota
	// Multi-channel Opus.
	MimeSubTypeMULTIOPUS
	MimeSubTypePCMA
	MimeSubTypePCMU
	MimeSubTypeISAC
	MimeSubTypeG722
	MimeSubTypeILBC
	MimeSubTypeSILK
)

// Video codecs:
const (
	MimeSubTypeVP8 MimeSubType = MimeSubTypeVideoCodec + iota
	MimeSubTypeVP9
	MimeSubTypeH264
	MimeSubTypeX_H264UC
	MimeSubTypeH265
)

// Complementary codecs:
const (
	MimeSubTypeCN MimeSubType = MimeSubTypeComplementaryCodec + iota
	MimeSubTypeTELEPHONE_EVENT
)

// Feature codecs:
const (
	MimeSubTypeRTX MimeSubType = MimeSubTypeFeatureCodec + iota
	MimeSubTypeULPFEC
	MimeSubTypeX_ULPFECUC
	MimeSubTypeFLEXFEC
	MimeSubTypeRED
)

type RtpCodecMimeType struct {
	Type     MimeType
	SubType  MimeSubType
	MimeType string
}

var rtpCodecMimeType2String = map[MimeType]string{
	MimeTypeAudio: "audio",
	MimeTypeVideo: "video",
}

var rtpCodecMimeString2Type = map[string]MimeType{}

var rtpCodecMimeSubType2String = map[MimeSubType]string{
	// Audio codecs:
	MimeSubTypeOPUS:      "opus",
	MimeSubTypeMULTIOPUS: "multiopus",
	MimeSubTypePCMA:      "PCMA",
	MimeSubTypePCMU:      "PCMU",
	MimeSubTypeISAC:      "ISAC",
	MimeSubTypeG722:      "G722",
	MimeSubTypeILBC:      "iLBC",
	MimeSubTypeSILK:      "SILK",
	// Video codecs:
	MimeSubTypeVP8:      "VP8",
	MimeSubTypeVP9:      "VP9",
	MimeSubTypeH264:     "H264",
	MimeSubTypeX_H264UC: "X-H264UC",
	MimeSubTypeH265:     "H265",
	// Complementary codecs:
	MimeSubTypeCN:              "CN",
	MimeSubTypeTELEPHONE_EVENT: "telephone-event",
	// Feature codecs:
	MimeSubTypeRTX:        "rtx",
	MimeSubTypeULPFEC:     "ulpfec",
	MimeSubTypeFLEXFEC:    "flexfec",
	MimeSubTypeX_ULPFECUC: "x-ulpfecuc",
	MimeSubTypeRED:        "red",
}

var rtpCodecMimeString2SubType = map[string]MimeSubType{}

func (r *RtpCodecMimeType) SetMimeType(mimeType string) error {
	// Force lowcase names.
	// Set mimeType.
	r.MimeType = strings.ToLower(mimeType)

	slashPos := strings.Split(r.MimeType, "/")
	if len(slashPos) != 2 {
		return errors.New("wrong codec MIME")
	}
	typeStr := slashPos[0]
	subTypeStr := slashPos[1]

	var ok bool
	// Set MIME type.
	r.Type, ok = rtpCodecMimeString2Type[typeStr]
	if !ok {
		return errors.New("unknown codec MIME type")
	}
	// Set MIME subtype.
	r.SubType, ok = rtpCodecMimeString2SubType[subTypeStr]
	if !ok {
		return errors.New("unknown codec MIME subtype")
	}
	return nil
}

func (r *RtpCodecMimeType) UpdateMimeType() {
	if r.Type == MimeTypeUnset {
		panic("type unset")
	}
	if r.SubType == MimeSubTypeUNSET {
		panic("subtype unset")
	}

	// Set mimeType.
	r.MimeType = r.Type2String() + "/" + r.SubType2String()
}

func (r RtpCodecMimeType) IsMediaCodec() bool {
	if r.SubType >= MimeSubTypeAudioCodec && r.SubType < MimeSubTypeComplementaryCodec {
		return true
	}
	return false
}

func (r RtpCodecMimeType) IsComplementaryCodec() bool {
	if r.SubType >= MimeSubTypeComplementaryCodec && r.SubType < MimeSubTypeFeatureCodec {
		return true
	}
	return false
}

func (r RtpCodecMimeType) IsFeatureCodec() bool {
	if r.SubType >= MimeSubTypeFeatureCodec {
		return true
	}
	return false
}

func (r RtpCodecMimeType) Type2String() string {
	return rtpCodecMimeType2String[r.Type]
}

func (r RtpCodecMimeType) SubType2String() string {
	return rtpCodecMimeSubType2String[r.SubType]
}
