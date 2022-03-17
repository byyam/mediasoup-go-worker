package webrtctransport

import (
	"encoding/json"
	"errors"

	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/server/config"
	"github.com/byyam/mediasoup-go-worker/example/server/workerapi"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/google/uuid"
	"github.com/jiyeyuran/go-protoo"
)

func (h *Handler) findProducerDump(targetId string) (*mediasoupdata.ProducerDump, error) {
	routerDump, err := workerapi.RouterDump(h.worker, workerchannel.InternalData{RouterId: GetRouterId(h.worker)})
	if err != nil {
		return nil, err
	}
	for _, transportId := range routerDump.TransportIds {
		transportDump, err := h.getTransportDump(transportId)
		if err != nil {
			return nil, err
		}
		for _, producerId := range transportDump.ProducerIds {
			if targetId != producerId {
				continue
			}
			producerDump, err := h.getProducerDump(transportId, producerId)
			if err != nil {
				return nil, err
			}
			return producerDump, nil
		}
	}
	return nil, errors.New("producer dump not found")
}

func (h *Handler) getTransportDump(transportId string) (*mediasoupdata.TransportDump, error) {
	transportDump, err := workerapi.TransportDump(h.worker, workerchannel.InternalData{
		RouterId:    GetRouterId(h.worker),
		TransportId: transportId,
	})
	return transportDump, err
}

func (h *Handler) getProducerDump(transportId, producerId string) (*mediasoupdata.ProducerDump, error) {
	producerDump, err := workerapi.ProducerDump(h.worker, workerchannel.InternalData{
		RouterId:    GetRouterId(h.worker),
		TransportId: transportId,
		ProducerId:  producerId,
	})
	return producerDump, err
}

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
	producerData, err := h.findProducerDump(GetProducerId(req.StreamId))
	if err != nil {
		return nil, ServerError(err)
	}

	routerRtpCapabilities, err := mediasoupdata.GenerateRouterRtpCapabilities(config.RouterOptions.MediaCodecs)
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
