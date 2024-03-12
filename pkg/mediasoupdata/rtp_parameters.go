package mediasoupdata

import (
	"errors"
	"strings"

	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
	"github.com/byyam/mediasoup-go-worker/pkg/h264"
)

/**
 * The RTP capabilities define what mediasoup or an endpoint can receive at
 * media level.
 */
type RtpCapabilities struct {
	/**
	 * Supported media and RTX codecs.
	 */
	Codecs []*RtpCodecCapability `json:"codecs,omitempty"`

	/**
	 * Supported RTP header extensions.
	 */
	HeaderExtensions []*RtpHeaderExtension `json:"headerExtensions,omitempty"`

	/**
	 * Supported FEC mechanisms.
	 */
	FecMechanisms []string `json:"fecMechanisms,omitempty"`
}

/**
 * Media kind ('audio' or 'video').
 */
type MediaKind string

const (
	MediaKind_Audio MediaKind = "audio"
	MediaKind_Video MediaKind = "video"
)

/**
 * Provides information on the capabilities of a codec within the RTP
 * capabilities. The list of media codecs supported by mediasoup and their
 * settings is defined in the supportedRtpCapabilities.ts file.
 *
 * Exactly one RtpCodecCapability will be present for each supported combination
 * of parameters that requires a distinct value of preferredPayloadType. For
 * example
 *
 * - Multiple H264 codecs, each with their own distinct 'packetization-mode' and
 *   'profile-level-id' values.
 * - Multiple VP9 codecs, each with their own distinct 'profile-id' value.
 *
 * RtpCodecCapability entries in the mediaCodecs array of RouterOptions do not
 * require preferredPayloadType field (if unset, mediasoup will choose a random
 * one). If given, make sure it's in the 96-127 range.
 */
type RtpCodecCapability struct {
	/**
	 * Media kind.
	 */
	Kind MediaKind `json:"kind"`

	/**
	 * The codec MIME media type/subtype (e.g. 'audio/opus', 'video/VP8').
	 */
	MimeType string `json:"mimeType"`

	/**
	 * The preferred RTP payload type.
	 */
	PreferredPayloadType uint8 `json:"preferredPayloadType,omitempty"`

	/**
	 * Codec clock rate expressed in Hertz.
	 */
	ClockRate int `json:"clockRate"`

	/**
	 * The int of channels supported (e.g. two for stereo). Just for audio.
	 * Default 1.
	 */
	Channels uint8 `json:"channels,omitempty"`

	/**
	 * Codec specific parameters. Some parameters (such as 'packetization-mode'
	 * and 'profile-level-id' in H264 or 'profile-id' in VP9) are critical for
	 * codec matching.
	 */
	Parameters RtpCodecSpecificParameters `json:"parameters,omitempty"`

	/**
	 * Transport layer and codec-specific feedback messages for this codec.
	 */
	RtcpFeedback []*FBS__RtpParameters.RtcpFeedbackT `json:"rtcpFeedback,omitempty"`
}

func (r RtpCodecCapability) isRtxCodec() bool {
	return strings.HasSuffix(strings.ToLower(r.MimeType), "/rtx")
}

/**
 * Direction of RTP header extension.
 */
type RtpHeaderExtensionDirection string

const (
	Direction_Sendrecv RtpHeaderExtensionDirection = "sendrecv"
	Direction_Sendonly RtpHeaderExtensionDirection = "sendonly"
	Direction_Recvonly RtpHeaderExtensionDirection = "recvonly"
	Direction_Inactive RtpHeaderExtensionDirection = "inactive"
)

/**
 * Provides information relating to supported header extensions. The list of
 * RTP header extensions supported by mediasoup is defined in the
 * supportedRtpCapabilities.ts file.
 *
 * mediasoup does not currently support encrypted RTP header extensions. The
 * direction field is just present in mediasoup RTP capabilities (retrieved via
 * router.rtpCapabilities or mediasoup.getSupportedRtpCapabilities()). It's
 * ignored if present in endpoints' RTP capabilities.
 */
