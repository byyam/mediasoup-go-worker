package rtc

import "testing"

func TestTransmissionCounter_GetBytes(t *testing.T) {
	transmissionCounter := newTransmissionCounter(1, 1, 200)
	transmissionCounter.GetBytes()
}
