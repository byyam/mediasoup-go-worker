package webrtctransport

import (
	"fmt"

	"github.com/jiyeyuran/go-protoo"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
)

func GetRouterId(w *mediasoup_go_worker.SimpleWorker) string {
	return fmt.Sprintf("router-%d", w.GetPid())
}

const (
	ErrCodeServer = 10000
)

func ServerError(err error) *protoo.Error {
	return protoo.NewError(ErrCodeServer, err.Error())
}

var (
	ErrUnknownMethod = protoo.NewError(100, "unknown method")
	ErrInvalidParam  = protoo.NewError(101, "invalid param")
)

func GetProducerId(streamId uint64) string {
	return fmt.Sprintf("%d", streamId)
}
