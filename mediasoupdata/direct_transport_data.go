package mediasoupdata

type DirectTransportOptions struct {
	/**
	 * Maximum allowed size for direct messages sent from DataProducers.
	 * Default 262144.
	 */
	MaxMessageSize uint32 `json:"maxMessageSize,omitempty"`

	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`
}
