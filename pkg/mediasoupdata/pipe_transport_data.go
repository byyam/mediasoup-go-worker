package mediasoupdata

type PipeTransportOptions struct {
	/**
	 * Listening IP address.
	 */
	ListenIp TransportListenIp `json:"listenIp,omitempty"`
	Port     uint16            `json:"port,omitempty"`

	SctpOptions
	/**
	 * Enable RTX and NACK for RTP retransmission. Useful if both Routers are
	 * located in different hosts and there is packet lost in the link. For this
	 * to work, both PipeTransports must enable this setting. Default false.
	 */
	EnableRtx bool `json:"enableRtx,omitempty"`

	/**
	 * Enable SRTP. Useful to protect the RTP and RTCP traffic if both Routers
	 * are located in different hosts. For this to work, connect() must be called
	 * with remote SRTP parameters. Default false.
	 */
	EnableSrtp bool `json:"enableSrtp,omitempty"`

	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`
}

type PipeTransportData struct {
	Tuple          TransportTuple  `json:"tuple,omitempty"`
	SctpParameters SctpParameters  `json:"sctpParameters,omitempty"`
	SctpState      SctpState       `json:"sctpState,omitempty"`
	Rtx            bool            `json:"rtx,omitempty"`
	SrtpParameters *SrtpParameters `json:"srtpParameters,omitempty"`
}
