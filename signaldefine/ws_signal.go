package signaldefine

import "github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"

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

	PublishOffer
}

type PublishResponse struct {
	TransportId string                `json:"transportId,omitempty"`
	Answer      WebRtcTransportAnswer `json:"webrtcTransportAnswer"`

	PublishAnswer
}

type UnPublishRequest struct {
	StreamId uint64 `json:"streamId"`
}

type UnPublishResponse struct{}

type SubscribeRequest struct {
	StreamId    uint64               `json:"streamId"`
	TransportId string               `json:"transportId,omitempty"`
	Offer       WebRtcTransportOffer `json:"webrtcTransportOffer"`

	SubscribeOffer
}

type SubscribeResponse struct {
	SubscribeId string                `json:"subscribeId"`
	TransportId string                `json:"transportId,omitempty"`
	Answer      WebRtcTransportAnswer `json:"webrtcTransportAnswer"`

	SubscribeAnswer
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
