package workerapi

import (
	"encoding/json"
	logger2 "github.com/byyam/mediasoup-go-worker/utils"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

var (
	logger = logger2.NewLogger("worker-api")
)

func CreateRouter(w *mediasoup_go_worker.SimpleWorker, routerId string) error {
	_, err := request(w, mediasoupdata.MethodWorkerCreateRouter, workerchannel.InternalData{RouterId: routerId})
	if err != nil {
		return err
	}
	return nil
}

type ParamCreateWebRtcTransport struct {
	RouterId    string
	TransportId string
	Options     mediasoupdata.WebRtcTransportOptions
}

func CreateWebRtcTransport(w *mediasoup_go_worker.SimpleWorker, param ParamCreateWebRtcTransport) (*mediasoupdata.WebrtcTransportData, error) {
	rsp, err := request(w, mediasoupdata.MethodRouterCreateWebRtcTransport, workerchannel.InternalData{
		RouterId:    param.RouterId,
		TransportId: param.TransportId,
	}, param.Options)
	if err != nil {
		return nil, err
	}
	var rspData mediasoupdata.WebrtcTransportData
	if err := json.Unmarshal(rsp.Data, &rspData); err != nil {
		return nil, err
	}
	return &rspData, nil
}

type ParamTransportConnect struct {
	RouterId    string
	TransportId string
	Options     mediasoupdata.TransportConnectOptions
}

func TransportConnect(w *mediasoup_go_worker.SimpleWorker, param ParamTransportConnect) error {
	rsp, err := request(w, mediasoupdata.MethodTransportConnect, workerchannel.InternalData{
		RouterId:    param.RouterId,
		TransportId: param.TransportId,
	}, param.Options)
	if err != nil {
		return err
	}
	var rspData mediasoupdata.TransportConnectData
	if err := json.Unmarshal(rsp.Data, &rspData); err != nil {
		return err
	}
	return nil
}

type ParamTransportProduce struct {
	RouterId    string
	TransportId string
	ProducerId  string
	Options     mediasoupdata.ProducerOptions
}

func TransportProduce(w *mediasoup_go_worker.SimpleWorker, param ParamTransportProduce) error {
	rsp, err := request(w, mediasoupdata.MethodTransportProduce, workerchannel.InternalData{
		RouterId:    param.RouterId,
		TransportId: param.TransportId,
		ProducerId:  param.ProducerId,
	}, param.Options)
	if err != nil {
		return err
	}
	var rspData mediasoupdata.ProducerData
	if err := json.Unmarshal(rsp.Data, &rspData); err != nil {
		return err
	}
	return nil
}

type ParamTransportConsume struct {
	RouterId    string
	TransportId string
	ProducerId  string
	ConsumerId  string
	Options     mediasoupdata.ConsumerOptions
}

func TransportConsume(w *mediasoup_go_worker.SimpleWorker, param ParamTransportConsume) (*mediasoupdata.ConsumerData, error) {
	rsp, err := request(w, mediasoupdata.MethodTransportConsume, workerchannel.InternalData{
		RouterId:    param.RouterId,
		TransportId: param.TransportId,
		ProducerId:  param.ProducerId,
		ConsumerId:  param.ConsumerId,
	}, param.Options)
	if err != nil {
		return nil, err
	}
	var rspData mediasoupdata.ConsumerData
	if err := json.Unmarshal(rsp.Data, &rspData); err != nil {
		return nil, err
	}
	return &rspData, nil
}

func ProducerClose(w *mediasoup_go_worker.SimpleWorker, internal workerchannel.InternalData) error {
	_, err := request(w, mediasoupdata.MethodProducerClose, internal)
	if err != nil {
		return err
	}
	return nil
}

func ProducerDump(w *mediasoup_go_worker.SimpleWorker, internal workerchannel.InternalData) (*mediasoupdata.ProducerDump, error) {
	rsp, err := request(w, mediasoupdata.MethodProducerDump, internal)
	if err != nil {
		return nil, err
	}
	var rspData mediasoupdata.ProducerDump
	if err := json.Unmarshal(rsp.Data, &rspData); err != nil {
		return nil, err
	}
	return &rspData, nil
}

func RouterDump(w *mediasoup_go_worker.SimpleWorker, internal workerchannel.InternalData) (*mediasoupdata.RouterDump, error) {
	rsp, err := request(w, mediasoupdata.MethodRouterDump, internal)
	if err != nil {
		return nil, err
	}
	var rspData mediasoupdata.RouterDump
	if err := json.Unmarshal(rsp.Data, &rspData); err != nil {
		return nil, err
	}
	return &rspData, nil
}

func TransportDump(w *mediasoup_go_worker.SimpleWorker, internal workerchannel.InternalData) (*mediasoupdata.TransportDump, error) {
	rsp, err := request(w, mediasoupdata.MethodTransportDump, internal)
	if err != nil {
		return nil, err
	}
	var rspData mediasoupdata.TransportDump
	if err := json.Unmarshal(rsp.Data, &rspData); err != nil {
		return nil, err
	}
	return &rspData, nil
}

func ConsumerClose(w *mediasoup_go_worker.SimpleWorker, internal workerchannel.InternalData) error {
	_, err := request(w, mediasoupdata.MethodConsumerClose, internal)
	if err != nil {
		return err
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// pipe-transport api

type ParamCreatePipeTransport struct {
	RouterId    string
	TransportId string
	Options     mediasoupdata.PipeTransportOptions
}

func CreatePipeTransport(w *mediasoup_go_worker.SimpleWorker, param ParamCreatePipeTransport) (*mediasoupdata.PipeTransportData, error) {
	rsp, err := request(w, mediasoupdata.MethodRouterCreatePipeTransport, workerchannel.InternalData{
		RouterId:    param.RouterId,
		TransportId: param.TransportId,
	}, param.Options)
	if err != nil {
		return nil, err
	}
	var rspData mediasoupdata.PipeTransportData
	if err := json.Unmarshal(rsp.Data, &rspData); err != nil {
		return nil, err
	}
	return &rspData, nil
}
