package wsconn

import (
	"testing"

	"github.com/byyam/mediasoup-go-worker/signaldefine"
)

func TestWsClient_Conn(t *testing.T) {
	c, err := NewWsClient(WsClientOpt{
		Path: "echo",
	})
	if err != nil {
		t.Fatal("conn error")
	}
	defer func() {
		_ = c.Close()
	}()
	// select {}
}

func TestWsClient_Request(t *testing.T) {
	req := signaldefine.UnPublishRequest{StreamId: 123}
	c, err := NewWsClient(WsClientOpt{
		Path: "echo",
	})
	if err != nil {
		t.Fatal("NewWsClient error", err)
	}
	rsp, err := c.Request(signaldefine.MethodUnPublish, req)
	if err != nil {
		t.Fatal("request unpublish error", err)
	}
	t.Logf("rsp:%+v", rsp)
}