type RtpHeaderExtension struct {
	/**
	 * Media kind. If empty string, it's valid for all kinds.
	 * Default any media kind.
	 */
	Kind MediaKind `json:"kind"`

	/*
	 * The URI of the RTP header extension, as defined in RFC 5285.
	 */
	Uri string `json:"uri"`

	/**
	 * The preferred numeric identifier that goes in the RTP packet. Must be
	 * unique.
	 */
	PreferredId int `json:"preferredId"`

	/**
	 * If true, it is preferred that the value in the header be encrypted as per
	 * RFC 6904. Default false.
	 */
	PreferredEncrypt bool `json:"preferredEncrypt,omitempty"`

	/**
	 * If 'sendrecv', mediasoup supports sending and receiving this RTP extension.
	 * 'sendonly' means that mediasoup can send (but not receive) it. 'recvonly'
	 * means that mediasoup can receive (but not send) it.
	 */
	Direction RtpHeaderExtensionDirection `json:"direction,omitempty"`
}

/**
 * The RTP send parameters describe a media stream received by mediasoup from
 * an endpoint through its corresponding mediasoup Producer. These parameters
 * may include a mid value that the mediasoup transport will use to match
 * received RTP packets based on their MID RTP extension value.
 *
 * mediasoup allows RTP send parameters with a single encoding and with multiple
 * encodings (simulcast). In the latter case, each entry in the encodings array
 * must include a ssrc field or a rid field (the RID RTP extension value). Check
 * the Simulcast and SVC sections for more information.
 *
 * The RTP receive parameters describe a media stream as sent by mediasoup to
 * an endpoint through its corresponding mediasoup Consumer. The mid value is
 * unset (mediasoup does not include the MID RTP extension into RTP packets
 * being sent to endpoints).
 *
 * There is a single entry in the encodings array (even if the corresponding
 * producer uses simulcast). The consumer sends a single and continuous RTP
 * stream to the endpoint and spatial/temporal layer selection is possible via
 * consumer.setPreferredLayers().
 *
 * As an exception, previous bullet is not true when consuming a stream over a
 * PipeTransport, in which all RTP streams from the associated producer are
 * forwarded verbatim through the consumer.
 *
 * The RTP receive parameters will always have their ssrc values randomly
 * generated for all of its  encodings (and optional rtx { ssrc XXXX } if the
 * endpoint supports RTX), regardless of the original RTP send parameters in
 * the associated producer. This applies even if the producer's encodings have
 * rid set.
 */
type RtpParameters struct {
	/**
	 * The MID RTP extension value as defined in the BUNDLE specification.
	 */
	Mid string `json:"mid,omitempty"`

	/**
	 * Media and RTX codecs in use.
	 */
	Codecs []*RtpCodecParameters `json:"codecs"`

	/**
	 * RTP header extensions in use.
	 */
	HeaderExtensions []*RtpHeaderExtensionParameters `json:"header_extensions,omitempty"`

	/**
	 * Transmitted RTP streams and their settings.
	 */
	Encodings []*RtpEncodingParameters `json:"encodings,omitempty"`

	/**
	 * Parameters used for RTCP.
	 */
	Rtcp *FBS__RtpParameters.RtcpParametersT `json:"rtcp,omitempty"`
}

func NewRtpParameters(fbs *FBS__RtpParameters.RtpParametersT) *RtpParameters {
	r := &RtpParameters{}
	_ = Clone(fbs, r)

	return r
}

func (r *RtpParameters) Init() error {
	if !r.Valid() {
		return errors.New("valid check failed")
	}
	for _, codec := range r.Codecs {
		if err := codec.Init(); err != nil {
			return err
		}
	}

	// Validate RTP parameters.
	r.validateCodecs()
	if err := r.validateEncodings(); err != nil {
		return err
	}
	return nil
}

func (r *RtpParameters) validateCodecs() {
	// todo
}

