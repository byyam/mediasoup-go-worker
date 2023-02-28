package webrtctransport

import (
	"encoding/json"

	"github.com/jiyeyuran/go-protoo"

	"github.com/byyam/mediasoup-go-worker/workerchannel"

	"github.com/google/uuid"

	"github.com/byyam/mediasoup-go-worker/signaldefine"

	"github.com/byyam/mediasoup-go-worker/cmd/server/demoutils"
	"github.com/byyam/mediasoup-go-worker/cmd/server/workerapi"
)

func (h *Handler) subscribeHandler(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.SubscribeRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	transportId := uuid.New().String()
	transportData, err := h.newTransport(req.Offer.DtlsParameters, transportId)
	if err != nil {
		h.logger.Error().Msgf("new transport on publish failed:%v", err)
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
	h.logger.Info().Msgf("consumer data:%+v", consumerData)

	return signaldefine.SubscribeResponse{
		SubscribeId: consumerId,
		TransportId: transportId,
		Answer: signaldefine.WebRtcTransportAnswer{
			IceParameters:  transportData.IceParameters,
			IceCandidates:  transportData.IceCandidates,
			DtlsParameters: transportData.DtlsParameters,
		},
		SubscribeAnswer: signaldefine.SubscribeAnswer{
			Kind:          consumerOptions.Kind,
			RtpParameters: consumerOptions.RtpParameters,
			AppData:       nil,
		},
	}, nil
}

func (h *Handler) unSubscribeHandler(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.UnSubscribeRequest
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
