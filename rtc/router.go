package rtc

import (
	"encoding/json"
	"sync"

	"github.com/pion/rtp"

	"github.com/byyam/mediasoup-go-worker/common"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type Router struct {
	id                   string
	logger               utils.Logger
	mapTransports        sync.Map
	mapProducerConsumers *utils.Hashmap
	mapProducers         sync.Map
}

func NewRouter(id string) *Router {
	return &Router{
		id:                   id,
		logger:               utils.NewLogger("router"),
		mapProducerConsumers: utils.NewHashMap(),
	}
}

func (r *Router) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		r.logger.Debug("method=%s,internal=%+v,response:%s", request.Method, request.Internal, response)
	}()

	switch request.Method {
	case mediasoupdata.MethodRouterCreateWebRtcTransport:
		var options mediasoupdata.WebRtcTransportOptions
		_ = json.Unmarshal(request.Data, &options)
		webrtcTransport, err := newWebrtcTransport(webrtcTransportParam{
			options: options,
			transportParam: transportParam{
				Id:                                   request.Internal.TransportId,
				OnTransportNewProducer:               r.OnTransportNewProducer,
				OnTransportProducerClosed:            r.OnTransportProducerClosed,
				OnTransportProducerRtpPacketReceived: r.OnTransportProducerRtpPacketReceived,
				OnTransportNewConsumer:               r.OnTransportNewConsumer,
				OnTransportConsumerClosed:            r.OnTransportConsumerClosed,
			},
		})
		if err != nil {
			response.Err = common.ErrCreateWebrtcTransport
			return
		}
		r.mapTransports.Store(request.Internal.TransportId, webrtcTransport)
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
		t, ok := r.mapTransports.Load(request.Internal.TransportId)
		if !ok {
			response.Err = common.ErrTransportNotFound
			return
		}
		transport := t.(ITransport)
		transport.HandleRequest(request, response)
	}
}

func (r *Router) Close() {
	r.mapTransports.Range(func(key, value interface{}) bool {
		transport := value.(ITransport)
		transport.Close()
		return true
	})
	r.logger.Warn("router:%s stop", r.id)
}

func (r *Router) FillJson() json.RawMessage {
	var transportIds []string
	r.mapTransports.Range(func(key, value interface{}) bool {
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

func (r *Router) OnTransportNewProducer(producer *Producer) error {
	if _, ok := r.mapProducers.Load(producer.id); ok {
		return common.ErrProducerExist
	}
	r.mapProducers.Store(producer.id, producer)

	return nil
}

func (r *Router) OnTransportProducerClosed(producerId string) {
	// close consumers
	value, ok := r.mapProducerConsumers.Get(producerId)
	if !ok {
		return
	}
	consumersMap, ok := value.(map[interface{}]interface{})
	if !ok {
		r.logger.Error("mapProducerConsumers get consumers failed")
		return
	}
	for _, v := range consumersMap {
		v.(IConsumer).Close()
	}
	// clear producer in map
	r.mapProducers.Delete(producerId)
	r.mapProducerConsumers.Erase(producerId)
}

func (r *Router) OnTransportNewConsumer(consumer IConsumer, producerId string) error {
	if _, ok := r.mapProducers.Load(producerId); !ok {
		return common.ErrProducerNotFound
	}
	r.mapProducerConsumers.Store(producerId, consumer.GetId(), consumer)
	r.logger.Debug("OnTransportNewConsumer store mapProducerConsumers, producerId:%s, consumerId:%s", producerId, consumer.GetId())

	return nil
}

func (r *Router) OnTransportConsumerClosed(producerId, consumerId string) {
	v, ok := r.mapProducerConsumers.Load(producerId, consumerId)
	if !ok {
		r.logger.Error("consumer not found[producerId:%s][consumerId:%s]", producerId, consumerId)
		return
	}
	v.(IConsumer).Close()
}

func (r *Router) OnTransportProducerRtpPacketReceived(producer *Producer, packet *rtp.Packet) {
	value, ok := r.mapProducerConsumers.Get(producer.id)
	if !ok {
		r.logger.Trace("no consumers to router RTP")
		return
	}
	consumersMap, ok := value.(map[interface{}]interface{})
	if !ok {
		r.logger.Error("mapProducerConsumers get consumers failed")
		return
	}
	for _, v := range consumersMap {
		v.(IConsumer).SendRtpPacket(packet)
	}
}
