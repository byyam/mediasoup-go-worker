package workerapi

import (
	"encoding/json"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Worker "github.com/byyam/mediasoup-go-worker/fbs/FBS/Worker"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

var (
	logger = zerowrapper.NewScope("worker-api")
)

func CreateRouter(w *mediasoup_go_worker.SimpleWorker, routerId string) error {
	_, err := requestFbs(w, workerchannel.InternalData{RouterId: routerId},
		&FBS__Request.RequestT{
			Method:    FBS__Request.MethodWORKER_CREATE_ROUTER,
			HandlerId: routerId,
			Body: &FBS__Request.BodyT{
				Type:  FBS__Request.BodyWorker_CreateRouterRequest,
				Value: &FBS__Worker.CreateRouterRequestT{RouterId: routerId},
			},
		},
	)
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
