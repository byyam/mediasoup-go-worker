package udpmux

import "testing"

func TestNewUdpMux(t *testing.T) {
	_, err := NewUdpMux("", "40000", nil)
	if err != nil {
		t.Error(err)
	}
	select {}
}
