package isignal

import "github.com/byyam/mediasoup-go-worker/mediasoupdata"

type CreatePipeTransportRequest struct {
	PipeTransportOffer
}

type CreatePipeTransportResponse struct {
	TransportId string `json:"id"`
	PipeTransportAnswer
}

type PipeTransportOffer struct {
	mediasoupdata.PipeTransportOptions
	EndPointIp   string `json:"endPointIpIp"`
	EndPointPort uint16 `json:"endPointPort"`
}

type PipeTransportAnswer struct {
	*mediasoupdata.PipeTransportData `json:"pipeTransportData,omitempty"`
}

type PublishOnPipeTransportRequest struct {
	StreamId    uint64 `json:"streamId"`
	TransportId string `json:"transportId"`
	PublishOffer
}

type PublishOnPipeTransportResponse struct {
	PublishAnswer
}

type SubscribeOnPipeTransportRequest struct {
	StreamId    uint64 `json:"streamId"`
	TransportId string `json:"transportId"`
	SubscribeOffer
}

type SubscribeOnPipeTransportResponse struct {
	SubscribeId string `json:"subscribeId"`
	SubscribeAnswer
}
