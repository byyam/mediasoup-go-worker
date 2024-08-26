package mediasoupdata

import (
	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
)

type RtpMapping struct {
	Codecs    []*RtpMappingCodec    `json:"codecs,omitempty"`
	Encodings []*RtpMappingEncoding `json:"encodings,omitempty"`
}

func (r *RtpMapping) Valid() bool {
	if len(r.Codecs) == 0 || len(r.Encodings) == 0 {
		return false
	}
	for _, codec := range r.Codecs {
		if !codec.Valid() {
			return false
		}
	}
	for _, encoding := range r.Encodings {
		if !encoding.Valid() {
			return false
		}
	}
	return true
}

func (r *RtpMapping) Convert() *FBS__RtpParameters.RtpMappingT {
	p := &FBS__RtpParameters.RtpMappingT{
		Codecs:    make([]*FBS__RtpParameters.CodecMappingT, 0),
		Encodings: make([]*FBS__RtpParameters.EncodingMappingT, 0),
	}

	for _, codec := range r.Codecs {
		c := codec.Convert()
		p.Codecs = append(p.Codecs, c)
	}

	for _, encoding := range r.Encodings {
		e := encoding.Convert()
		p.Encodings = append(p.Encodings, e)
	}
	return p
}

func (r *RtpMapping) Set(fbs *FBS__RtpParameters.RtpMappingT) {
	if r.Codecs == nil {
		r.Codecs = make([]*RtpMappingCodec, 0)
	}
	if r.Encodings == nil {
		r.Encodings = make([]*RtpMappingEncoding, 0)
	}
	for _, codec := range fbs.Codecs {
		c := &RtpMappingCodec{}
		c.Set(codec)
		r.Codecs = append(r.Codecs, c)
	}
	for _, encoding := range fbs.Encodings {
		e := &RtpMappingEncoding{}
		e.Set(encoding)
		r.Encodings = append(r.Encodings, e)
	}
}

type RtpMappingCodec struct {
	PayloadType       byte `json:"payloadType"`
	MappedPayloadType byte `json:"mappedPayloadType"`
}

func (r *RtpMappingCodec) Valid() bool {
	if r.MappedPayloadType == 0 || r.PayloadType == 0 {
		return false
	}
	return true
}

func (r *RtpMappingCodec) Convert() *FBS__RtpParameters.CodecMappingT {
	p := &FBS__RtpParameters.CodecMappingT{
		PayloadType:       r.PayloadType,
		MappedPayloadType: r.MappedPayloadType,
	}
	return p
}

func (r *RtpMappingCodec) Set(fbs *FBS__RtpParameters.CodecMappingT) {
	r.PayloadType = fbs.PayloadType
	r.MappedPayloadType = fbs.MappedPayloadType
}

type RtpMappingEncoding struct {
	Ssrc            uint32 `json:"ssrc,omitempty"`
	Rid             string `json:"rid,omitempty"`
	ScalabilityMode string `json:"scalabilityMode,omitempty"`
	MappedSsrc      uint32 `json:"mappedSsrc"`
}

func (r *RtpMappingEncoding) Valid() bool {
	return true
}

func (r *RtpMappingEncoding) Convert() *FBS__RtpParameters.EncodingMappingT {
	p := &FBS__RtpParameters.EncodingMappingT{
		Rid:             r.Rid,
		ScalabilityMode: r.ScalabilityMode,
		MappedSsrc:      r.MappedSsrc,
	}
	if r.Ssrc != 0 {
		p.Ssrc = &r.Ssrc
	}
	return p
}

func (r *RtpMappingEncoding) Set(fbs *FBS__RtpParameters.EncodingMappingT) {
	r.ScalabilityMode = fbs.ScalabilityMode
	r.MappedSsrc = fbs.MappedSsrc
	if fbs.Ssrc != nil {
		r.Ssrc = *fbs.Ssrc
	}
	r.Rid = fbs.Rid
}
