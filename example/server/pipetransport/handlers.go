package pipetransport

import (
	"encoding/json"
	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/example/internal/demoutils"
	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/server/basehandler"
	"github.com/byyam/mediasoup-go-worker/example/server/workerapi"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	basehandler.BaseHandler
	logger utils.Logger
}

func NewHandler(worker *mediasoup_go_worker.SimpleWorker) *Handler {
	return &Handler{
		BaseHandler: basehandler.BaseHandler{
			Worker: worker,
		},
		logger: utils.NewLogger("http-handler"),
	}
}

func (h *Handler) HandlePipeTransportCreateAndConnect(w http.ResponseWriter, r *http.Request) {
	var req isignal.CreatePipeTransportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Info("HandlePipeTransportCreateAndConnect:%+v", req)
	transportId := uuid.New().String()
	// create pipe transport
	transportData, err := workerapi.CreatePipeTransport(h.Worker, workerapi.ParamCreatePipeTransport{
		RouterId:    demoutils.GetRouterId(h.Worker),
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
		Ip:   req.EndPointIp,
		Port: req.EndPointPort,
	}
	if err := workerapi.TransportConnect(h.Worker, workerapi.ParamTransportConnect{
		RouterId:    demoutils.GetRouterId(h.Worker),
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

func (h *Handler) HandlePublish(w http.ResponseWriter, r *http.Request) {
	var req isignal.PublishOnPipeTransportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Info("HandlePublish:%+v", req)
	// publish on pipe transport
	produceOptions, err := basehandler.ProducerOptions(req.Kind, req.StreamId, req.RtpParameters)
	if err != nil {
		h.logger.Error("create producer options failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := workerapi.TransportProduce(h.Worker, workerapi.ParamTransportProduce{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: req.TransportId,
		ProducerId:  produceOptions.Id,
		Options:     *produceOptions,
	}); err != nil {
		h.logger.Error("publish on pipe transport failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var rsp isignal.PublishOnPipeTransportResponse
	_ = json.NewEncoder(w).Encode(rsp)
}

func (h *Handler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	var req isignal.SubscribeOnPipeTransportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Info("HandleSubscribe:%+v", req)
	// subscribe on pipe transport
	consumerOptions, err := h.ConsumerOptions(req.StreamId, req.RtpCapabilities)
	if err != nil {
		h.logger.Error("create consumer options failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	consumerId := uuid.New().String()
	consumerData, err := workerapi.TransportConsume(h.Worker, workerapi.ParamTransportConsume{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: req.TransportId,
		ProducerId:  consumerOptions.ProducerId,
		ConsumerId:  consumerId,
		Options:     *consumerOptions,
	})
	if err != nil {
		h.logger.Error("subscribe on pipe transport failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Debug("subscribe on pipe-transport done, data:%+v", consumerData)
	var rsp isignal.SubscribeOnPipeTransportResponse
	rsp.SubscribeId = consumerId
	rsp.SubscribeAnswer.Kind = consumerOptions.Kind
	rsp.SubscribeAnswer.RtpParameters = consumerOptions.RtpParameters
	_ = json.NewEncoder(w).Encode(rsp)
}
