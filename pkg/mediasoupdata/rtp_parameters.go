package mediasoupdata

import (
	"encoding/json"
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
	ClockRate uint32 `json:"clockRate"`

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

const (
	AptString = "apt"
)

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
	Rtcp *RtcpParameters `json:"rtcp,omitempty"`
}

func (r *RtpParameters) Convert() *FBS__RtpParameters.RtpParametersT {
	p := &FBS__RtpParameters.RtpParametersT{
		Mid:              r.Mid,
		Codecs:           make([]*FBS__RtpParameters.RtpCodecParametersT, 0),
		HeaderExtensions: make([]*FBS__RtpParameters.RtpHeaderExtensionParametersT, 0),
		Encodings:        make([]*FBS__RtpParameters.RtpEncodingParametersT, 0),
		Rtcp:             r.Rtcp.Convert(),
	}
	for _, c := range r.Codecs {
		p.Codecs = append(p.Codecs, c.Convert())
	}
	for _, h := range r.HeaderExtensions {
		p.HeaderExtensions = append(p.HeaderExtensions, h.Convert())
	}
	for _, e := range r.Encodings {
		p.Encodings = append(p.Encodings, e.Convert())
	}
	return p
}

func (r *RtpParameters) Set(fbs *FBS__RtpParameters.RtpParametersT) {
	r.Mid = fbs.Mid
	r.Rtcp.Set(fbs.Rtcp)
	for _, f := range fbs.Codecs {
		c := &RtpCodecParameters{}
		c.Set(f)
		r.Codecs = append(r.Codecs, c)
	}
	for _, f := range fbs.HeaderExtensions {
		c := &RtpHeaderExtensionParameters{}
		c.Set(f)
		r.HeaderExtensions = append(r.HeaderExtensions, c)
	}
	for _, f := range fbs.Encodings {
		c := &RtpEncodingParameters{}
		c.Set(f)
		r.Encodings = append(r.Encodings, c)
	}
}

func NewRtpParameters(fbs *FBS__RtpParameters.RtpParametersT) *RtpParameters {
	r := &RtpParameters{
		Rtcp: &RtcpParameters{},
	}
	//_ = Clone(fbs, r)
	r.Set(fbs)

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

	// Validate RTP parameters.
	if err := r.validateCodecs(); err != nil {
		return err
	}
	if err := r.validateEncodings(); err != nil {
		return err
	}
	return nil
}

func (r *RtpParameters) validateCodecs() error {
	for _, codec := range r.Codecs {
		// todo: check duplicated payloadType
		codecParameters := &RtpCodecMimeType{}
		if err := codecParameters.SetMimeType(codec.MimeType); err != nil {
			return err
		}
		switch codecParameters.SubType {
		// A RTX codec must have 'apt' parameter pointing to a non RTX codec.
		case MimeSubTypeRTX:
			apt := GetIntegerByName(codec.Parameters.Convert(), AptString)
			for _, codec := range r.Codecs {
				if apt == int32(codec.PayloadType) {
					codecParameters := &RtpCodecMimeType{}
					if err := codecParameters.SetMimeType(codec.MimeType); err != nil {
						return err
					}
					if codecParameters.SubType == MimeSubTypeRTX {
						return errors.New("apt in RTX codec points to a RTX codec")
					} else if codecParameters.SubType == MimeSubTypeULPFEC {
						return errors.New("apt in RTX codec points to a ULPFEC codec")
					} else if codecParameters.SubType == MimeSubTypeFLEXFEC {
						return errors.New("apt in RTX codec points to a FLEXFEC codec")
					}
					return nil
				}
			}

		default:

		}
	}
	return nil
}

func GetIntegerByName(params []*FBS__RtpParameters.ParameterT, name string) int32 {
	for _, param := range params {
		if param.Name == name {
			return GetInteger(param)
		}
	}
	return 0
}

