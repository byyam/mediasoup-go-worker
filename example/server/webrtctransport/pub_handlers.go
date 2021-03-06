package webrtctransport

import (
	"encoding/json"
	"github.com/byyam/mediasoup-go-worker/example/internal/demoutils"
	"github.com/byyam/mediasoup-go-worker/example/server/basehandler"
	"github.com/byyam/mediasoup-go-worker/workerchannel"

	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/server/democonf"
	"github.com/byyam/mediasoup-go-worker/example/server/workerapi"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/google/uuid"
	"github.com/jiyeyuran/go-protoo"
)

func (h *Handler) newTransport(dtlsParameters mediasoupdata.DtlsParameters, transportId string) (*mediasoupdata.WebrtcTransportData, error) {
	// create transport
	listenIp := mediasoupdata.TransportListenIp{
		Ip:          conf.Settings.RtcListenIp,
		AnnouncedIp: conf.Settings.RtcListenIp,
	}
	webrtcTransportOptions := democonf.WebrtcTransportOptions
	webrtcTransportOptions.ListenIps = append(webrtcTransportOptions.ListenIps, listenIp)
	transportData, err := workerapi.CreateWebRtcTransport(h.Worker, workerapi.ParamCreateWebRtcTransport{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportId,
		Options:     democonf.WebrtcTransportOptions,
	})
	if err != nil {
		return nil, err
	}
	h.logger.Debug("transport data:%+v", transportData)
	// connect transport
	transportConnectOptions := mediasoupdata.TransportConnectOptions{
		DtlsParameters: &dtlsParameters,
	}
	if err := workerapi.TransportConnect(h.Worker, workerapi.ParamTransportConnect{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportId,
		Options:     transportConnectOptions,
	}); err != nil {
		return nil, err
	}
	return transportData, nil
}

func (h *Handler) publishHandler(message protoo.Message) (interface{}, *protoo.Error) {
	var req isignal.PublishRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	transportId := uuid.New().String()
	transportData, err := h.newTransport(req.Offer.DtlsParameters, transportId)
	if err != nil {
		h.logger.Error("new transport on publish failed:%v", err)
		return nil, demoutils.ServerError(err)
	}
	// produce
	produceOptions, err := basehandler.ProducerOptions(req.Kind, req.StreamId, req.RtpParameters)
	if err != nil {
		return nil, demoutils.ServerError(err)
	}
	if err := workerapi.TransportProduce(h.Worker, workerapi.ParamTransportProduce{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportId,
		ProducerId:  produceOptions.Id,
		Options:     *produceOptions,
	}); err != nil {
		return nil, demoutils.ServerError(err)
	}

	return isignal.PublishResponse{
		TransportId: transportId,
		Answer: isignal.WebRtcTransportAnswer{
			IceParameters:  transportData.IceParameters,
			IceCandidates:  transportData.IceCandidates,
			DtlsParameters: transportData.DtlsParameters,
		},
	}, nil
}

func (h *Handler) unPublishHandler(message protoo.Message) (interface{}, *protoo.Error) {
	var req isignal.UnPublishRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	// get producer data
	transportData, producerData, err := h.FindProducer(demoutils.GetProducerId(req.StreamId))
	if err != nil {
		return nil, demoutils.ServerError(err)
	}

	if err := workerapi.ProducerClose(h.Worker, workerchannel.InternalData{
		RouterId:    demoutils.GetRouterId(h.Worker),
		TransportId: transportData.Id,
		ProducerId:  producerData.Id,
	}); err != nil {
		return nil, demoutils.ServerError(err)
	}
	return nil, nil
}
