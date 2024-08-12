package rtctime

import (
	"fmt"
	"testing"
	"time"
)

func Test_ToNtpTime(t *testing.T) {
	ntp := NtpTime(16840898373497280725)
	timeValue := ntp.Time()
	fmt.Printf("time:%+v", timeValue)
}

func TestNtpTime32(t *testing.T) {
	now := time.Now()
	npt32 := NtpTime32(now)
	fmt.Printf("npt32:%+v\n", npt32)
}
