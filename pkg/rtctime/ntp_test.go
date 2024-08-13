package rtctime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeToNTPConversation(t *testing.T) {
	for i, cc := range []struct {
		ts uint64
	}{
		{
			ts: 0,
		},
		{
			ts: 65535,
		},
		{
			ts: 16606669245815957503,
		},
		{
			ts: 9487534653230284800,
		},
	} {
		t.Run(fmt.Sprintf("TimeToNTP/%v", i), func(t *testing.T) {
			fmt.Printf("i=%d, value=%0x TimeMs2Ntp:%+v\n", i, cc.ts, TimeMs2Ntp(cc.ts))
			assert.Equal(t, cc.ts, Ntp2TimeMs(TimeMs2Ntp(cc.ts)))
		})
	}
}

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
