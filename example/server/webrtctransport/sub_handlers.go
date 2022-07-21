package webrtctransport

import (
	"encoding/json"
	"github.com/byyam/mediasoup-go-worker/example/internal/demoutils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"

	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/server/workerapi"
	"github.com/google/uuid"
	"github.com/jiyeyuran/go-protoo"
)

func (h *Handler) subscribeHandler(message protoo.Message) (interface{}, *protoo.Error) {
	var req isignal.SubscribeRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	transportId := uuid.New().String()
	transportData, err := h.newTransport(req.Offer.DtlsParameters, transportId)
	if err != nil {
		h.logger.Error("new transport on publish failed:%v", err)
		return nil, demoutils.ServerError(err)
	}
	consumerId := uuid.New().String()
	consumerOptions, err := h.ConsumerOptions(req.StreamId, req.RtpCapabilities)
	if err != nil {
		return nil, demoutils.ServerError(err)
	}
	consumerData, err := workerapi.TransportConsume(h.Worker, workerapi.ParamTransportConsume{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportId,
		ProducerId:  consumerOptions.ProducerId,
		ConsumerId:  consumerId,
		Options:     *consumerOptions,
	})
	if err != nil {
		return nil, demoutils.ServerError(err)
	}
	h.logger.Info("consumer data:%+v", consumerData)

	return isignal.SubscribeResponse{
		SubscribeId: consumerId,
		TransportId: transportId,
		Answer: isignal.WebRtcTransportAnswer{
			IceParameters:  transportData.IceParameters,
			IceCandidates:  transportData.IceCandidates,
			DtlsParameters: transportData.DtlsParameters,
		},
		SubscribeAnswer: isignal.SubscribeAnswer{
			Kind:          consumerOptions.Kind,
			RtpParameters: consumerOptions.RtpParameters,
			AppData:       nil,
		},
	}, nil
}

func (h *Handler) unSubscribeHandler(message protoo.Message) (interface{}, *protoo.Error) {
	var req isignal.UnSubscribeRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	// get producer data
	transportData, producerData, err := h.FindProducer(demoutils.GetProducerId(req.StreamId))
	if err != nil {
		return nil, demoutils.ServerError(err)
	}

	if err := workerapi.ConsumerClose(h.Worker, workerchannel.InternalData{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportData.Id,
		ProducerId:  producerData.Id,
		ConsumerId:  req.SubscribeId,
	}); err != nil {
		return nil, demoutils.ServerError(err)
	}
	return nil, nil
}
