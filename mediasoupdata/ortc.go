package mediasoupdata

type RtpMapping struct {
	Codecs    []RtpMappingCodec    `json:"codecs,omitempty"`
	Encodings []RtpMappingEncoding `json:"encodings,omitempty"`
}

type RtpMappingCodec struct {
	PayloadType       byte `json:"payloadType"`
	MappedPayloadType byte `json:"mappedPayloadType"`
}

type RtpMappingEncoding struct {
	Ssrc            uint32 `json:"ssrc,omitempty"`
	Rid             string `json:"rid,omitempty"`
	ScalabilityMode string `json:"scalabilityMode,omitempty"`
	MappedSsrc      uint32 `json:"mappedSsrc"`
}