func GetInteger(param *FBS__RtpParameters.ParameterT) int32 {
	dataDump := &FBS__RtpParameters.Integer32T{}
	if err := Clone(param.Value.Value, dataDump); err != nil {
		return 0
	}
	return dataDump.Value
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

		if encoding.CodecPayloadType == 0 {
			encoding.CodecPayloadType = firstMediaPayloadType
		} else {
			var exist bool
			for _, codec := range r.Codecs {
				if codec.PayloadType == encoding.CodecPayloadType {
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
		if codec.PayloadType == payloadType {
			return codec
		}
	}
	panic("no valid codec payload type for the given encoding")
}

func (r *RtpParameters) GetRtxCodecForEncoding(encoding *RtpEncodingParameters) *RtpCodecParameters {
	payloadType := encoding.CodecPayloadType
	for _, codec := range r.Codecs {
		if codec.RtpCodecMimeType.IsFeatureCodec() && codec.Parameters.Apt == payloadType {
			return codec
		}
	}
	return nil
}

func (r *RtpParameters) CheckRTCPFeedbackType(name string) bool {
	for _, codec := range r.Codecs {
		for _, rtcpFeedback := range codec.RtcpFeedback {
			if rtcpFeedback.Type == name {
				return true
			}
		}
	}
	return false
}

/**
 * Provides information on codec settings within the RTP parameters. The list
 * of media codecs supported by mediasoup and their settings is defined in the
 * supportedRtpCapabilities.ts file.
 */
type RtpCodecParameters struct {
	/**
	 * The codec MIME media type/subtype (e.g. 'audio/opus', 'video/VP8').
	 */
	MimeType string `json:"mimeType"`

	/**
	 * The value that goes in the RTP Payload Type Field. Must be unique.
	 */
	PayloadType uint8 `json:"payloadType"`

	/**
	 * Codec clock rate expressed in Hertz.
	 */
	ClockRate uint32 `json:"clockRate"`

	/**
	 * The int of channels supported (e.g. two for stereo). Just for audio.
	 * Default 1.
	 */
	Channels uint8 `json:"channels,omitempty"`

	/**
	 * Codec-specific parameters available for signaling. Some parameters (such
	 * as 'packetization-mode' and 'profile-level-id' in H264 or 'profile-id' in
	 * VP9) are critical for codec matching.
	 */
	Parameters RtpCodecSpecificParameters `json:"parameters,omitempty"`

	/**
	 * Transport layer and codec-specific feedback messages for this codec.
	 */
	RtcpFeedback []*FBS__RtpParameters.RtcpFeedbackT `json:"rtcpFeedback,omitempty"`

	RtpCodecMimeType RtpCodecMimeType `json:"-"`
}

func (r *RtpCodecParameters) Convert() *FBS__RtpParameters.RtpCodecParametersT {
	p := &FBS__RtpParameters.RtpCodecParametersT{
		MimeType:     r.MimeType,
		PayloadType:  r.PayloadType,
		ClockRate:    r.ClockRate,
		Channels:     &r.Channels,
		Parameters:   r.Parameters.Convert(),
		RtcpFeedback: r.RtcpFeedback,
	}
	return p
}

func (r *RtpCodecParameters) Set(fbs *FBS__RtpParameters.RtpCodecParametersT) {
	r.MimeType = fbs.MimeType
	r.PayloadType = fbs.PayloadType
	r.ClockRate = fbs.ClockRate
	r.Channels = *fbs.Channels
	r.RtcpFeedback = fbs.RtcpFeedback
	r.Parameters.Set(fbs.Parameters)
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
	//r.setSpecificParameters()

	if err := r.CheckCodec(); err != nil {
		return err
	}
	return nil
}

//	func (r *RtpCodecParameters) setSpecificParameters() {
//		for _, p := range r.Parameters {
//			switch p.Name {
//			case AptString:
//				r.SpecificParameters.Apt = byte(GetInteger(p))
//			}
//		}
//	}
func (r *RtpCodecParameters) CheckCodec() error {
	switch r.RtpCodecMimeType.SubType {
	case MimeSubTypeRTX:
		if r.Parameters.Apt <= 0 {
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

func (r *RtpCodecSpecificParameters) Convert() []*FBS__RtpParameters.ParameterT {
	p := make([]*FBS__RtpParameters.ParameterT, 0)
	content, _ := json.Marshal(r)
	_ = json.Unmarshal(content, &p)
	return p
}

func (r *RtpCodecSpecificParameters) Set(fbs []*FBS__RtpParameters.ParameterT) {
	content, _ := json.Marshal(fbs)
	_ = json.Unmarshal(content, r)
}

/**
 * Provides information on RTCP feedback messages for a specific codec. Those
 * messages can be transport layer feedback messages or codec-specific feedback
 * messages. The list of RTCP feedbacks supported by mediasoup is defined in the
 * supportedRtpCapabilities.ts file.
 */
//type RtcpFeedback struct {
//	/**
//	 * RTCP feedback type.
//	 */
//	Type string `json:"type"`
//
//	/**
//	 * RTCP feedback parameter.
//	 */
//	Parameter string `json:"parameter,omitempty"`
//}
//
//func (r *RtcpFeedback) Convert() *FBS__RtpParameters.RtcpFeedbackT {
//	p := &FBS__RtpParameters.RtcpFeedbackT{
//		Type:      r.Type,
//		Parameter: r.Parameter,
//	}
//	return p
//}
//
/**
 * Provides information relating to an encoding, which represents a media RTP
 * stream and its associated RTX stream (if any).
 */
type RtpEncodingParameters struct {
	/**
	 * The media SSRC.
	 */
	Ssrc uint32 `json:"ssrc,omitempty"`

	/**
	 * The RID RTP extension value. Must be unique.
	 */
	Rid string `json:"rid,omitempty"`

	/**
	 * Codec payload type this encoding affects. If unset, first media codec is
	 * chosen.
	 */
	CodecPayloadType    byte `json:"codecPayloadType,omitempty"`
	HasCodecPayloadType bool `json:"-"`

	/**
	 * RTX stream information. It must contain a numeric ssrc field indicating
	 * the RTX SSRC.
	 */
	Rtx *FBS__RtpParameters.RtxT `json:"rtx,omitempty"`

	/**
	 * It indicates whether discontinuous RTP transmission will be used. Useful
	 * for audio (if the codec supports it) and for video screen sharing (when
	 * static content is being transmitted, this option disables the RTP
	 * inactivity checks in mediasoup). Default false.
	 */
	Dtx bool `json:"dtx,omitempty"`

	/**
	 * int of spatial and temporal layers in the RTP stream (e.g. 'L1T3').
	 * See webrtc-svc.
	 */
	ScalabilityMode string `json:"scalabilityMode,omitempty"`

	/**
	 * Others.
	 */
	ScaleResolutionDownBy int    `json:"scaleResolutionDownBy,omitempty"`
	MaxBitrate            uint32 `json:"maxBitrate,omitempty"`

	SpatialLayers  uint8 `json:"-"` // default 1
	TemporalLayers uint8 `json:"-"` // default 1
}

func (r *RtpEncodingParameters) Convert() *FBS__RtpParameters.RtpEncodingParametersT {
	p := &FBS__RtpParameters.RtpEncodingParametersT{
		Ssrc:             &r.Ssrc,
		Rid:              r.Rid,
		CodecPayloadType: &r.CodecPayloadType,
		Rtx:              r.Rtx,
		Dtx:              r.Dtx,
		ScalabilityMode:  r.ScalabilityMode,
		MaxBitrate:       &r.MaxBitrate,
	}
	return p
}

func (r *RtpEncodingParameters) Set(fbs *FBS__RtpParameters.RtpEncodingParametersT) {
	if fbs.Ssrc != nil {
		r.Ssrc = *fbs.Ssrc
	}
	r.Rid = fbs.Rid
	if fbs.CodecPayloadType != nil {
		r.CodecPayloadType = *fbs.CodecPayloadType
	}
	r.Rtx = fbs.Rtx
	r.Dtx = fbs.Dtx
	r.ScalabilityMode = fbs.ScalabilityMode
	if fbs.MaxBitrate != nil {
		r.MaxBitrate = *fbs.MaxBitrate
	}
}

//func (r *RtpEncodingParameters) Init() error {
//
//	mode := ParseScalabilityMode(r.ScalabilityMode)
//	r.ParsedScalabilityMode.SpatialLayers = mode.SpatialLayers
//	r.ParsedScalabilityMode.TemporalLayers = mode.TemporalLayers
//	r.ParsedScalabilityMode.Ksvc = mode.Ksvc
//
//	return nil
//}

//// RtpEncodingRtx represents the associated RTX stream for RTP stream.
//type RtpEncodingRtx struct {
//	Ssrc uint32 `json:"ssrc"`
//}

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

func (r *RtpHeaderExtensionParameters) Convert() *FBS__RtpParameters.RtpHeaderExtensionParametersT {
	p := &FBS__RtpParameters.RtpHeaderExtensionParametersT{
		Uri:        r.Uri,
		Id:         r.Id,
		Encrypt:    r.Encrypt,
		Parameters: r.Parameters,
	}
	return p
}

func (r *RtpHeaderExtensionParameters) Set(fbs *FBS__RtpParameters.RtpHeaderExtensionParametersT) {
	r.SpecificParameters.Set(fbs.Parameters)
	r.Uri = fbs.Uri
	r.Id = fbs.Id
	r.Encrypt = fbs.Encrypt
	r.Parameters = fbs.Parameters
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

func (r *RtcpParameters) Convert() *FBS__RtpParameters.RtcpParametersT {
	p := &FBS__RtpParameters.RtcpParametersT{
		Cname:       r.Cname,
		ReducedSize: *r.ReducedSize,
	}
	return p
}

func (r *RtcpParameters) Set(fbs *FBS__RtpParameters.RtcpParametersT) {
	r.Cname = fbs.Cname
	r.ReducedSize = &fbs.ReducedSize
}
