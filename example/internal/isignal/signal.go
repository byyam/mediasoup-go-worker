package isignal

import "github.com/byyam/mediasoup-go-worker/mediasoupdata"

const (
	MethodPublish     = "publish"
	MethodUnPublish   = "unPublish"
	MethodSubscribe   = "subscribe"
	MethodUnSubscribe = "unSubscribe"
)

type PublishRequest struct {
	StreamId    uint64               `json:"streamId"`
	TransportId string               `json:"transportId,omitempty"`
	Offer       WebRtcTransportOffer `json:"webrtcTransportOffer"`

	Kind          mediasoupdata.MediaKind     `json:"kind"`
	RtpParameters mediasoupdata.RtpParameters `json:"rtpParameters"`
	AppData       interface{}                 `json:"appData"`
}

type PublishResponse struct {
	TransportId string                `json:"transportId,omitempty"`
	Answer      WebRtcTransportAnswer `json:"webrtcTransportAnswer"`
}

type UnPublishRequest struct {
	StreamId uint64 `json:"streamId"`
}

type UnPublishResponse struct{}

type SubscribeRequest struct {
	StreamId        uint64                         `json:"streamId"`
	TransportId     string                         `json:"transportId,omitempty"`
	Offer           WebRtcTransportOffer           `json:"webrtcTransportOffer"`
	Kind            mediasoupdata.MediaKind        `json:"kind"`
	RtpCapabilities *mediasoupdata.RtpCapabilities `json:"rtpCapabilities"`
	AppData         interface{}                    `json:"appData"`
}

type SubscribeResponse struct {
	SubscribeId   string                      `json:"subscribeId"`
	TransportId   string                      `json:"transportId,omitempty"`
	Answer        WebRtcTransportAnswer       `json:"webrtcTransportAnswer"`
	Kind          mediasoupdata.MediaKind     `json:"kind"`
	RtpParameters mediasoupdata.RtpParameters `json:"rtpParameters"`
	AppData       interface{}                 `json:"appData"`
}

type UnSubscribeRequest struct {
	StreamId    uint64 `json:"streamId"`
	SubscribeId string `json:"subscribeId"`
}

type UnSubscribeResponse struct{}

type WebRtcTransportOffer struct {
	ForceTcp       bool                         `json:"forceTcp"`
	DtlsParameters mediasoupdata.DtlsParameters `json:"dtlsParameters"`
}

type WebRtcTransportAnswer struct {
	IceParameters  mediasoupdata.IceParameters  `json:"iceParameters"`
	IceCandidates  []mediasoupdata.IceCandidate `json:"iceCandidates"`
	DtlsParameters mediasoupdata.DtlsParameters `json:"dtlsParameters"`
}
