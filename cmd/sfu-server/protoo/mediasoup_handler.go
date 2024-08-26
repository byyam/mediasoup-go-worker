package protoo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jiyeyuran/go-protoo"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/basehandler"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/democonf"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/demoutils"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/room"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/workerapi"
	"github.com/byyam/mediasoup-go-worker/conf"
	FBS__Router "github.com/byyam/mediasoup-go-worker/fbs/FBS/Router"
	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
	FBS__SctpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/SctpParameters"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/signaldefine"
)

func (h *ProtooHandler) GetRouterRtpCapabilities(message protoo.Message) (interface{}, *protoo.Error) {
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

func (h *ProtooHandler) CreateWebRtcTransport(message protoo.Message) (interface{}, *protoo.Error) {
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
	options.Base.NumSctpStreams = &FBS__SctpParameters.NumSctpStreamsT{
		Os:  req.SctpCapabilities.NumStreams.OS,
		Mis: req.SctpCapabilities.NumStreams.MIS,
	}
	rspData, err := workerapi.CreateWebRtcTransport(h.Worker, h.queryParams.RouterId, &FBS__Router.CreateWebRtcTransportRequestT{
		TransportId: transportId,
		Options:     &options,
	})
	if err != nil {
		return nil, demoutils.ServerError(err)
	}
	h.logger.Debug().Msgf("transport data:%+v", rspData)

	// update peer info
	getPeer, err := room.GetPeer(h.queryParams.RoomId, h.queryParams.PeerId)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetPeer failed")
		return nil, demoutils.ServerError(err)
	}
	if req.Producing {
		getPeer.ProducingWebRtcTransportId = transportId
	} else if req.Consuming {
		getPeer.ConsumingWebRtcTransportId = transportId
	}

	rsp := signaldefine.CreateWebRtcTransportResponse{
		Id:             transportId,
		IceParameters:  &mediasoupdata.IceParameters{},
		IceCandidates:  make([]*mediasoupdata.IceCandidate, 0),
		DtlsParameters: &mediasoupdata.DtlsParameters{},
		SctpParameters: &mediasoupdata.SctpParameters{},
	}
	if rspData.DtlsParameters != nil {
		rsp.DtlsParameters.Set(rspData.DtlsParameters)
	}
	if rspData.IceParameters != nil {
		rsp.IceParameters.Set(rspData.IceParameters)
	}
	if rspData.Base != nil && rspData.Base.SctpParameters != nil {
		rsp.SctpParameters.Set(rspData.Base.SctpParameters)
	}
	for _, iceCandidate := range rspData.IceCandidates {
		candidate := &mediasoupdata.IceCandidate{}
		candidate.Set(iceCandidate)
		rsp.IceCandidates = append(rsp.IceCandidates, candidate)
	}

	return rsp, nil
}

func (h *ProtooHandler) Join(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.JoinRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}

	getRoom, err := room.GetRoom(h.queryParams.RoomId)
	if err != nil {
		return nil, demoutils.ServerError(err)
	}

	getPeer, exist := getRoom.Peers[h.queryParams.PeerId]
	if !exist {
		return nil, demoutils.ServerError(fmt.Errorf("peer %s not exist", h.queryParams.PeerId))
	}
	getPeer.Joined = true
	// produce

	for id, otherPeer := range getRoom.Peers {
		if id == h.queryParams.PeerId {
			continue
		}
		if !otherPeer.Joined {
			continue
		}
		// consume

	}

	return signaldefine.JoinResponse{
		Peers: make([]signaldefine.JoinedPeerInfo, 0),
	}, nil
}

func (h *ProtooHandler) ConnectWebRtcTransport(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.ConnectWebRtcTransportRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	// connect webrtc-transport
	rsp, err := workerapi.ConnectWebRtcTransport(h.Worker,
		h.queryParams.RouterId,
		req.TransportId,
		&FBS__WebRtcTransport.ConnectRequestT{DtlsParameters: req.DtlsParameters.Convert()})
	if err != nil {
		h.logger.Error().Err(err).Msg("ConnectWebRtcTransport failed")
		return nil, demoutils.ServerError(err)
	}

	return signaldefine.ConnectWebRtcTransportResponse{
		DtlsRole: strings.ToLower(FBS__WebRtcTransport.EnumNamesDtlsRole[rsp.DtlsLocalRole]),
	}, nil
}

func (h *ProtooHandler) Produce(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.ProduceRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}
	produceId := uuid.New().String()
	produceOptions, err := basehandler.ProduceOptions(req.Kind, produceId, req.RtpParameters)
	if err != nil {
		return nil, demoutils.ServerError(err)
	}
	content, _ := json.Marshal(produceOptions)
	h.logger.Debug().Msgf("[Produce] produceOptions:%s", string(content))
	_, err = workerapi.TransportProduce(h.Worker,
		h.queryParams.RouterId,
		req.TransportId,
		&FBS__Transport.ProduceRequestT{
			ProducerId:    produceId,
			Kind:          FBS__RtpParameters.EnumValuesMediaKind[strings.ToUpper(string(req.Kind))],
			RtpParameters: req.RtpParameters.Convert(),
			RtpMapping:    produceOptions.RtpMapping.Convert(),
		},
	)
	if err != nil {
		h.logger.Error().Err(err).Msg("Produce failed")
		return nil, demoutils.ServerError(err)
	}

	return signaldefine.ProduceResponse{
		Id: produceId,
	}, nil
}

func (h *ProtooHandler) GetTransportStats(message protoo.Message) (interface{}, *protoo.Error) {
	var req signaldefine.GetTransportStatRequest
	if err := json.Unmarshal(message.Data, &req); err != nil {
		return nil, demoutils.ServerError(err)
	}

	rsp, err := workerapi.GetTransportStats(h.Worker,
		h.queryParams.RouterId,
		req.TransportId,
	)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetTransportStats failed")
		return nil, demoutils.ServerError(err)
	}

	return rsp, nil
}
