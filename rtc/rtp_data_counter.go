package rtc

import (
	"github.com/byyam/mediasoup-go-worker/pkg/ratecalculator"
	"github.com/byyam/mediasoup-go-worker/pkg/rtctime"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
)

type RtpDataCounter struct {
	rate    *ratecalculator.RateCalculator
	packets int64
}

func NewRtpDataCounter(windowSizeMs int) *RtpDataCounter {
	size := 2500
	if windowSizeMs > 0 {
		size = windowSizeMs
	}
	return &RtpDataCounter{
		rate:    ratecalculator.NewRateCalculator(size, 0, 0, nil),
		packets: 0,
	}
}

func (p RtpDataCounter) GetBytes() int64 {
	return p.rate.GetBytes()
}

func (p RtpDataCounter) GetPacketCount() int64 {
	return p.packets
}

func (p RtpDataCounter) GetBitrate(nowMs uint64) uint32 {
	return p.rate.GetRate(nowMs)
}

func (p *RtpDataCounter) Update(packet *rtpparser.Packet) {
	nowMs := rtctime.GetTimeMs()
	p.packets++
	p.rate.Update(packet.GetLen(), nowMs)
}
