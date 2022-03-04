package mediasoupdata

type RouterOptions struct {
	/**
	 * Router media codecs.
	 */
	MediaCodecs []*RtpCodecCapability `json:"mediaCodecs,omitempty"`

	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`
}
