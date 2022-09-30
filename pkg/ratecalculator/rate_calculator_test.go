package ratecalculator

import (
	"math/rand"
	"testing"
	"time"

	"github.com/byyam/mediasoup-go-worker/pkg/rtctime"
)

func BenchmarkRateCalculator_Update(b *testing.B) {
	rate := NewRateCalculator(2500, 0, 0, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nowMs := rtctime.GetTimeMs()
		rate.Update(300, nowMs)
	}
}

func TestRateCalculator_Update(t *testing.T) {
	rate := NewRateCalculator(100, 0, 10, nil)
	for i := 0; i < 20; i++ {
		nowMs := rtctime.GetTimeMs()
		jitterMs := nowMs + int64(rand.Intn(10)) - int64(rand.Intn(10))
		rate.Update(100, jitterMs)
		time.Sleep(time.Millisecond * 10)
	}
}
