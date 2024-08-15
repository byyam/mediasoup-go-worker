package mediasoupdata

import (
	FBS__SctpAssociation "github.com/byyam/mediasoup-go-worker/fbs/FBS/SctpAssociation"
	FBS__SctpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/SctpParameters"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
)

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

	SctpOptions
	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`
}

type WebrtcTransportData struct {
	// always be 'controlled'
	IceParameters    *FBS__WebRtcTransport.IceParametersT  `json:"iceParameters"`
	IceCandidates    []*FBS__WebRtcTransport.IceCandidateT `json:"iceCandidates"`
	DtlsParameters   *FBS__WebRtcTransport.DtlsParametersT `json:"dtlsParameters"`
	SctpParameters   *FBS__SctpParameters.SctpParametersT  `json:"sctpParameters"`
	IceRole          FBS__WebRtcTransport.IceRole          `json:"iceRole,omitempty"`
	IceState         FBS__WebRtcTransport.IceState         `json:"iceState,omitempty"`
	IceSelectedTuple *FBS__Transport.TupleT                `json:"iceSelectedTuple,omitempty"`
	DtlsState        FBS__WebRtcTransport.DtlsState        `json:"dtlsState,omitempty"`
	DtlsRemoteCert   string                                `json:"dtlsRemoteCert,omitempty"`
	SctpState        *FBS__SctpAssociation.SctpState       `json:"sctpState,omitempty"`
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
