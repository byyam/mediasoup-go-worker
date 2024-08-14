package webrtctransport

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jiyeyuran/go-protoo"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/democonf"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/demoutils"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/workerapi"
	"github.com/byyam/mediasoup-go-worker/conf"
	FBS__Router "github.com/byyam/mediasoup-go-worker/fbs/FBS/Router"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/signaldefine"
)

func (h *Handler) GetRouterRtpCapabilities(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.GetRouterRtpCapabilitiesRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	routerRtpCapabilities, err := mediasoupdata.GenerateRouterRtpCapabilities(democonf.RouterOptions.MediaCodecs)
	if err != nil {
		h.logger.Error().Msgf("GenerateRouterRtpCapabilities failed:%+v", err)
		return nil, demoutils.ServerError(err)
	}

	return signaldefine.GetRouterRtpCapabilitiesResponse{
		RtpCapabilities: routerRtpCapabilities,
	}, nil
}

func (h *Handler) CreateWebRtcTransport(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.CreateWebRtcTransportRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	// create webrtc transport
	transportId := uuid.New().String()
	options := democonf.WebrtcTransportOptionsFBS
	listenIp := &FBS__WebRtcTransport.ListenT{
		Type: FBS__WebRtcTransport.ListenListenIndividual,
		Value: &FBS__Transport.ListenInfoT{
			Protocol:       FBS__Transport.ProtocolUDP,
			Ip:             conf.Settings.RtcListenIp,
			AnnouncedIp:    conf.Settings.RtcListenIp,
			Port:           0,
			Flags:          nil,
			SendBufferSize: 0,
			RecvBufferSize: 0,
		},
	}
	options.Listen = listenIp
	transportData, err := workerapi.CreateWebRtcTransport(h.Worker, demoutils.GetRouterId(h.Worker), &FBS__Router.CreateWebRtcTransportRequestT{
		TransportId: transportId,
		Options:     &options,
	})
	if err != nil {
		return nil, demoutils.ServerError(err)
	}
	h.logger.Debug().Msgf("transport data:%+v", transportData)
	return signaldefine.CreateWebRtcTransportResponse{
		Id:             transportId,
		IceCandidates:  transportData.IceCandidates,
		IceParameters:  transportData.IceParameters,
		DtlsParameters: transportData.DtlsParameters,
		SctpParameters: transportData.SctpParameters,
	}, nil
}

func (h *Handler) Join(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.JoinRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}

	return signaldefine.JoinResponse{}, nil
}

func (h *Handler) ConnectWebRtcTransport(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.ConnectWebRtcTransportRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}

	return signaldefine.ConnectWebRtcTransportResponse{}, nil
}
