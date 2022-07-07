package rtc

import (
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"testing"
)

func TestPipeTransport(t *testing.T) {
	_, err := newPipeTransport(pipeTransportParam{
		options: mediasoupdata.PipeTransportOptions{
			ListenIp: mediasoupdata.TransportListenIp{
				Ip:          "127.0.0.1",
				AnnouncedIp: "",
			},
			Port: 40001,
		},
		transportParam: transportParam{},
	})
	if err != nil {
		t.Fatal(err)
	}
	select {}
}
