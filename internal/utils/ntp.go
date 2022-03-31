package utils

import "math"

const (
	UnixNtpOffset     = 0x83AA7E80
	NtpFractionalUnit = uint64(1) << 32
)

type Ntp struct {
	Seconds   uint32
	Fractions uint32
}

func TimeMs2Ntp(ms int64) *Ntp {
	return &Ntp{
		Seconds:   uint32(ms / 1000),
		Fractions: uint32((float64(ms%1000) / 1000) * float64(NtpFractionalUnit)),
	}
}

func Ntp2TimeMs(ntp Ntp) int64 {
	return int64(ntp.Seconds*1000) + int64(math.Round(float64(ntp.Fractions)*1000)/float64(NtpFractionalUnit))
}
