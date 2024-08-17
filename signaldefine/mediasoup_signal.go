package signaldefine

import (
	"strings"

	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
)

const (
	MethodGetRouterRtpCapabilities = "getRouterRtpCapabilities"
	MethodJoin                     = "join"
	MethodCreateWebRtcTransport    = "createWebRtcTransport"
	MethodConnectWebRtcTransport   = "connectWebRtcTransport"
	MethodProduce                  = "produce"
	MethodRestartIce               = "restartIce"
	MethodCloseProducer            = "closeProducer"
	MethodGetTransportStats        = "getTransportStats"
)

type ClientDevice struct {
	SdkVersion  string `json:"sdkVersion"`
	Platform    int    `json:"platform"`
	OsVersion   string `json:"osVersion"`
	DeviceModel string `json:"deviceModel"`
	AppId       int    `json:"appId"`
	SceneId     int    `json:"sceneId"`
	DeviceId    string `json:"deviceId"`
}

type GetRouterRtpCapabilitiesRequest struct {
	Device ClientDevice `json:"device"`
}

type GetRouterRtpCapabilitiesResponse struct {
	mediasoupdata.RtpCapabilities
}

type JoinRequest struct {
	DisplayName      string
	Device           map[string]interface{}
	RtpCapabilities  *mediasoupdata.RtpCapabilities
	SctpCapabilities *mediasoupdata.SctpCapabilities
	PublishAudio     *bool `json:"publishAudio,omitempty"`
	PublishVideo     *bool `json:"publishVideo,omitempty"`
}

type JoinedPeerInfo struct {
	Id          string                 `json:"id"`
	DisplayName string                 `json:"displayName"`
	Device      map[string]interface{} `json:"device"`
}

type JoinResponse struct {
	Peers []JoinedPeerInfo `json:"peers"`
}

type CreateWebRtcTransportRequest struct {
	ForceTcp         bool
	Producing        bool
	Consuming        bool
	SctpCapabilities *mediasoupdata.SctpCapabilities
}

type CreateWebRtcTransportResponse struct {
	Id string `json:"id"`

	IceParameters  *mediasoupdata.IceParameters  `json:"iceParameters"`
	IceCandidates  []*mediasoupdata.IceCandidate `json:"iceCandidates"`
	DtlsParameters *mediasoupdata.DtlsParameters `json:"dtlsParameters"`
	SctpParameters *mediasoupdata.SctpParameters `json:"sctpParameters"`
}

type ConnectWebRtcTransportRequest struct {
	TransportId    string
	DtlsParameters mediasoupdata.DtlsParameters
}

type ConnectWebRtcTransportResponse struct {
	DtlsRole string `json:"dtlsRole"`
}

type ProduceRequest struct {
	TransportId   string
	Kind          mediasoupdata.MediaKind
	RtpParameters mediasoupdata.RtpParameters
	AppData       map[string]interface{} `json:"appData,omitempty"`
}

type ProduceResponse struct {
	Id string `json:"id"`
}

type TransportStats struct {
	TransportId string `json:"transportId"`
	Timestamp   uint64 `json:"timestamp"`
	//SctpState                *FBS__SctpAssociation.SctpState `json:"sctp_state"`
	BytesReceived            uint64   `json:"bytesReceived"`
	RecvBitrate              uint32   `json:"recvBitrate"`
	BytesSent                uint64   `json:"bytesSent"`
	SendBitrate              uint32   `json:"sendBitrate"`
	RtpBytesReceived         uint64   `json:"rtpBytesReceived"`
	RtpRecvBitrate           uint32   `json:"rtpRecvBitrate"`
	RtpBytesSent             uint64   `json:"rtpBytesSent"`
	RtpSendBitrate           uint32   `json:"rtpSendBitrate"`
	RtxBytesReceived         uint64   `json:"rtxBytesReceived"`
	RtxRecvBitrate           uint32   `json:"rtxRecvBitrate"`
	RtxBytesSent             uint64   `json:"rtxBytesSent"`
	RtxSendBitrate           uint32   `json:"rtxSendBitrate"`
	ProbationBytesSent       uint64   `json:"probationBytesSent"`
	ProbationSendBitrate     uint32   `json:"probationSendBitrate"`
	AvailableOutgoingBitrate *uint32  `json:"availableOutgoingBitrate,omitempty"`
	AvailableIncomingBitrate *uint32  `json:"availableIncomingBitrate,omitempty"`
	MaxIncomingBitrate       *uint32  `json:"maxIncomingBitrate,omitempty"`
	MaxOutgoingBitrate       *uint32  `json:"maxOutgoingBitrate,omitempty"`
	MinOutgoingBitrate       *uint32  `json:"minOutgoingBitrate,omitempty"`
	RtpPacketLossReceived    *float64 `json:"rtpPacketLossReceived,omitempty"`
	RtpPacketLossSent        *float64 `json:"rtpPacketLossSent,omitempty"`
}

