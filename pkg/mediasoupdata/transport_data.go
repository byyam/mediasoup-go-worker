package mediasoupdata

import (
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
)

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

type TransportOptions struct {
	DirectTransportOptions
	SctpOptions
	InitialAvailableOutgoingBitrate uint32 `json:"initialAvailableOutgoingBitrate,omitempty"`
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
	// webrtc transport
	DtlsLocalRole FBS__WebRtcTransport.DtlsRole `json:"dtlsLocalRole,omitempty"`
	// pipe transport
	Tuple TransportTuple `json:"tuple,omitempty"`
}

type TransportType string

const (
	TransportType_Direct TransportType = "DirectTransport"
	TransportType_Plain  TransportType = "PlainTransport"
	TransportType_Pipe   TransportType = "PipeTransport"
	TransportType_Webrtc TransportType = "WebrtcTransport"
)

type TransportStat struct {
	// Common to all Transports.
	Type                     string    `json:"type,omitempty"`
	TransportId              string    `json:"transportId,omitempty"`
	Timestamp                int64     `json:"timestamp,omitempty"`
	SctpState                SctpState `json:"sctpState,omitempty"`
	BytesReceived            int64     `json:"bytesReceived,omitempty"`
	RecvBitrate              int64     `json:"recvBitrate,omitempty"`
	BytesSent                int64     `json:"bytesSent,omitempty"`
	SendBitrate              int64     `json:"sendBitrate,omitempty"`
	RtpBytesReceived         int64     `json:"rtpBytesReceived,omitempty"`
	RtpRecvBitrate           int64     `json:"rtpRecvBitrate,omitempty"`
	RtpBytesSent             int64     `json:"rtpBytesSent,omitempty"`
	RtpSendBitrate           int64     `json:"rtpSendBitrate,omitempty"`
	RtxBytesReceived         int64     `json:"rtxBytesReceived,omitempty"`
	RtxRecvBitrate           int64     `json:"rtxRecvBitrate,omitempty"`
	RtxBytesSent             int64     `json:"rtxBytesSent,omitempty"`
	RtxSendBitrate           int64     `json:"rtxSendBitrate,omitempty"`
	ProbationBytesSent       int64     `json:"probationBytesSent,omitempty"`
	ProbationSendBitrate     int64     `json:"probationSendBitrate,omitempty"`
	AvailableOutgoingBitrate int64     `json:"availableOutgoingBitrate,omitempty"`
	AvailableIncomingBitrate int64     `json:"availableIncomingBitrate,omitempty"`
	MaxIncomingBitrate       int64     `json:"maxIncomingBitrate,omitempty"`
	RtpPacketLossReceived    float64   `json:"rtpPacketLossReceived,omitempty"`
	RtpPacketLossSent        float64   `json:"rtpPacketLossSent,omitempty"`

	*WebRtcTransportSpecificStat
}
