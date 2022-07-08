package pipetransport

import (
	"encoding/json"
	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	utils2 "github.com/byyam/mediasoup-go-worker/example/server/utils"
	"github.com/byyam/mediasoup-go-worker/example/server/workerapi"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	worker *mediasoup_go_worker.SimpleWorker
	logger utils.Logger
}

func NewHandler(worker *mediasoup_go_worker.SimpleWorker) *Handler {
	return &Handler{
		worker: worker,
		logger: utils.NewLogger("http-handler"),
	}
}

func (h *Handler) HandleCreatePipeTransport(w http.ResponseWriter, r *http.Request) {
	var req isignal.CreatePipeTransportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Info("HandleCreatePipeTransport:%+v", req)
	transportId := uuid.New().String()
	// create pipe transport
	transportData, err := workerapi.CreatePipeTransport(h.worker, workerapi.ParamCreatePipeTransport{
		RouterId:    utils2.GetRouterId(h.worker),
		TransportId: transportId,
		Options:     req.PipeTransportOptions,
	})
	if err != nil {
		h.logger.Error("create pipe transport failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Debug("create pipe-transport done, data:%+v", transportData)
	// connect pipe transport
	transportConnectOptions := mediasoupdata.TransportConnectOptions{
		Ip:   req.RemoteIp,
		Port: req.RemotePort,
	}
	if err := workerapi.TransportConnect(h.worker, workerapi.ParamTransportConnect{
		RouterId:    utils2.GetRouterId(h.worker),
		TransportId: transportId,
		Options:     transportConnectOptions,
	}); err != nil {
		h.logger.Error("connect pipe transport failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var rsp isignal.CreatePipeTransportResponse
	rsp.TransportId = transportId
	rsp.PipeTransportData = transportData
	_ = json.NewEncoder(w).Encode(rsp)
}
