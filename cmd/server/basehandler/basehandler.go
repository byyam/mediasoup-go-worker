package basehandler

import (
	"errors"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"

	"github.com/byyam/mediasoup-go-worker/cmd/server/democonf"
	"github.com/byyam/mediasoup-go-worker/cmd/server/demoutils"
	"github.com/byyam/mediasoup-go-worker/cmd/server/workerapi"
)

var (
	logger = zerowrapper.NewScope("base-handler")
)

type BaseHandler struct {
	Worker *mediasoup_go_worker.SimpleWorker
}

func (h *BaseHandler) ConsumerOptions(streamId uint64, capabilities *mediasoupdata.RtpCapabilities) (*mediasoupdata.ConsumerOptions, error) {
	// get producer data
	_, producerData, err := h.FindProducer(demoutils.GetProducerId(streamId))
	if err != nil {
		return nil, demoutils.ServerError(err)
	}

	routerRtpCapabilities, err := mediasoupdata.GenerateRouterRtpCapabilities(democonf.RouterOptions.MediaCodecs)
	if err != nil {
		logger.Error().Msgf("GenerateRouterRtpCapabilities failed:%+v", err)
		return nil, demoutils.ServerError(err)
	}
	consumableRtpParameters, err := mediasoupdata.GetConsumableRtpParameters(
		mediasoupdata.MediaKind(producerData.Kind), producerData.RtpParameters, routerRtpCapabilities, producerData.RtpMapping)
	if err != nil {
		return nil, demoutils.ServerError(err)
	}

	rtpParameters, err := mediasoupdata.GetConsumerRtpParameters(consumableRtpParameters, *capabilities, false)
	if err != nil {
		return nil, demoutils.ServerError(err)
	}

	// consume
	consumerOptions := &mediasoupdata.ConsumerOptions{
		ProducerId:             demoutils.GetProducerId(streamId),
		RtpCapabilities:        mediasoupdata.RtpCapabilities{},
		Paused:                 false,
		Mid:                    "",
		PreferredLayers:        nil,
		Pipe:                   false,
		AppData:                nil,
		Kind:                   mediasoupdata.MediaKind(producerData.Kind),
		Type:                   mediasoupdata.ConsumerType(producerData.Type),
		RtpParameters:          rtpParameters,
		ConsumableRtpEncodings: consumableRtpParameters.Encodings,
	}
	return consumerOptions, nil
}

func (h *BaseHandler) FindProducer(targetId string) (*mediasoupdata.TransportDump, *mediasoupdata.ProducerDump, error) {
	routerDump, err := workerapi.RouterDump(h.Worker, workerchannel.InternalData{RouterId: demoutils.GetRouterId(h.Worker)})
	if err != nil {
		return nil, nil, err
	}
	for _, transportId := range routerDump.TransportIds {
		transportDump, err := h.GetTransportDump(transportId)
		if err != nil {
			return nil, nil, err
		}
		for _, producerId := range transportDump.ProducerIds {
			if targetId != producerId {
				continue
			}
			producerDump, err := h.GetProducerDump(transportId, producerId)
			if err != nil {
				return nil, nil, err
			}
			return transportDump, producerDump, nil
		}
	}
	return nil, nil, errors.New("producer not found")
}

func (h *BaseHandler) GetTransportDump(transportId string) (*mediasoupdata.TransportDump, error) {
	transportDump, err := workerapi.TransportDump(h.Worker, workerchannel.InternalData{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportId,
	})
	return transportDump, err
}

func (h *BaseHandler) GetProducerDump(transportId, producerId string) (*mediasoupdata.ProducerDump, error) {
	producerDump, err := workerapi.ProducerDump(h.Worker, workerchannel.InternalData{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportId,
		ProducerId:  producerId,
	})
	return producerDump, err
}
