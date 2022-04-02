package workerchannel

import (
	"encoding/json"
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
