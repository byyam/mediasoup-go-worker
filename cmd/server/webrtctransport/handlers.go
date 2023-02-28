package webrtctransport

import (
	"github.com/jiyeyuran/go-protoo"
	"github.com/rs/zerolog"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/byyam/mediasoup-go-worker/signaldefine"

	"github.com/byyam/mediasoup-go-worker/cmd/server/basehandler"
	"github.com/byyam/mediasoup-go-worker/cmd/server/demoutils"
)

type Handler struct {
	basehandler.BaseHandler
	logger zerolog.Logger
}

func NewHandler(worker *mediasoup_go_worker.SimpleWorker) *Handler {
	return &Handler{
		BaseHandler: basehandler.BaseHandler{
			Worker: worker,
		},
		logger: zerowrapper.NewScope("websocket-handler"),
	}
}

func (h *Handler) HandleProtooMessage(req protoo.Message) *protoo.Message {
	var data interface{}
	var err *protoo.Error
	switch req.Method {
	case signaldefine.MethodPublish:
		data, err = h.publishHandler(req)
	case signaldefine.MethodUnPublish:
		data, err = h.unPublishHandler(req)
	case signaldefine.MethodSubscribe:
		data, err = h.subscribeHandler(req)
	case signaldefine.MethodUnSubscribe:
		data, err = h.unSubscribeHandler(req)
	default:
		err = demoutils.ErrUnknownMethod
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
