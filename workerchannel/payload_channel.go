package workerchannel

type PayloadChannel struct {
	consumerFd int
	producerFd int
}

func NewPayloadChannel(consumerFd, producerFd int) *PayloadChannel {
	return &PayloadChannel{
		consumerFd: consumerFd,
		producerFd: producerFd,
	}
}
