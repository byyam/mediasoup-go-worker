package rtc

import (
	"sync"

	"github.com/pion/rtp"
)

type RtpListener struct {
	ssrcTable sync.Map
	midTable  sync.Map
	ridTable  sync.Map
}

func newRtpListener() *RtpListener {
	return &RtpListener{}
}

func (r *RtpListener) AddProducer(producer *Producer) {

}

func (r *RtpListener) GetProducer(packet *rtp.Packet) *Producer {

	return nil
}
