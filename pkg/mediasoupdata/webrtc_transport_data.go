package mediasoupdata

import (
	"crypto"
	"regexp"
	"strings"

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

func (c *IceParameters) Set(fbs *FBS__WebRtcTransport.IceParametersT) {
	c.UsernameFragment = fbs.UsernameFragment
	c.Password = fbs.Password
	c.IceLite = fbs.IceLite
}

type IceCandidate struct {
	Foundation string            `json:"foundation"`
	Priority   uint32            `json:"priority"`
	Ip         string            `json:"ip"`
	Protocol   TransportProtocol `json:"protocol"`
	Port       uint16            `json:"port"`
	// always "host"
	Type string `json:"type,omitempty"`
	// "passive" | undefined
	TcpType string `json:"tcpType,omitempty"`
}

func (c *IceCandidate) Set(fbs *FBS__WebRtcTransport.IceCandidateT) {
	c.Foundation = fbs.Foundation
	c.Priority = fbs.Priority
	c.Ip = fbs.Ip
	c.Protocol = TransportProtocol(strings.ToLower(FBS__Transport.EnumNamesProtocol[fbs.Protocol]))
	c.Port = fbs.Port
	c.Type = strings.ToLower(FBS__WebRtcTransport.EnumNamesIceCandidateType[fbs.Type])
	if fbs.TcpType != nil {
		c.TcpType = strings.ToLower(FBS__WebRtcTransport.EnumNamesIceCandidateTcpType[*fbs.TcpType])
	}
}

type DtlsParameters struct {
	Role         string            `json:"role,omitempty"`
	Fingerprints []DtlsFingerprint `json:"fingerprints"`
}

func (c *DtlsParameters) Convert() *FBS__WebRtcTransport.DtlsParametersT {
	d := new(FBS__WebRtcTransport.DtlsParametersT)
	switch strings.ToLower(c.Role) {
	case DtlsRole_Client:
		d.Role = FBS__WebRtcTransport.DtlsRoleCLIENT
	case DtlsRole_Server:
		d.Role = FBS__WebRtcTransport.DtlsRoleSERVER
	case DtlsRole_Auto:
		d.Role = FBS__WebRtcTransport.DtlsRoleAUTO
	}
	re := regexp.MustCompile("-")
	d.Fingerprints = make([]*FBS__WebRtcTransport.FingerprintT, 0)
	for _, f := range c.Fingerprints {
		cleanedHash := re.ReplaceAllString(strings.ToUpper(f.Algorithm), "")
		d.Fingerprints = append(d.Fingerprints, &FBS__WebRtcTransport.FingerprintT{
			Algorithm: FBS__WebRtcTransport.EnumValuesFingerprintAlgorithm[cleanedHash],
			Value:     f.Value,
		})
	}
	return d
}

func (c *DtlsParameters) Set(fbs *FBS__WebRtcTransport.DtlsParametersT) {
	c.Role = strings.ToLower(FBS__WebRtcTransport.EnumNamesDtlsRole[fbs.Role])
	c.Fingerprints = make([]DtlsFingerprint, 0)
	for _, f := range fbs.Fingerprints {
		var hash crypto.Hash
		switch f.Algorithm {
		case FBS__WebRtcTransport.FingerprintAlgorithmSHA1:
			hash = crypto.SHA1
		case FBS__WebRtcTransport.FingerprintAlgorithmSHA224:
			hash = crypto.SHA224
		case FBS__WebRtcTransport.FingerprintAlgorithmSHA256:
			hash = crypto.SHA256
		case FBS__WebRtcTransport.FingerprintAlgorithmSHA384:
			hash = crypto.SHA384
		case FBS__WebRtcTransport.FingerprintAlgorithmSHA512:
			hash = crypto.SHA512
		}
		c.Fingerprints = append(c.Fingerprints, DtlsFingerprint{
			Algorithm: strings.ToLower(hash.String()),
			Value:     f.Value,
		})
	}
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

const (
	DtlsRole_Auto   string = "auto"
	DtlsRole_Client string = "client"
	DtlsRole_Server string = "server"
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
	DtlsState        string          `json:"dtlsState"`
	IceSelectedTuple *TransportTuple `json:"iceSelectedTuple,omitempty"`
}
