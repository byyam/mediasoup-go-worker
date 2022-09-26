package mediasoupdata

type RtpMapping struct {
	Codecs    []RtpMappingCodec    `json:"codecs,omitempty"`
	Encodings []RtpMappingEncoding `json:"encodings,omitempty"`
}

func (r RtpMapping) Valid() bool {
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

type RtpMappingCodec struct {
	PayloadType       byte `json:"payloadType"`
	MappedPayloadType byte `json:"mappedPayloadType"`
}

func (r RtpMappingCodec) Valid() bool {
	if r.MappedPayloadType == 0 || r.PayloadType == 0 {
		return false
	}
	return true
}

type RtpMappingEncoding struct {
	Ssrc            uint32 `json:"ssrc,omitempty"`
	Rid             string `json:"rid,omitempty"`
	ScalabilityMode string `json:"scalabilityMode,omitempty"`
	MappedSsrc      uint32 `json:"mappedSsrc"`
}

func (r RtpMappingEncoding) Valid() bool {
	return true
}
