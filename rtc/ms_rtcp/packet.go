package ms_rtcp

import "time"

const (
	MaxAudioIntervalMs = time.Millisecond * 5000
	MaxVideoIntervalMs = time.Millisecond * 1000
)

func CountRequestedPackets() uint32 {
	return 0
}
