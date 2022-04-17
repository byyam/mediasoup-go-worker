package rtc

import (
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
)

const (
	InactivityCheckInterval        = 1500 // In ms.
	InactivityCheckIntervalWithDtx = 5000 // In ms.
)

type TransmissionCounter struct {
	spatialLayerCounters [][]*RtpDataCounter
}

func newTransmissionCounter(spatialLayers, temporalLayers uint8, windowSize int) *TransmissionCounter {
	t := &TransmissionCounter{
		// Reserve vectors capacity.
		spatialLayerCounters: make([][]*RtpDataCounter, spatialLayers),
	}

	for idx, _ := range t.spatialLayerCounters {
		temporalLayerCounter := make([]*RtpDataCounter, 0)
		for tIdx := uint8(0); tIdx < temporalLayers; tIdx++ {
			r := NewRtpDataCounter(windowSize)
			temporalLayerCounter = append(temporalLayerCounter, r)
		}
		t.spatialLayerCounters[idx] = temporalLayerCounter
	}

	return t
}

func (p *TransmissionCounter) Update(packet *rtpparser.Packet) {
	// todo: support svc
	counter := p.spatialLayerCounters[0][0]
	counter.Update(packet)
}

func (p *TransmissionCounter) GetBitrate(nowMs int64) (rate uint32) {
	for _, spatialLayers := range p.spatialLayerCounters {
		for _, temporalLayer := range spatialLayers {
			rate += temporalLayer.GetBitrate(nowMs)
		}
	}
	return
}

func (p *TransmissionCounter) GetBytes() (bytes int64) {
	for _, spatialLayers := range p.spatialLayerCounters {
		for _, temporalLayer := range spatialLayers {
			bytes += temporalLayer.GetBytes()
		}
	}
	return
}

func (p *TransmissionCounter) GetPacketCount() (packetCount int64) {
	for _, spatialLayers := range p.spatialLayerCounters {
		for _, temporalLayer := range spatialLayers {
			packetCount += temporalLayer.GetPacketCount()
		}
	}
	return
}