func (t *TransportStats) Set(fbs *FBS__Transport.StatsT) {
	t.TransportId = fbs.TransportId
	t.Timestamp = fbs.Timestamp
	t.BytesReceived = fbs.BytesReceived
	t.RecvBitrate = fbs.RecvBitrate
	t.BytesSent = fbs.BytesSent
	t.SendBitrate = fbs.SendBitrate
	t.RtpBytesReceived = fbs.RtpBytesReceived
	t.RtpRecvBitrate = fbs.RtpRecvBitrate
	t.RtpBytesSent = fbs.RtpBytesSent
	t.RtpSendBitrate = fbs.RtpSendBitrate
	t.RtxBytesReceived = fbs.RtxBytesReceived
	t.RtxRecvBitrate = fbs.RtxRecvBitrate
	t.RtxBytesSent = fbs.RtxBytesSent
	t.RtxSendBitrate = fbs.RtxSendBitrate
	t.ProbationBytesSent = fbs.ProbationBytesSent
	t.ProbationSendBitrate = fbs.ProbationSendBitrate
	t.AvailableOutgoingBitrate = fbs.AvailableOutgoingBitrate
	t.AvailableIncomingBitrate = fbs.AvailableIncomingBitrate
	t.MaxIncomingBitrate = fbs.MaxIncomingBitrate
	t.MaxOutgoingBitrate = fbs.MaxOutgoingBitrate
	t.MinOutgoingBitrate = fbs.MinOutgoingBitrate
	t.RtpPacketLossReceived = fbs.RtpPacketLossReceived
	t.RtpPacketLossSent = fbs.RtpPacketLossSent
}

type TransportTuple struct {
	LocalIp    string `json:"localIp"`
	LocalPort  uint16 `json:"localPort"`
	RemoteIp   string `json:"remoteIp"`
	RemotePort uint16 `json:"remotePort"`
	Protocol   string `json:"protocol"`
}

func (t *TransportTuple) Set(fbs *FBS__Transport.TupleT) {
	t.LocalIp = fbs.LocalIp
	t.LocalPort = fbs.LocalPort
	t.RemoteIp = fbs.RemoteIp
	t.RemotePort = fbs.RemotePort
	t.Protocol = strings.ToLower(FBS__Transport.EnumNamesProtocol[fbs.Protocol])
}

type GetTransportStatRequest struct {
	TransportId string
}

type GetTransportStatResponse struct {
	TransportStats
	Type             string         `json:"type"`
	IceRole          string         `json:"iceRole"`
	IceState         string         `json:"iceState"`
	IceSelectedTuple TransportTuple `json:"iceSelectedTuple"`
	DtlsState        string         `json:"dtlsState"`
}

func (r *GetTransportStatResponse) Set(typ string, fbs *FBS__WebRtcTransport.GetStatsResponseT) {
	r.TransportStats.Set(fbs.Base)
	r.Type = typ
	r.IceRole = strings.ToLower(FBS__WebRtcTransport.EnumNamesIceRole[fbs.IceRole])
	r.IceState = strings.ToLower(FBS__WebRtcTransport.EnumNamesIceState[fbs.IceState])
	r.IceSelectedTuple.Set(fbs.IceSelectedTuple)
	r.DtlsState = strings.ToLower(FBS__WebRtcTransport.EnumNamesDtlsState[fbs.DtlsState])
}
