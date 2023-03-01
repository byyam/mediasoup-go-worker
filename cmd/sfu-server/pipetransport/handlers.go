package pipetransport

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/byyam/mediasoup-go-worker/signaldefine"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/basehandler"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/demoutils"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/workerapi"
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
		logger: zerowrapper.NewScope("http-handler"),
	}
}

func (h *Handler) HandlePipeTransportCreateAndConnect(w http.ResponseWriter, r *http.Request) {
	var req signaldefine.CreatePipeTransportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Info().Msgf("HandlePipeTransportCreateAndConnect:%+v", req)
	transportId := uuid.New().String()
	// create pipe transport
	transportData, err := workerapi.CreatePipeTransport(h.Worker, workerapi.ParamCreatePipeTransport{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportId,
		Options:     req.PipeTransportOptions,
	})
	if err != nil {
		h.logger.Error().Msgf("create pipe transport failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Debug().Msgf("create pipe-transport done, data:%+v", transportData)
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
		h.logger.Error().Msgf("connect pipe transport failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var rsp signaldefine.CreatePipeTransportResponse
	rsp.TransportId = transportId
	rsp.PipeTransportData = transportData
	_ = json.NewEncoder(w).Encode(rsp)
}

func (h *Handler) HandlePublish(w http.ResponseWriter, r *http.Request) {
	var req signaldefine.PublishOnPipeTransportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Info().Msgf("HandlePublish:%+v", req)
	// publish on pipe transport
	produceOptions, err := basehandler.ProducerOptions(req.Kind, req.StreamId, req.RtpParameters)
	if err != nil {
		h.logger.Error().Msgf("create producer options failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := workerapi.TransportProduce(h.Worker, workerapi.ParamTransportProduce{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: req.TransportId,
		ProducerId:  produceOptions.Id,
		Options:     *produceOptions,
	}); err != nil {
		h.logger.Error().Msgf("publish on pipe transport failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var rsp signaldefine.PublishOnPipeTransportResponse
	_ = json.NewEncoder(w).Encode(rsp)
}

func (h *Handler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	var req signaldefine.SubscribeOnPipeTransportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Info().Msgf("HandleSubscribe:%+v", req)
	// subscribe on pipe transport
	consumerOptions, err := h.ConsumerOptions(req.StreamId, req.RtpCapabilities)
	if err != nil {
		h.logger.Error().Msgf("create consumer options failed:%v", err)
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
		h.logger.Error().Msgf("subscribe on pipe transport failed:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Debug().Msgf("subscribe on pipe-transport done, data:%+v", consumerData)
	var rsp signaldefine.SubscribeOnPipeTransportResponse
	rsp.SubscribeId = consumerId
	rsp.SubscribeAnswer.Kind = consumerOptions.Kind
	rsp.SubscribeAnswer.RtpParameters = consumerOptions.RtpParameters
	_ = json.NewEncoder(w).Encode(rsp)
}
