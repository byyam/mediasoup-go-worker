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
	spatialLayers        uint8
	temporalLayers       uint8
}

func newTransmissionCounter(spatialLayers, temporalLayers uint8, windowSize int) *TransmissionCounter {
	t := &TransmissionCounter{
		spatialLayers:  spatialLayers,
		temporalLayers: temporalLayers,
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
	spatialLayer := packet.GetSpatialLayer()
	temporalLayer := packet.GetTemporalLayer()
	// Sanity check. Do not allow spatial layers higher than defined.
	if spatialLayer > p.spatialLayers-1 {
		spatialLayer = p.spatialLayers - 1
	}
	// Sanity check. Do not allow temporal layers higher than defined.
	if temporalLayer > p.temporalLayers-1 {
		temporalLayer = p.temporalLayers - 1
	}
	counter := p.spatialLayerCounters[spatialLayer][temporalLayer]
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
