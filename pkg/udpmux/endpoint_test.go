package udpmux

import (
	"github.com/pion/rtp"
	"testing"
	"time"
)

func TestNewEndPoint(t *testing.T) {
	mux, err := NewUdpMux("127.0.0.1", "40000", nil)
	if err != nil {
		t.Error(err)
	}
	e, err := mux.AddEndPoint("127.0.0.1", "40001")
	if err != nil {
		t.Error(err)
	}
	e.OnRead(func(data []byte) {
		rtpPacket := &rtp.Packet{}
		_ = rtpPacket.Unmarshal(data)
		t.Log("OnRead", rtpPacket.String())
		_, _ = e.Write(data)
	})
	time.AfterFunc(time.Second*10, func() {
		t.Logf("close end point")
		e.Close()
	})
	select {}
}
