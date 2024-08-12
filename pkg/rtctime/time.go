package rtctime

import "time"

func GetTimeMs() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func GetTimeUs() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}

func GetTimeNs() uint64 {
	return uint64(time.Now().UnixNano())
}