func (r *RtpParameters) validateEncodings() error {
	firstMediaPayloadType := uint8(0)
	{
		var exist bool
		for _, codec := range r.Codecs {
			if codec.RtpCodecMimeType.IsMediaCodec() {
				firstMediaPayloadType = codec.PayloadType
				exist = true
			}
		}
		if !exist {
			return errors.New("no media codecs found")
		}
	}

	// Iterate all the encodings, set the first payloadType in all of them with
	// codecPayloadType unset, and check that others point to a media codec.
	//
	// Also, don't allow multiple SVC spatial layers into an encoding if there
	// are more than one encoding (simulcast).
	for _, encoding := range r.Encodings {
		if encoding.SpatialLayers == 0 {
			encoding.SpatialLayers = 1
		}
		if encoding.TemporalLayers == 0 {
			encoding.TemporalLayers = 1
		}
		if encoding.SpatialLayers > 1 && len(r.Encodings) > 1 {
			return errors.New("cannot use both simulcast and encodings with multiple SVC spatial layers")
		}

		if encoding.CodecPayloadType == nil {
			encoding.CodecPayloadType = &firstMediaPayloadType
		} else {
			var exist bool
			for _, codec := range r.Codecs {
				if codec.PayloadType == *encoding.CodecPayloadType {
					// Must be a media codec.
					if codec.RtpCodecMimeType.IsMediaCodec() {
						exist = true
						break
					}
					return errors.New("invalid codecPayloadType")
				}
			}
			if !exist {
				return errors.New("unknown codecPayloadType")
			}
		}
	}
	return nil
}

type RtpParametersType string

const (
	RtpParametersType_Simple    RtpParametersType = "simple"
	RtpParametersType_Simulcast                   = "simulcast"
	RtpParametersType_Svc                         = "svc"
	RtpParametersType_Pipe                        = "pipe"
	RtpParametersType_None                        = "none"
)

func (r *RtpParameters) Valid() bool {
	// encodings are mandatory.
	if r.Encodings == nil || len(r.Encodings) == 0 {
		return false
	}
	// codecs are mandatory.
	if r.Codecs == nil || len(r.Codecs) == 0 {
		return false
	}
	return true
}

func (r *RtpParameters) GetType() FBS__RtpParameters.Type {
	if len(r.Encodings) == 1 {
		if r.Encodings[0].SpatialLayers > 1 || r.Encodings[0].TemporalLayers > 1 {
			return FBS__RtpParameters.TypeSVC
		}
		return FBS__RtpParameters.TypeSIMPLE
	} else if len(r.Encodings) > 1 {
		return FBS__RtpParameters.TypeSIMULCAST
	}
	return FBS__RtpParameters.TypeSIMPLE
}

func (r *RtpParameters) GetCodecForEncoding(encoding *RtpEncodingParameters) *RtpCodecParameters {
	payloadType := encoding.CodecPayloadType
	for _, codec := range r.Codecs {
		if codec.PayloadType == *payloadType {
			return codec
		}
	}
	panic("no valid codec payload type for the given encoding")
}

func (r *RtpParameters) GetRtxCodecForEncoding(encoding *RtpEncodingParameters) *RtpCodecParameters {
	payloadType := encoding.CodecPayloadType
	for _, codec := range r.Codecs {
		if codec.RtpCodecMimeType.IsFeatureCodec() && codec.SpecificParameters.Apt == *payloadType {
			return codec
		}
	}
	return nil
}

/**
 * Provides information on codec settings within the RTP parameters. The list
 * of media codecs supported by mediasoup and their settings is defined in the
 * supportedRtpCapabilities.ts file.
 */
type RtpCodecParameters struct {
	FBS__RtpParameters.RtpCodecParametersT
	/**
	 * Codec-specific parameters available for signaling. Some parameters (such
	 * as 'packetization-mode' and 'profile-level-id' in H264 or 'profile-id' in
	 * VP9) are critical for codec matching.
	 */
	SpecificParameters RtpCodecSpecificParameters `json:"specific_parameters,omitempty"`

	RtpCodecMimeType RtpCodecMimeType `json:"-"`
}

func (r *RtpCodecParameters) Init() error {
	if err := r.RtpCodecMimeType.SetMimeType(r.MimeType); err != nil {
		return err
	}
	if r.PayloadType <= 0 {
		return errors.New("missing payloadType")
	}
	if r.ClockRate <= 0 {
		return errors.New("missing clockRate")
	}
	r.setSpecificParameters()

	if err := r.CheckCodec(); err != nil {
		return err
	}
	return nil
}

func (r *RtpCodecParameters) setSpecificParameters() {
	for _, p := range r.Parameters {
		switch p.Name {
		case "apt":
			r.SpecificParameters.Apt = 1
		}
	}
}

