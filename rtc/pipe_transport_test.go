package rtc

import (
	"testing"

	mediasoupdata2 "github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
)

func TestPipeTransport(t *testing.T) {
	_, err := newPipeTransport(pipeTransportParam{
		options: mediasoupdata2.PipeTransportOptions{
			ListenIp: mediasoupdata2.TransportListenIp{
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
