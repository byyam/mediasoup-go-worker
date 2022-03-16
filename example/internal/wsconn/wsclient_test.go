package wsconn

import (
	"testing"

	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
)

func TestWsClient_Conn(t *testing.T) {
	c, err := NewWsClient(WsClientOpt{
		Path: "echo",
	}).Conn()
	if err != nil {
		t.Fatal("conn error")
	}
	defer func() {
		_ = c.Close()
	}()
	// select {}
}

func TestWsClient_Request(t *testing.T) {
	req := isignal.UnPublishRequest{StreamId: 123}
	rsp, err := NewWsClient(WsClientOpt{
		Path: "echo",
	}).Request(isignal.MethodUnPublish, req)
	if err != nil {
		t.Fatal("request unpublish error", err)
	}
	t.Logf("rsp:%+v", rsp)
}
