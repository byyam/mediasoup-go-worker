package udpmux

import "testing"

func TestNewUdpMux(t *testing.T) {
	_, err := NewUdpMux("udp", "", 40000)
	if err != nil {
		t.Error(err)
	}
	select {}
}
