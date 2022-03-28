package rtc

import "github.com/pion/rtp"

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
		temporalLayerCounter := make([]*RtpDataCounter, temporalLayers)
		for tIdx := uint8(0); tIdx < temporalLayers; tIdx++ {
			r := newRtpDataCounter(windowSize)
			temporalLayerCounter = append(temporalLayerCounter, r)
		}
		t.spatialLayerCounters[idx] = temporalLayerCounter
	}

	return t
}

func (p *TransmissionCounter) Update(packet *rtp.Packet) {
	// todo
}
