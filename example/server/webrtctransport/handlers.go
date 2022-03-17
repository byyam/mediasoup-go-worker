package webrtctransport

import (
	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/jiyeyuran/go-protoo"
)

type Handler struct {
	worker *mediasoup_go_worker.SimpleWorker
	logger utils.Logger
}

func NewHandler(worker *mediasoup_go_worker.SimpleWorker) *Handler {
	return &Handler{
		worker: worker,
		logger: utils.NewLogger("websocket-handler"),
	}
}

func (h *Handler) HandleProtooMessage(req protoo.Message) *protoo.Message {
	var data interface{}
	var err *protoo.Error
	switch req.Method {
	case isignal.MethodPublish:
		data, err = h.publishHandler(req)
	case isignal.MethodSubscribe:
		data, err = h.subscribeHandler(req)
	default:
		err = ErrUnknownMethod
	}
	// create response protoo message
	if err != nil {
		rsp := protoo.CreateErrorResponse(req, err)
		return &rsp
	} else {
		rsp := protoo.CreateSuccessResponse(req, data)
		return &rsp
	}
}
