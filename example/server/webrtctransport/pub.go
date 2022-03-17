package webrtctransport

import (
	"encoding/json"

	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/server/config"
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
	webrtcTransportOptions := config.WebrtcTransportOptions
	webrtcTransportOptions.ListenIps = append(webrtcTransportOptions.ListenIps, listenIp)
	transportData, err := workerapi.CreateWebRtcTransport(h.worker, workerapi.ParamCreateWebRtcTransport{
		RouterId:    GetRouterId(h.worker),
		TransportId: transportId,
		Options:     config.WebrtcTransportOptions,
	})
	if err != nil {
		return nil, err
	}
	h.logger.Debug("transport data:%+v", transportData)
	// connect transport
	transportConnectOptions := mediasoupdata.TransportConnectOptions{
		DtlsParameters: &dtlsParameters,
	}
	if err := workerapi.TransportConnect(h.worker, workerapi.ParamTransportConnect{
		RouterId:    GetRouterId(h.worker),
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
		return nil, ServerError(err)
	}
	transportId := uuid.New().String()
	transportData, err := h.newTransport(req.Offer.DtlsParameters, transportId)
	if err != nil {
		h.logger.Error("new transport on publish failed:%v", err)
		return nil, ServerError(err)
	}
	// produce
	routerRtpCapabilities, err := mediasoupdata.GenerateRouterRtpCapabilities(config.RouterOptions.MediaCodecs)
	if err != nil {
		h.logger.Error("GenerateRouterRtpCapabilities failed:%+v", err)
		return nil, ServerError(err)
	}
	for _, c := range routerRtpCapabilities.Codecs {
		h.logger.Debug("routerRtpCapabilities:%+v", c)
	}
	rtpMapping, err := mediasoupdata.GetProducerRtpParametersMapping(
		req.RtpParameters, routerRtpCapabilities)
	if err != nil {
		h.logger.Error("GetProducerRtpParametersMapping:%+v", err)
		return nil, ServerError(err)
	}

	produceOptions := mediasoupdata.ProducerOptions{
		Id:                   GetProducerId(req.StreamId),
		Kind:                 req.Kind,
		RtpParameters:        req.RtpParameters,
		Paused:               false,
		KeyFrameRequestDelay: 0,
		AppData:              nil,
		RtpMapping:           rtpMapping,
	}
	if err := workerapi.TransportProduce(h.worker, workerapi.ParamTransportProduce{
		RouterId:    GetRouterId(h.worker),
		TransportId: transportId,
		ProducerId:  produceOptions.Id,
		Options:     produceOptions,
	}); err != nil {
		return nil, ServerError(err)
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

func (h *Handler) unPublishHandler(req protoo.Message) (interface{}, *protoo.Error) {
	return nil, nil
}
