package rtctime

import (
	"math"
	"time"
)

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

// -------------------------------------

var (
	ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
)

type NtpTime uint64

func (t NtpTime) Duration() time.Duration {
	sec := (t >> 32) * 1e9
	frac := (t & 0xffffffff) * 1e9
	nsec := frac >> 32
	if uint32(frac) >= 0x80000000 {
		nsec++
	}
	return time.Duration(sec + nsec)
}

func (t NtpTime) Time() time.Time {
	return ntpEpoch.Add(t.Duration())
}

func ToNtpTime(t time.Time) NtpTime {
	nsec := uint64(t.Sub(ntpEpoch))
	sec := nsec / 1e9
	nsec = (nsec - sec*1e9) << 32
	frac := nsec / 1e9
	if nsec%1e9 >= 1e9/2 {
		frac++
	}
	return NtpTime(sec<<32 | frac)
}

// ------------------------------------------
