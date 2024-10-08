package demoutils

import (
	"fmt"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"

	"github.com/jiyeyuran/go-protoo"
)

const (
	ErrCodeServer = 10000
)

func ServerError(err error) *protoo.Error {
	return protoo.NewError(ErrCodeServer, err.Error())
}

var (
	ErrUnknownMethod = protoo.NewError(100, "unknown method")
	ErrInvalidParam  = protoo.NewError(101, "invalid param")
	ErrServerError   = protoo.NewError(102, "server error")
)

func GetProducerId(streamId uint64) string {
	return fmt.Sprintf("%d", streamId)
}

func GetRouterId(w *mediasoup_go_worker.SimpleWorker) string {
	return fmt.Sprintf("router-%d", w.GetPid())
}
