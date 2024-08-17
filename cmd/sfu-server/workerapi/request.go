package workerapi

import (
	"encoding/json"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"

	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

func request(worker *mediasoup_go_worker.SimpleWorker, method string, internal workerchannel.InternalData, data ...interface{}) (*workerchannel.ResponseData, error) {
	req := workerchannel.RequestData{
		Method:   method,
		Internal: internal,
		Request:  new(FBS__Request.RequestT),
	}
	if len(data) > 0 {
		rawData, err := json.Marshal(data[0])
		if err != nil {
			return nil, err
		}
		//req.Data = rawData
		if err = json.Unmarshal(rawData, req.Request); err != nil {
			return nil, err
		}
		req.Method = FBS__Request.EnumNamesMethod[req.Request.Method]
	}
	rsp := worker.OnChannelRequest(req)
	logger.Info().Msgf("request done, req:[%+v] rsp:[%+v]", req, rsp)
	if rsp.Err != nil {
		return nil, rsp.Err
	}
	return &rsp, nil
}

func requestFbs(worker *mediasoup_go_worker.SimpleWorker, internal workerchannel.InternalData, fbsRequest *FBS__Request.RequestT) (*workerchannel.ResponseData, error) {
	req := workerchannel.RequestData{
		MethodType: fbsRequest.Method,
		Method:     FBS__Request.EnumNamesMethod[fbsRequest.Method],
		HandlerId:  fbsRequest.HandlerId, // set handlerID
		Internal:   internal,
		Request:    fbsRequest,
	}

	// for logging
	reqData, _ := json.Marshal(req.Request.Body)
	logger.Info().Msgf("[requestFbs] request --> \n[%+v] \nreqData: \n[%s]",
		req, string(reqData))

	rsp := worker.OnChannelRequest(req)

	// for logging
	if rsp.Data == nil {
		rsp.Data, _ = json.Marshal(rsp.RspBody)
	}
	logger.Info().Msgf("[requestFbs] response <-- \n[%+v] \nrspData: \n[%s]",
		rsp, string(rsp.Data))

	if rsp.Err != nil {
		return nil, rsp.Err
	}
	return &rsp, nil
}
