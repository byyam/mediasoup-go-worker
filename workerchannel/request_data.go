package workerchannel

import "encoding/json"

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
	Method       string
	InternalData InternalData
	Data         json.RawMessage
}

type ResponseData struct {
	Err  error
	Data json.RawMessage
}
