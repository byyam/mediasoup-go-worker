package workerapi

import (
	"encoding/json"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"

	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

func request(worker *mediasoup_go_worker.SimpleWorker, method string, internal workerchannel.InternalData, data ...interface{}) (*workerchannel.ResponseData, error) {
	req := workerchannel.RequestData{
		Method:   method,
		Internal: internal,
	}
	if len(data) > 0 {
		rawData, err := json.Marshal(data[0])
		if err != nil {
			return nil, err
		}
		req.Data = rawData
	}
	rsp := worker.OnChannelRequest(req)
	logger.Info().Msgf("request done, req:[%+v] rsp:[%+v]", req, rsp)
	if rsp.Err != nil {
		return nil, rsp.Err
	}
	return &rsp, nil
}