func (r *RtpCodecParameters) CheckCodec() error {
	switch r.RtpCodecMimeType.SubType {
	case MimeSubTypeRTX:
		if r.SpecificParameters.Apt <= 0 {
			return errors.New("missing apt parameter in RTX codec")
		}
	default:
		return nil
	}
	return nil
}

func (r *RtpCodecParameters) isRtxCodec() bool {
	return strings.HasSuffix(strings.ToLower(r.MimeType), "/rtx")
}

/**
 * RtpCodecSpecificParameters the Codec-specific parameters available for signaling. Some parameters (such
 * as 'packetization-mode' and 'profile-level-id' in H264 or 'profile-id' in
 * VP9) are critical for codec matching.
 */
type RtpCodecSpecificParameters struct {
	h264.RtpParameter          // used by h264 codec
	ProfileId           uint8  `json:"profile-id,omitempty"`   // used by vp9  https://www.webmproject.org/vp9/profiles/
	Apt                 byte   `json:"apt,omitempty"`          // used by rtx codec
	SpropStereo         uint8  `json:"sprop-stereo,omitempty"` // used by audio, 1 or 0
	Useinbandfec        uint8  `json:"useinbandfec,omitempty"` // used by audio, 1 or 0
	Usedtx              uint8  `json:"usedtx,omitempty"`       // used by audio, 1 or 0
	Maxplaybackrate     uint32 `json:"maxplaybackrate,omitempty"`
	XGoogleMinBitrate   uint32 `json:"x-google-min-bitrate,omitempty"`
	XGoogleMaxBitrate   uint32 `json:"x-google-max-bitrate,omitempty"`
	XGoogleStartBitrate uint32 `json:"x-google-start-bitrate,omitempty"`
	ChannelMapping      string `json:"channel_mapping,omitempty"`
	NumStreams          uint8  `json:"num_streams,omitempty"`
	CoupledStreams      uint8  `json:"coupled_streams,omitempty"`
}

/**
 * Provides information on RTCP feedback messages for a specific codec. Those
 * messages can be transport layer feedback messages or codec-specific feedback
 * messages. The list of RTCP feedbacks supported by mediasoup is defined in the
 * supportedRtpCapabilities.ts file.
 */
type RtcpFeedback struct {
	/**
	 * RTCP feedback type.
	 */
	Type string `json:"type"`

	/**
	 * RTCP feedback parameter.
	 */
	Parameter string `json:"parameter,omitempty"`
}

/**
 * Provides information relating to an encoding, which represents a media RTP
 * stream and its associated RTX stream (if any).
 */
type RtpEncodingParameters struct {
	FBS__RtpParameters.RtpEncodingParametersT

	SpatialLayers  uint8 `json:"-"` // default 1
	TemporalLayers uint8 `json:"-"` // default 1
}

// RtpEncodingRtx represents the associated RTX stream for RTP stream.
type RtpEncodingRtx struct {
	Ssrc uint32 `json:"ssrc"`
}

/**
 * Defines a RTP header extension within the RTP parameters. The list of RTP
 * header extensions supported by mediasoup is defined in the
 * supportedRtpCapabilities.ts file.
 *
 * mediasoup does not currently support encrypted RTP header extensions and no
 * parameters are currently considered.
 */
type RtpHeaderExtensionParameters struct {
	FBS__RtpParameters.RtpHeaderExtensionParametersT

	/**
	 * Configuration parameters for the header extension.
	 */
	SpecificParameters *RtpCodecSpecificParameters `json:"specific_parameters,omitempty"`
}

/**
 * Provides information on RTCP settings within the RTP parameters.
 *
 * If no cname is given in a producer's RTP parameters, the mediasoup transport
 * will choose a random one that will be used into RTCP SDES messages sent to
 * all its associated consumers.
 *
 * mediasoup assumes reducedSize to always be true.
 */
type RtcpParameters struct {
	/**
	 * The Canonical Name (CNAME) used by RTCP (e.g. in SDES messages).
	 */
	Cname string `json:"cname,omitempty"`

	/**
	 * Whether reduced size RTCP RFC 5506 is configured (if true) or compound RTCP
	 * as specified in RFC 3550 (if false). Default true.
	 */
	ReducedSize *bool `json:"reducedSize,omitempty"`

	/**
	 * Whether RTCP-mux is used. Default true.
	 */
	Mux *bool `json:"mux,omitempty"`
}
