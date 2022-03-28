package rtc

import (
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/pion/rtp"
)

type RtpDataCounter struct {
	rate    *RateCalculator
	packets int
}

func newRtpDataCounter(windowSizeMs int) *RtpDataCounter {
	size := 2500
	if windowSizeMs > 0 {
		size = windowSizeMs
	}
	return &RtpDataCounter{
		rate:    newRateCalculator(size, 0, 0),
		packets: 0,
	}
}

func (p RtpDataCounter) GetBytes() int {
	return p.rate.GetBytes()
}

func (p RtpDataCounter) GetPacketCount() int {
	return p.packets
}

func (p RtpDataCounter) GetBitrate(nowMs int64) uint32 {
	return p.rate.GetRate(nowMs)
}

func (p *RtpDataCounter) Update(packet *rtp.Packet) {
	nowMs := utils.GetTimeMs()
	p.packets++
	p.rate.Update(len(packet.Payload), nowMs) // todo: packet size
}
