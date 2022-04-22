package rtc

import (
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
)

type StorageItem struct {
	packet     *rtpparser.Packet
	store      [MtuSize + 100]uint8
	resentAtMs int64
	sentTimes  uint8
	rtxEncoded bool
}
