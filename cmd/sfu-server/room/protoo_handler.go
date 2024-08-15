package room

import (
	"fmt"

	"github.com/jiyeyuran/go-protoo"
	"github.com/rs/zerolog"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/byyam/mediasoup-go-worker/signaldefine"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/basehandler"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/demoutils"
)

type QueryParams struct {
	RoomId string `json:"roomId"`
	PeerId string `json:"peerId"`
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
		h.logger.Info().Msgf("[HandleProtooMessage] done, method:%s, err:%v", req.Method, err)
	}()
	switch req.Method {
	case signaldefine.MethodGetRouterRtpCapabilities:
		data, err = h.GetRouterRtpCapabilities(req)
	case signaldefine.MethodCreateWebRtcTransport:
		data, err = h.CreateWebRtcTransport(req)
	case signaldefine.MethodJoin:
		data, err = h.Join(req)
	case signaldefine.MethodConnectWebRtcTransport:
		data, err = h.ConnectWebRtcTransport(req)
	default:
		h.logger.Warn().Msgf("[HandleProtooMessage] unknown signal method: %s", req.Method)
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
