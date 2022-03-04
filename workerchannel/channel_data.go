package workerchannel

import (
	"encoding/json"

	"github.com/ragsagar/netstringer"
)

type channelData struct {
	Id int64 `json:"id,omitempty"`
	// request
	Method   string          `json:"method,omitempty"`
	Internal json.RawMessage `json:"internal,omitempty"`
	// response
	Accepted bool   `json:"accepted,omitempty"`
	Error    string `json:"error,omitempty"`
	Reason   string `json:"reason,omitempty"`
	// notification
	TargetId string `json:"targetId,omitempty"`
	Event    string `json:"event,omitempty"`
	// common data
	Data json.RawMessage `json:"data,omitempty"`
}

func (c *channelData) Marshal() ([]byte, error) {
	jsonByte, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	buf := netstringer.Encode(jsonByte)
	return buf, nil
}
