package rtctime

import (
	"fmt"
	"math"
	"time"
)

const (
	UnixNtpOffset     = 0x83AA7E80 // Seconds from Jan 1, 1900 to Jan 1, 1970.
	NtpFractionalUnit = uint64(1) << 32
)

// need uint64 to avoid overflow
type Ntp struct {
	Seconds   uint64
	Fractions uint64
}

func TimeMs2Ntp(ms uint64) *Ntp {
	return &Ntp{
		Seconds:   ms / 1000,
		Fractions: uint64((float64(ms%1000) / 1000) * float64(NtpFractionalUnit)),
	}
}

func Ntp2TimeMs(ntp *Ntp) uint64 {
	return ntp.Seconds*1000 + uint64(math.Round(float64(ntp.Fractions)*1000/float64(NtpFractionalUnit)))
}

func NtpTime32(t time.Time) uint32 {
	// seconds since 1st January 1900
	s := (float64(t.UnixNano()) / 1000000000.0) + 2208988800

	integerPart := uint32(s)
	fractionalPart := uint32((s - float64(integerPart)) * 0xFFFFFFFF)
	fmt.Printf("intergerPart:%d, fractionPart:%d\n", integerPart, fractionalPart)

	// higher 32 bits are the integer part, lower 32 bits are the fractional part
	return uint32(((uint64(integerPart)<<32 | uint64(fractionalPart)) >> 16) & 0xFFFFFFFF)
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
