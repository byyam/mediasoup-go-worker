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
	TypeStr  string
	SubType  MimeSubType
	MimeType string
}

func (r *RtpCodecMimeType) SetMimeType(mimeType string) error {
	// Force lowcase names.
	r.MimeType = strings.ToLower(mimeType)

	slashPos := strings.Split(r.MimeType, "/")
	if len(slashPos) != 2 {
		return errors.New("wrong codec MIME")
	}
	r.TypeStr = slashPos[0]
	subType := slashPos[1]

	// Set MIME type.
	switch r.TypeStr {
	case "audio":
		r.Type = MimeTypeAudio
	case "video":
		r.Type = MimeTypeVideo
	default:
		return errors.New("unknown codec MIME type")
	}

	switch subType {
	// Audio codecs:
	case "opus":
		r.SubType = MimeSubTypeOPUS
	case "multiopus":
		r.SubType = MimeSubTypeMULTIOPUS
	case "pcma":
		r.SubType = MimeSubTypePCMA
	case "pcmu":
		r.SubType = MimeSubTypePCMU
	case "isac":
		r.SubType = MimeSubTypeISAC
	case "g722":
		r.SubType = MimeSubTypeG722
	case "ilbc":
		r.SubType = MimeSubTypeILBC
	case "silk":
		r.SubType = MimeSubTypeSILK
	// Video codecs:
	case "vp8":
		r.SubType = MimeSubTypeVP8
	case "vp9":
		r.SubType = MimeSubTypeVP9
	case "h264":
		r.SubType = MimeSubTypeH264
	case "h265":
		r.SubType = MimeSubTypeH265
	// Complementary codecs:
	case "cn":
		r.SubType = MimeSubTypeCN
	case "telephone-event":
		r.SubType = MimeSubTypeTELEPHONE_EVENT
	// Feature codecs:
	case "rtx":
		r.SubType = MimeSubTypeRTX
	case "ulpfec":
		r.SubType = MimeSubTypeULPFEC
	case "flexfec":
		r.SubType = MimeSubTypeFLEXFEC
	case "x-ulpfecuc":
		r.SubType = MimeSubTypeX_ULPFECUC
	case "red":
		r.SubType = MimeSubTypeRED
	default:
		return errors.New("unknown codec MIME subtype")
	}

	return nil
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
