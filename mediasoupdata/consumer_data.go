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
