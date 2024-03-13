package workerchannel

import (
	"encoding/json"
	"fmt"

	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Response "github.com/byyam/mediasoup-go-worker/fbs/FBS/Response"
)

type InternalData struct {
	RouterId       string `json:"routerId,omitempty"`
	TransportId    string `json:"transportId,omitempty"`
	ProducerId     string `json:"producerId,omitempty"`
	ConsumerId     string `json:"consumerId,omitempty"`
	DataProducerId string `json:"dataProducerId,omitempty"`
	DataConsumerId string `json:"dataConsumerId,omitempty"`
	RtpObserverId  string `json:"rtpObserverId,omitempty"`
}

func (i *InternalData) Unmarshal(data json.RawMessage) error {
	return json.Unmarshal(data, i)
}

type RequestData struct {
	MethodType FBS__Request.Method
	Method     string
	HandlerId  string `json:"handlerId,omitempty"`
	Internal   InternalData
	Data       json.RawMessage
	// FBS request
	Request *FBS__Request.RequestT
}

func (d RequestData) String() string {
	return fmt.Sprintf("Id:%d,HandlerId:%s,Method:%s",
		d.Request.Id,
		d.Request.HandlerId,
		FBS__Request.EnumNamesMethod[d.MethodType])
}

type ResponseData struct {
	Id         uint32
	MethodType FBS__Request.Method
	Err        error
	Data       json.RawMessage
	// FBS data
	RspBody *FBS__Response.BodyT
}

func (d ResponseData) String() string {

	var bodyType string
	if d.RspBody != nil {
		bodyType = FBS__Response.EnumNamesBody[d.RspBody.Type]
	}
	return fmt.Sprintf("Id:%d,Method:%s,BodyType:%s,Err:%v",
		d.Id,
		FBS__Request.EnumNamesMethod[d.MethodType],
		bodyType,
		d.Err)
}
