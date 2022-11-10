package rtc

import (
	"fmt"
	"sync"

	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
)

type RtpListener struct {
	ssrcTable sync.Map
	midTable  sync.Map
	ridTable  sync.Map
}

func newRtpListener() *RtpListener {
	return &RtpListener{}
}

func (r *RtpListener) AddProducer(producer *Producer) error {
	// Add entries into the ssrcTable.
	for _, encoding := range producer.RtpParameters.Encodings {
		// Check encoding.ssrc.
		if encoding.Ssrc != 0 {
			if _, ok := r.ssrcTable.Load(encoding.Ssrc); ok {
				return fmt.Errorf("ssrc already exists in RTP listener [ssrc:%d]", encoding.Ssrc)
			}
			r.ssrcTable.Store(encoding.Ssrc, producer)
		}
		// Check encoding.rtx.ssrc.
		if encoding.Rtx != nil && encoding.Rtx.Ssrc != 0 {
			if _, ok := r.ssrcTable.Load(encoding.Rtx.Ssrc); ok {
				return fmt.Errorf("RTX ssrc already exists in RTP listener [ssrc:%d]", encoding.Ssrc)
			}
			r.ssrcTable.Store(encoding.Rtx.Ssrc, producer)
		}
	}
	// Add entries into midTable.
	if producer.RtpParameters.Mid != "" {
		mid := producer.RtpParameters.Mid
		if _, ok := r.midTable.Load(mid); ok {
			return fmt.Errorf("MID already exists in RTP listener [mid:%s]", mid)
		}
		r.midTable.Store(mid, producer)
	}
	// Add entries into ridTable.
	for _, encoding := range producer.RtpParameters.Encodings {
		rid := encoding.Rid
		if rid == "" {
			continue
		}
		// Just fail if no MID is given.
		if _, ok := r.ridTable.Load(rid); ok && producer.RtpParameters.Mid == "" {
			return fmt.Errorf("RID already exists in RTP listener and no MID is given [rid:%s]", rid)
		}
		r.ridTable.Store(rid, producer)
	}
	return nil
}

func (r *RtpListener) GetProducer(packet *rtpparser.Packet) *Producer {
	// First lookup into the SSRC table.
	if producer := r.GetProducerBySSRC(packet.SSRC); producer != nil {
		return producer
	}
	// Otherwise lookup into the MID table.
	if producer := r.GetProducerByMID(packet); producer != nil {
		return producer
	}
	// Otherwise lookup into the RID table.
	if producer := r.GetProducerByRID(packet); producer != nil {
		return producer
	}
	return nil
}

// GetProducerBySSRC is over-write to GetProducer in mediasoup
func (r *RtpListener) GetProducerBySSRC(ssrc uint32) *Producer {
	value, ok := r.ssrcTable.Load(ssrc)
	if !ok {
		return nil
	}
	return value.(*Producer)
}

func (r *RtpListener) GetProducerByMID(packet *rtpparser.Packet) *Producer {
	value, ok := r.midTable.Load(packet.GetMid())
	if !ok {
		return nil
	}
	// Fill the ssrc table.
	// NOTE: We may be overriding an exiting SSRC here, but we don't care.
	producer := value.(*Producer)
	r.ssrcTable.Store(packet.SSRC, producer)

	return producer
}

func (r *RtpListener) GetProducerByRID(packet *rtpparser.Packet) *Producer {
	value, ok := r.ridTable.Load(packet.GetRid())
	if !ok {
		return nil
	}
	// Fill the ssrc table.
	// NOTE: We may be overriding an exiting SSRC here, but we don't care.
	producer := value.(*Producer)
	r.ssrcTable.Store(packet.SSRC, producer)

	return producer
}
