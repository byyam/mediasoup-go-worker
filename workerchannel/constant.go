package workerchannel

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
