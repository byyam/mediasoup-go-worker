package mediasoupdata

type ProducerOptions struct {
	/**
	 * Producer id (just for Router.pipeToRouter() method).
	 */
	Id string `json:"id,omitempty"`

	/**
	 * Media kind ('audio' or 'video').
	 */
	Kind MediaKind `json:"kind,omitempty"`

	/**
	 * RTP parameters defining what the endpoint is sending.
	 */
	RtpParameters RtpParameters `json:"rtpParameters,omitempty"`

	/**
	 * Whether the producer must start in paused mode. Default false.
	 */
	Paused bool `json:"paused,omitempty"`

	/**
	 * Just for video. Time (in ms) before asking the sender for a new key frame
	 * after having asked a previous one. Default 0.
	 */
	KeyFrameRequestDelay uint32 `json:"keyFrameRequestDelay,omitempty"`

	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`

	RtpMapping RtpMapping `json:"rtpMapping"`
}

func (o ProducerOptions) Valid() bool {
	if !o.RtpMapping.Valid() || !o.RtpParameters.Valid() {
		return false
	}
	return true
}

/**
 * Producer type.
 */
type ProducerType = RtpParametersType

const (
	ProducerType_Simple    ProducerType = RtpParametersType_Simple
	ProducerType_Simulcast              = RtpParametersType_Simulcast
	ProducerType_Svc                    = RtpParametersType_Svc
)

type ProducerData struct {
	Kind                    MediaKind     `json:"kind,omitempty"`
	Type                    ProducerType  `json:"type,omitempty"`
	RtpParameters           RtpParameters `json:"rtpParameters,omitempty"`
	ConsumableRtpParameters RtpParameters `json:"consumableRtpParameters,omitempty"`
}
