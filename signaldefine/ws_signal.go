package signaldefine

import (
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
)

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
	IceParameters  *FBS__WebRtcTransport.IceParametersT  `json:"iceParameters"`
	IceCandidates  []*FBS__WebRtcTransport.IceCandidateT `json:"iceCandidates"`
	DtlsParameters *FBS__WebRtcTransport.DtlsParametersT `json:"dtlsParameters"`
}
