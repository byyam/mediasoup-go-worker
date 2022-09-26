package mediasoupdata

type WebRtcTransportOptions struct {
	/**
	 * Listening IP address or addresses in order of preference (first one is the
	 * preferred one).
	 */
	ListenIps []TransportListenIp `json:"listenIps,omitempty"`

	/**
	 * Listen in UDP. Default true.
	 */
	EnableUdp *bool `json:"enableUdp,omitempty"`

	/**
	 * Listen in TCP. Default false.
	 */
	EnableTcp bool `json:"enableTcp,omitempty"`

	/**
	 * Prefer UDP. Default false.
	 */
	PreferUdp bool `json:"preferUdp,omitempty"`

	/**
	 * Prefer TCP. Default false.
	 */
	PreferTcp bool `json:"preferTcp,omitempty"`

	/**
	 * Initial available outgoing bitrate (in bps). Default 600000.
	 */
	InitialAvailableOutgoingBitrate uint32 `json:"initialAvailableOutgoingBitrate,omitempty"`

	/**
	 * Create a SCTP association. Default false.
	 */
	EnableSctp bool `json:"enableSctp,omitempty"`

	/**
	 * SCTP streams uint32.
	 */
	NumSctpStreams NumSctpStreams `json:"numSctpStreams,omitempty"`

	/**
	 * Maximum allowed size for SCTP messages sent by DataProducers.
	 * Default 262144.
	 */
	MaxSctpMessageSize int `json:"maxSctpMessageSize,omitempty"`

	/**
	 * Maximum SCTP send buffer used by DataConsumers.
	 * Default 262144.
	 */
	SctpSendBufferSize int `json:"sctpSendBufferSize,omitempty"`

	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`
}

type WebrtcTransportData struct {
	// always be 'controlled'
	IceRole          string         `json:"iceRole,omitempty"`
	IceParameters    IceParameters  `json:"iceParameters,omitempty"`
	IceCandidates    []IceCandidate `json:"iceCandidates,omitempty"`
	IceState         IceState       `json:"iceState,omitempty"`
	IceSelectedTuple TransportTuple `json:"iceSelectedTuple,omitempty"`
	DtlsParameters   DtlsParameters `json:"dtlsParameters,omitempty"`
	DtlsState        DtlsState      `json:"dtlsState,omitempty"`
	DtlsRemoteCert   string         `json:"dtlsRemoteCert,omitempty"`
	SctpParameters   SctpParameters `json:"sctpParameters,omitempty"`
	SctpState        SctpState      `json:"sctpState,omitempty"`
}

type IceParameters struct {
	UsernameFragment string `json:"usernameFragment"`
	Password         string `json:"password"`
	IceLite          bool   `json:"iceLite,omitempty"`
}

type IceCandidate struct {
	Foundation string            `json:"foundation"`
	Priority   uint32            `json:"priority"`
	Ip         string            `json:"ip"`
	Protocol   TransportProtocol `json:"protocol"`
	Port       uint32            `json:"port"`
	// always "host"
	Type string `json:"type,omitempty"`
	// "passive" | undefined
	TcpType string `json:"tcpType,omitempty"`
}

type DtlsParameters struct {
	Role         DtlsRole          `json:"role,omitempty"`
	Fingerprints []DtlsFingerprint `json:"fingerprints"`
}

/**
 * The hash function algorithm (as defined in the "Hash function Textual Names"
 * registry initially specified in RFC 4572 Section 8) and its corresponding
 * certificate fingerprint value (in lowercase hex string as expressed utilizing
 * the syntax of "fingerprint" in RFC 4572 Section 5).
 */
type DtlsFingerprint struct {
	Algorithm string `json:"algorithm"`
	Value     string `json:"value"`
}

type IceState string

const (
	IceState_New          IceState = "new"
	IceState_Connected    IceState = "connected"
	IceState_Completed    IceState = "completed"
	IceState_Disconnected IceState = "disconnected"
	IceState_Closed       IceState = "closed"
)

type DtlsRole string

const (
	DtlsRole_Auto   DtlsRole = "auto"
	DtlsRole_Client DtlsRole = "client"
	DtlsRole_Server DtlsRole = "server"
)

type DtlsState string

const (
	DtlsState_New        = "new"
	DtlsState_Connecting = "connecting"
	DtlsState_Connected  = "connected"
	DtlsState_Failed     = "failed"
	DtlsState_Closed     = "closed"
)

type WebRtcTransportSpecificStat struct {
	IceRole          string          `json:"iceRole"`
	IceState         IceState        `json:"iceState"`
	DtlsState        DtlsRole        `json:"dtlsState"`
	IceSelectedTuple *TransportTuple `json:"iceSelectedTuple,omitempty"`
}
