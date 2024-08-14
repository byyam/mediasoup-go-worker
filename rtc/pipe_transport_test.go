package rtc

import (
	"testing"

	FBS__PipeTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/PipeTransport"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
)

func TestPipeTransport(t *testing.T) {
	_, err := newPipeTransport(pipeTransportParam{
		optionsFBS: &FBS__PipeTransport.PipeTransportOptionsT{
			ListenInfo: &FBS__Transport.ListenInfoT{
				Ip:          "127.0.0.1",
				Port:        40001,
				AnnouncedIp: "",
			},
		},
		transportParam: transportParam{},
	})
	if err != nil {
		t.Fatal(err)
	}
	select {}
}
