package workerapi

import (
	"encoding/json"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	FBS__PipeTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/PipeTransport"
	FBS__Producer "github.com/byyam/mediasoup-go-worker/fbs/FBS/Producer"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Router "github.com/byyam/mediasoup-go-worker/fbs/FBS/Router"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	FBS__Worker "github.com/byyam/mediasoup-go-worker/fbs/FBS/Worker"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/signaldefine"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

var (
	logger = zerowrapper.NewScope("worker-api")
)

func CreateRouter(w *mediasoup_go_worker.SimpleWorker, routerId string) error {
	_, err := requestFbs(w, workerchannel.InternalData{RouterId: routerId},
		&FBS__Request.RequestT{
			Id:        GetRid(),
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

func CreateWebRtcTransport(w *mediasoup_go_worker.SimpleWorker, routerId string, param *FBS__Router.CreateWebRtcTransportRequestT) (*FBS__WebRtcTransport.DumpResponseT, error) {
	logger.Info().Msgf("[CreateWebRtcTransport] routerId:%s, requestT:%+v", routerId, param)
	rsp, err := requestFbs(w, workerchannel.InternalData{
		RouterId:    routerId,
		TransportId: param.TransportId,
	}, &FBS__Request.RequestT{
		Id:        GetRid(),
		Method:    FBS__Request.MethodROUTER_CREATE_WEBRTCTRANSPORT,
		HandlerId: routerId,
		Body: &FBS__Request.BodyT{
			Type:  FBS__Request.BodyRouter_CreateWebRtcTransportRequest,
			Value: param,
		},
	})
	if err != nil {
		return nil, err
	}
	rspData := rsp.RspBody.Value.(*FBS__WebRtcTransport.DumpResponseT)

	return rspData, nil
}

func ConnectWebRtcTransport(w *mediasoup_go_worker.SimpleWorker, routerId, transportId string, param *FBS__WebRtcTransport.ConnectRequestT) (*FBS__WebRtcTransport.ConnectResponseT, error) {
	rsp, err := requestFbs(w, workerchannel.InternalData{
		RouterId:    routerId,
		TransportId: transportId,
	}, &FBS__Request.RequestT{
		Id:        GetRid(),
		Method:    FBS__Request.MethodWEBRTCTRANSPORT_CONNECT,
		HandlerId: transportId,
		Body: &FBS__Request.BodyT{
			Type:  FBS__Request.BodyWebRtcTransport_ConnectRequest,
			Value: param,
		},
	})
	if err != nil {
		return nil, err
	}
	rspData := rsp.RspBody.Value.(*FBS__WebRtcTransport.ConnectResponseT)

	return rspData, nil
}

func GetTransportStats(w *mediasoup_go_worker.SimpleWorker, routerId, transportId string) (signaldefine.GetTransportStatResponse, error) {
	rsp, err := requestFbs(w, workerchannel.InternalData{
		RouterId:    routerId,
		TransportId: transportId,
	}, &FBS__Request.RequestT{
		Id:        GetRid(),
		Method:    FBS__Request.MethodTRANSPORT_GET_STATS,
		HandlerId: transportId,
		Body:      nil,
	})
	if err != nil {
		return nil, err
	}
	logger.Info().Msgf("[GetTransportStats] rsp.Data:%s", rsp.Data)
	ret := make([]*signaldefine.TransportStatResponse, 0)
	switch rsp.RspBody.Value.(type) {
	case *FBS__WebRtcTransport.GetStatsResponseT:
		rspData := rsp.RspBody.Value.(*FBS__WebRtcTransport.GetStatsResponseT)
		logger.Info().Msgf("[GetTransportStats] rspData:%+v", rspData)
		stat := &signaldefine.TransportStatResponse{}
		stat.Set("webrtc-transport", rspData)
		ret = append(ret, stat)
	}

	return ret, nil
}

func GetProducerStats(w *mediasoup_go_worker.SimpleWorker, routerId, producerId string) (*FBS__Producer.GetStatsResponseT, error) {
	rsp, err := requestFbs(w, workerchannel.InternalData{
		RouterId:   routerId,
		ProducerId: producerId,
	}, &FBS__Request.RequestT{
		Id:        GetRid(),
		Method:    FBS__Request.MethodPRODUCER_GET_STATS,
		HandlerId: producerId,
		Body:      nil,
	})
	if err != nil {
		return nil, err
	}

	rspData := rsp.RspBody.Value.(*FBS__Producer.GetStatsResponseT)
	logger.Info().Msgf("[GetProducerStats] rsp.Data:%s", rsp.Data)

	return rspData, nil
}

func ConnectPipeTransport(w *mediasoup_go_worker.SimpleWorker, routerId, transportId string, param *FBS__PipeTransport.ConnectRequestT) error {
	// todo
	rsp, err := request(w, mediasoupdata.MethodTransportConnect, workerchannel.InternalData{
		RouterId:    routerId,
		TransportId: transportId,
	}, param)
	if err != nil {
		return err
	}
	logger.Info().Msgf("[ConnectPipeTransport] rsp.Data:%s", rsp.Data)
	var rspData FBS__PipeTransport.ConnectResponseT
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

func TransportProduce(w *mediasoup_go_worker.SimpleWorker, routerId, transportId string, param *FBS__Transport.ProduceRequestT) (*FBS__Transport.ProduceResponseT, error) {
	rsp, err := requestFbs(w, workerchannel.InternalData{
		RouterId:    routerId,
		TransportId: transportId,
		ProducerId:  param.ProducerId,
	}, &FBS__Request.RequestT{
		Id:        GetRid(),
		Method:    FBS__Request.MethodTRANSPORT_PRODUCE,
		HandlerId: transportId,
		Body: &FBS__Request.BodyT{
			Type:  FBS__Request.BodyTransport_ProduceRequest,
			Value: param,
		},
	})
	if err != nil {
		return nil, err
	}
	rspData := rsp.RspBody.Value.(*FBS__Transport.ProduceResponseT)

	return rspData, nil
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

func CloseRouter(w *mediasoup_go_worker.SimpleWorker, routerId string) error {
	_, err := requestFbs(w, workerchannel.InternalData{RouterId: routerId},
		&FBS__Request.RequestT{
			Id:        GetRid(),
			Method:    FBS__Request.MethodWORKER_CLOSE_ROUTER,
			HandlerId: routerId,
			Body: &FBS__Request.BodyT{
				Type:  FBS__Request.BodyWorker_CloseRouterRequest,
				Value: &FBS__Worker.CloseRouterRequestT{RouterId: routerId},
			},
		},
	)
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
