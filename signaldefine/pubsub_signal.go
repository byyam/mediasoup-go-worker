package signaldefine

import (
	"encoding/json"

	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
)

type PublishOffer struct {
	Kind          mediasoupdata.MediaKind     `json:"kind"`
	RtpParameters mediasoupdata.RtpParameters `json:"rtpParameters"`
	AppData       interface{}                 `json:"appData"`
}

type PublishAnswer struct {
}

type SubscribeOffer struct {
	Kind            mediasoupdata.MediaKind        `json:"kind"`
	RtpCapabilities *mediasoupdata.RtpCapabilities `json:"rtpCapabilities"`
	AppData         interface{}                    `json:"appData"`
}

type SubscribeAnswer struct {
	Kind          mediasoupdata.MediaKind     `json:"kind"`
	RtpParameters mediasoupdata.RtpParameters `json:"rtpParameters"`
	AppData       interface{}                 `json:"appData"`
}

func (p SubscribeAnswer) String() string {
	jsonByte, _ := json.Marshal(&p)
	return string(jsonByte)
}
