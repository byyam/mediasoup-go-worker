package workerchannel

import (
	"encoding/json"
	"fmt"
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
	Method    string
	HandlerId string `json:"handlerId,omitempty"`
	Internal  InternalData
	Data      json.RawMessage
}

func (d RequestData) String() string {
	return fmt.Sprintf("HandlerId:%s,Method:%s,Internal:%+v,Data:%s", d.HandlerId, d.Method, d.Internal, string(d.Data))
}

type ResponseData struct {
	Err  error
	Data json.RawMessage
}

func (d ResponseData) String() string {
	return fmt.Sprintf("Err:%v,Data:%s", d.Err, string(d.Data))
}
