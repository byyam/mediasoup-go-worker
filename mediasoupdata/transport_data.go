package mediasoupdata

type TransportListenIp struct {
	/**
	 * Listening IPv4 or IPv6.
	 */
	Ip string `json:"ip,omitempty"`

	/**
	 * Announced IPv4 or IPv6 (useful when running mediasoup behind NAT with
	 * private IP).
	 */
	AnnouncedIp string `json:"announcedIp,omitempty"`
}

/**
 * Transport protocol.
 */
type TransportProtocol string

const (
	TransportProtocol_Udp TransportProtocol = "udp"
	TransportProtocol_Tcp TransportProtocol = "tcp"
)

type TransportTraceEventType string

const (
	TransportTraceEventType_Probation TransportTraceEventType = "probation"
	TransportTraceEventType_Bwe       TransportTraceEventType = "bwe"
)

type TransportTuple struct {
	LocalIp    string `json:"localIp,omitempty"`
	LocalPort  uint16 `json:"localPort,omitempty"`
	RemoteIp   string `json:"remoteIp,omitempty"`
	RemotePort uint16 `json:"remotePort,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
}

type SctpState string

const (
	SctpState_New        = "new"
	SctpState_Connecting = "connecting"
	SctpState_Connected  = "connected"
	SctpState_Failed     = "failed"
	SctpState_Closed     = "closed"
)

type TransportConnectOptions struct {
	// pipe and plain transport
	Ip             string          `json:"ip,omitempty"`
	Port           uint16          `json:"port,omitempty"`
	SrtpParameters *SrtpParameters `json:"srtpParameters,omitempty"`

	// plain transport
	RtcpPort uint16 `json:"rtcpPort,omitempty"`

	// webrtc transport
	DtlsParameters *DtlsParameters `json:"dtlsParameters,omitempty"`
}

type TransportConnectData struct {
	DtlsLocalRole DtlsRole `json:"dtlsLocalRole,omitempty"`
}

type TransportType string

const (
	TransportType_Direct TransportType = "DirectTransport"
	TransportType_Plain  TransportType = "PlainTransport"
	TransportType_Pipe   TransportType = "PipeTransport"
	TransportType_Webrtc TransportType = "WebrtcTransport"
)
