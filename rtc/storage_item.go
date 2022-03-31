package rtc

import "github.com/pion/rtp"

type StorageItem struct {
	packet     *rtp.Packet
	store      [MtuSize + 100]uint8
	resentAtMs uint64
	sentTimes  uint8
	rtxEncoded bool
}
