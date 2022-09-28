package mediasoupdata

/**
 * Consumer type.
 */
type ConsumerType = RtpParametersType

const (
	ConsumerType_Simple    ConsumerType = RtpParametersType_Simple
	ConsumerType_Simulcast              = RtpParametersType_Simulcast
	ConsumerType_Svc                    = RtpParametersType_Svc
)

type ConsumerOptions struct {
	/**
	 * The id of the Producer to consume.
	 */
	ProducerId string `json:"producerId,omitempty"`

	/**
	 * RTP capabilities of the consuming endpoint.
	 */
	RtpCapabilities RtpCapabilities `json:"rtpCapabilities,omitempty"`

	/**
	 * Whether the Consumer must start in paused mode. Default false.
	 *
	 * When creating a video Consumer, it's recommended to set paused to true,
	 * then transmit the Consumer parameters to the consuming endpoint and, once
	 * the consuming endpoint has created its local side Consumer, unpause the
	 * server side Consumer using the resume() method. This is an optimization
	 * to make it possible for the consuming endpoint to render the video as far
	 * as possible. If the server side Consumer was created with paused false,
	 * mediasoup will immediately request a key frame to the remote Producer and
	 * suych a key frame may reach the consuming endpoint even before it's ready
	 * to consume it, generating “black” video until the device requests a keyframe
	 * by itself.
	 */
	Paused bool `json:"paused,omitempty"`

	/**
	 * The MID for the Consumer. If not specified, a sequentially growing
	 * number will be assigned.
	 */
	Mid string `json:"mid,omitempty"`

	/**
	 * Preferred spatial and temporal layer for simulcast or SVC media sources.
	 * If unset, the highest ones are selected.
	 */
	PreferredLayers *ConsumerLayers `json:"preferredLayers,omitempty"`

	/**
	 * Whether this Consumer should consume all RTP streams generated by the
	 * Producer.
	 */
	Pipe bool

	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`

	// input
	Kind                   MediaKind               `json:"kind,omitempty"`
	Type                   ConsumerType            `json:"type,omitempty"`
	RtpParameters          RtpParameters           `json:"rtpParameters,omitempty"`
	ConsumableRtpEncodings []RtpEncodingParameters `json:"consumableRtpEncodings,omitempty"`
}

type ConsumerLayers struct {
	/**
	 * The spatial layer index (from 0 to N).
	 */
	SpatialLayer uint8 `json:"spatialLayer"`

	/**
	 * The temporal layer index (from 0 to N).
	 */
	TemporalLayer uint8 `json:"temporalLayer"`
}

type ConsumerScore struct {
	/**
	 * The score of the RTP stream of the consumer.
	 */
	Score uint16 `json:"score"`

	/**
	 * The score of the currently selected RTP stream of the producer.
	 */
	ProducerScore uint16 `json:"producerScore"`

	/**
	 * The scores of all RTP streams in the producer ordered by encoding (just
	 * useful when the producer uses simulcast).
	 */
	ProducerScores []uint16 `json:"producerScores,omitempty"`
}

type ConsumerData struct {
	Paused         bool          `json:"paused,omitempty"`
	ProducerPaused bool          `json:"producerPaused,omitempty"`
	Score          ConsumerScore `json:"score,omitempty"`
}

type ConsumerStat = ProducerStat