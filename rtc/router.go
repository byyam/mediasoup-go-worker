package rtc

import (
	"encoding/json"
	"sync"

	"github.com/byyam/mediasoup-go-worker/common"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"

	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type Router struct {
	id           string
	logger       utils.Logger
	transportMap sync.Map
}

func NewRouter(id string) *Router {
	return &Router{
		id:     id,
		logger: utils.NewLogger("router"),
	}
}

func (r *Router) HandleRequest(request workerchannel.RequestData) (response workerchannel.ResponseData) {

	switch request.Method {
	case mediasoupdata.MethodRouterCreateWebRtcTransport:
		var options mediasoupdata.WebRtcTransportOptions
		_ = json.Unmarshal(request.Data, &options)
		webrtcTransport, err := newWebrtcTransport(request.InternalData.TransportId, options)
		if err != nil {
			response.Err = common.ErrCreateWebrtcTransport
			return
		}
		r.transportMap.Store(request.InternalData.TransportId, webrtcTransport)
		response.Data = webrtcTransport.FillJson()

	case mediasoupdata.MethodRouterCreatePlainTransport:

	case mediasoupdata.MethodRouterCreatePipeTransport:

	case mediasoupdata.MethodRouterCreateDirectTransport:

	case mediasoupdata.MethodRouterCreateActiveSpeakerObserver:

	case mediasoupdata.MethodRouterCreateAudioLevelObserver:

	case mediasoupdata.MethodRouterDump:
		response.Data = r.FillJson()

	case mediasoupdata.MethodRouterClose:
		r.Close()
	default:
		r, ok := r.transportMap.Load(request.InternalData.TransportId)
		if !ok {
			response.Err = common.ErrTransportNotFound
			return
		}
		transport := r.(ITransport)
		response = transport.HandleRequest(request)
	}
	return
}

func (r *Router) Close() {
	r.transportMap.Range(func(key, value interface{}) bool {
		transport := value.(ITransport)
		transport.Close()
		return true
	})
	r.logger.Warn("router:%s stop", r.id)
}

func (r *Router) FillJson() json.RawMessage {
	var transportIds []string
	r.transportMap.Range(func(key, value interface{}) bool {
		transportIds = append(transportIds, key.(string))
		return true
	})
	dumpData := mediasoupdata.RouterDump{
		Id:                               r.id,
		TransportIds:                     transportIds,
		RtpObserverIds:                   nil,
		MapProducerIdConsumerIds:         nil,
		MapConsumerIdProducerId:          nil,
		MapProducerIdObserverIds:         nil,
		MapDataProducerIdDataConsumerIds: nil,
		MapDataConsumerIdDataProducerId:  nil,
	}
	data, _ := json.Marshal(&dumpData)
	r.logger.Debug("dumpData:%+v", dumpData)
	return data
}
