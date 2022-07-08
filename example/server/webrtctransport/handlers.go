package webrtctransport

import (
	"errors"
	utils2 "github.com/byyam/mediasoup-go-worker/example/server/utils"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/server/workerapi"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
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
	case isignal.MethodUnPublish:
		data, err = h.unPublishHandler(req)
	case isignal.MethodSubscribe:
		data, err = h.subscribeHandler(req)
	case isignal.MethodUnSubscribe:
		data, err = h.unSubscribeHandler(req)
	default:
		err = utils2.ErrUnknownMethod
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

func (h *Handler) findProducer(targetId string) (*mediasoupdata.TransportDump, *mediasoupdata.ProducerDump, error) {
	routerDump, err := workerapi.RouterDump(h.worker, workerchannel.InternalData{RouterId: utils2.GetRouterId(h.worker)})
	if err != nil {
		return nil, nil, err
	}
	for _, transportId := range routerDump.TransportIds {
		transportDump, err := h.getTransportDump(transportId)
		if err != nil {
			return nil, nil, err
		}
		for _, producerId := range transportDump.ProducerIds {
			if targetId != producerId {
				continue
			}
			producerDump, err := h.getProducerDump(transportId, producerId)
			if err != nil {
				return nil, nil, err
			}
			return transportDump, producerDump, nil
		}
	}
	return nil, nil, errors.New("producer not found")
}

func (h *Handler) getTransportDump(transportId string) (*mediasoupdata.TransportDump, error) {
	transportDump, err := workerapi.TransportDump(h.worker, workerchannel.InternalData{
		RouterId:    utils2.GetRouterId(h.worker),
		TransportId: transportId,
	})
	return transportDump, err
}

func (h *Handler) getProducerDump(transportId, producerId string) (*mediasoupdata.ProducerDump, error) {
	producerDump, err := workerapi.ProducerDump(h.worker, workerchannel.InternalData{
		RouterId:    utils2.GetRouterId(h.worker),
		TransportId: transportId,
		ProducerId:  producerId,
	})
	return producerDump, err
}
