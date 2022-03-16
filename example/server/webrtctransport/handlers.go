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

func (h *Handler) PublishHandler(req protoo.Message) *protoo.Message {
	rspData := isignal.PublishResponse{
		TransportId: "demoId",
	}
	rsp := protoo.CreateSuccessResponse(req, rspData)
	return &rsp
}

func (h *Handler) UnPublishHandler(req protoo.Message) *protoo.Message {
	h.worker.FillJson()
	rspData := isignal.UnPublishResponse{}
	rsp := protoo.CreateSuccessResponse(req, rspData)
	return &rsp
}
