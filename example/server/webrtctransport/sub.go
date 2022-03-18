package webrtctransport

import (
	"encoding/json"

	"github.com/byyam/mediasoup-go-worker/workerchannel"

	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/server/democonf"
	"github.com/byyam/mediasoup-go-worker/example/server/workerapi"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/google/uuid"
	"github.com/jiyeyuran/go-protoo"
)

func (h *Handler) subscribeHandler(message protoo.Message) (interface{}, *protoo.Error) {
	var req isignal.SubscribeRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, ServerError(err)
	}
	transportId := uuid.New().String()
	transportData, err := h.newTransport(req.Offer.DtlsParameters, transportId)
	if err != nil {
		h.logger.Error("new transport on publish failed:%v", err)
		return nil, ServerError(err)
	}
	// get producer data
	_, producerData, err := h.findProducer(GetProducerId(req.StreamId))
	if err != nil {
		return nil, ServerError(err)
	}

	routerRtpCapabilities, err := mediasoupdata.GenerateRouterRtpCapabilities(democonf.RouterOptions.MediaCodecs)
	if err != nil {
		h.logger.Error("GenerateRouterRtpCapabilities failed:%+v", err)
		return nil, ServerError(err)
	}
	consumableRtpParameters, err := mediasoupdata.GetConsumableRtpParameters(
		mediasoupdata.MediaKind(producerData.Kind), producerData.RtpParameters, routerRtpCapabilities, producerData.RtpMapping)
	if err != nil {
		return nil, ServerError(err)
	}

	rtpParameters, err := mediasoupdata.GetConsumerRtpParameters(consumableRtpParameters, *req.RtpCapabilities, false)
	if err != nil {
		return nil, ServerError(err)
	}

	// consume
	consumerId := uuid.New().String()
	consumerOptions := mediasoupdata.ConsumerOptions{
		ProducerId:             GetProducerId(req.StreamId),
		RtpCapabilities:        mediasoupdata.RtpCapabilities{},
		Paused:                 false,
		Mid:                    "",
		PreferredLayers:        nil,
		Pipe:                   false,
		AppData:                nil,
		Kind:                   mediasoupdata.MediaKind(producerData.Kind),
		Type:                   mediasoupdata.ConsumerType(producerData.Type),
		RtpParameters:          rtpParameters,
		ConsumableRtpEncodings: nil,
	}
	consumerData, err := workerapi.TransportConsume(h.worker, workerapi.ParamTransportConsume{
		RouterId:    GetRouterId(h.worker),
		TransportId: transportId,
		ProducerId:  consumerOptions.ProducerId,
		ConsumerId:  consumerId,
		Options:     consumerOptions,
	})
	if err != nil {
		return nil, ServerError(err)
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
		Kind:          mediasoupdata.MediaKind(producerData.Kind),
		RtpParameters: consumerOptions.RtpParameters,
		AppData:       nil,
	}, nil
}

func (h *Handler) unSubscribeHandler(message protoo.Message) (interface{}, *protoo.Error) {
	var req isignal.UnSubscribeRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, ServerError(err)
	}
	// get producer data
	transportData, producerData, err := h.findProducer(GetProducerId(req.StreamId))
	if err != nil {
		return nil, ServerError(err)
	}

	if err := workerapi.ConsumerClose(h.worker, workerchannel.InternalData{
		RouterId:    GetRouterId(h.worker),
		TransportId: transportData.Id,
		ProducerId:  producerData.Id,
		ConsumerId:  req.SubscribeId,
	}); err != nil {
		return nil, ServerError(err)
	}
	return nil, nil
}
