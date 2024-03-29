package signaldefine

import (
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
	Id             string                       `json:"id"`
	IceParameters  mediasoupdata.IceParameters  `json:"iceParameters"`
	IceCandidates  []mediasoupdata.IceCandidate `json:"iceCandidates"`
	DtlsParameters mediasoupdata.DtlsParameters `json:"dtlsParameters"`
	SctpParameters mediasoupdata.SctpParameters `json:"sctpParameters"`
}

type ConnectWebRtcTransportRequest struct {
	TransportId    string
	DtlsParameters mediasoupdata.DtlsParameters
}

type ConnectWebRtcTransportResponse struct {
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
