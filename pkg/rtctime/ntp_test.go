package rtctime

import (
	"fmt"
	"testing"
)

func Test_ToNtpTime(t *testing.T) {
	ntp := NtpTime(16840898373497280725)
	timeValue := ntp.Time()
	fmt.Printf("time:%+v", timeValue)
}
