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
	for _, encoding := range producer.RtpParameters.Encodings {
		r.ssrcTable.Store(encoding.Ssrc, producer)
	}
	// todo: rtx,mid,rid
}

func (r *RtpListener) GetProducer(packet *rtp.Packet) *Producer {
	value, ok := r.ssrcTable.Load(packet.SSRC)
	if !ok {
		return nil
	}
	return value.(*Producer)
}