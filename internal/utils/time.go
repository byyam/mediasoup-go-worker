package utils

import "time"

func GetTimeMs() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetTimeUs() int64 {
	return time.Now().UnixNano() / 1000
}

func GetTimeNs() int64 {
	return time.Now().UnixNano()
}
