package protoo

import (
	"context"
	"fmt"
	"time"

	"github.com/jiyeyuran/go-protoo"
	"github.com/rs/zerolog"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/pkg/taskloop"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/byyam/mediasoup-go-worker/signaldefine"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/basehandler"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/demoutils"
)

const timeout = time.Second * 3

type QueryParams struct {
	RoomId   string `json:"roomId"`
	PeerId   string `json:"peerId"`
	RouterId string `json:"routerId"`
}

type ProtooHandler struct {
	basehandler.BaseHandler
	logger      zerolog.Logger
	queryParams *QueryParams
}

func NewHandler(worker *mediasoup_go_worker.SimpleWorker, queryParams *QueryParams) *ProtooHandler {
	return &ProtooHandler{
		BaseHandler: basehandler.BaseHandler{
			Worker: worker,
		},
		queryParams: queryParams,
		logger:      zerowrapper.NewScope(fmt.Sprintf("protoo[%s-%s]", queryParams.RoomId, queryParams.PeerId)),
	}
}

// mediasoup defined proto
func (h *ProtooHandler) HandleProtooMessage(req protoo.Message) *protoo.Message {
	var data interface{}
	var err *protoo.Error
	defer func() {
		if err == nil {
			h.logger.Info().Msgf("[HandleProtooMessage] success, method:%s", req.Method)
		} else {
			h.logger.Error().Msgf("[HandleProtooMessage] failed, method:%s, error:%+v", req.Method, err)
		}
	}()

	var protooFn func(message protoo.Message) (interface{}, *protoo.Error)
	switch req.Method {
	case signaldefine.MethodGetRouterRtpCapabilities:
		protooFn = h.GetRouterRtpCapabilities
	case signaldefine.MethodCreateWebRtcTransport:
		protooFn = h.CreateWebRtcTransport
	case signaldefine.MethodJoin:
		protooFn = h.Join
	case signaldefine.MethodConnectWebRtcTransport:
		protooFn = h.ConnectWebRtcTransport
	case signaldefine.MethodGetTransportStats:
		protooFn = h.GetTransportStats
	case signaldefine.MethodProduce:
		protooFn = h.Produce
	case signaldefine.MethodGetProducerStats:
		protooFn = h.GetProducerStats
	default:
		h.logger.Warn().Msgf("[HandleProtooMessage] unknown signal method: %s", req.Method)
		err = demoutils.ErrUnknownMethod
	}
	// default
	if protooFn == nil {
		protooFn = func(message protoo.Message) (interface{}, *protoo.Error) {
			h.logger.Warn().Msgf("[HandleProtooMessage] unknown signal method: %s", message.Method)
			return nil, demoutils.ErrUnknownMethod
		}
	}

	// run task
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	taskErr := taskloop.RunTask(ctx, func() error {
		data, err = protooFn(req)
		return nil
	})
	if taskErr != nil {
		h.logger.Error().Err(taskErr).Msgf("taskloop run failed")
		err = demoutils.ErrServerError
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
