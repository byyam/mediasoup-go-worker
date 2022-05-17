package ratecalculator

import (
	utils2 "github.com/byyam/mediasoup-go-worker/utils"
	"math/rand"
	"testing"
	"time"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
)

func BenchmarkRateCalculator_Update(b *testing.B) {
	rate := NewRateCalculator(2500, 0, 0, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nowMs := utils.GetTimeMs()
		rate.Update(300, nowMs)
	}
}

func TestRateCalculator_Update(t *testing.T) {
	utils2.DefaultLevel = utils2.TraceLevel
	rate := NewRateCalculator(100, 0, 10, utils2.NewLogger("RateCalculator"))
	for i := 0; i < 20; i++ {
		nowMs := utils.GetTimeMs()
		jitterMs := nowMs + int64(rand.Intn(10)) - int64(rand.Intn(10))
		rate.Update(100, jitterMs)
		time.Sleep(time.Millisecond * 10)
	}
}
