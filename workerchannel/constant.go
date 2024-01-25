package workerchannel

import (
	"github.com/byyam/mediasoup-go-worker/fbs/FBS/Notification"
)

const (
	NativeJsonVersion = "3.9.0"
	NativeVersion     = "3.10.5"
	FlatBufferVersion = "3.13.0"
)

// start from CustomerPipeStart
const (
	ConsumerChannelFd        = 3 // read fd
	ProducerChannelFd        = 4 // write fd
	PayloadConsumerChannelFd = 5
	PayloadProducerChannelFd = 6
)

var EventMap = map[Notification.Event]string{
	Notification.EventWORKER_RUNNING: "running",
}
