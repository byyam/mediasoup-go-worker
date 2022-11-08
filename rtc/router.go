package rtc

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/byyam/mediasoup-go-worker/mserror"

	"github.com/byyam/mediasoup-go-worker/internal/hashmap"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type Router struct {
	id                   string
	logger               zerolog.Logger
	mapTransports        sync.Map
	mapProducerConsumers *hashmap.Hashmap
	mapProducers         sync.Map
	mapConsumerProducer  sync.Map
	mapRtpObservers      sync.Map
}

func NewRouter(id string) *Router {
	r := &Router{
		id:                   id,
		logger:               zerowrapper.NewScope("router", id),
		mapProducerConsumers: hashmap.NewHashMap(),
	}
	workerchannel.RegisterHandler(id, r.HandleRequest)
	return r
}

func (r *Router) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		r.logger.Info().Str("request", request.String()).Str("response", response.String()).Msg("handle channel request done")
	}()

	switch request.Method {
	case mediasoupdata.MethodRouterCreateWebRtcTransport:
		var options mediasoupdata.WebRtcTransportOptions
		_ = json.Unmarshal(request.Data, &options)
		webrtcTransport, err := newWebrtcTransport(webrtcTransportParam{
			options: options,
			transportParam: transportParam{
				Options: mediasoupdata.TransportOptions{
					SctpOptions:                     options.SctpOptions,
					InitialAvailableOutgoingBitrate: options.InitialAvailableOutgoingBitrate,
				},
				Id:                                     request.Internal.TransportId,
				OnTransportNewProducer:                 r.OnTransportNewProducer,
				OnTransportProducerClosed:              r.OnTransportProducerClosed,
				OnTransportProducerRtpPacketReceived:   r.OnTransportProducerRtpPacketReceived,
				OnTransportNewConsumer:                 r.OnTransportNewConsumer,
				OnTransportConsumerClosed:              r.OnTransportConsumerClosed,
				OnTransportConsumerKeyFrameRequested:   r.OnTransportConsumerKeyFrameRequested,
				OnTransportNeedWorstRemoteFractionLost: r.OnTransportNeedWorstRemoteFractionLost,
			},
		})
		if err != nil {
			r.logger.Error().Err(err).Msg("createWebrtcTransport")
			response.Err = mserror.ErrCreateWebrtcTransport
			return
		}
		r.mapTransports.Store(request.Internal.TransportId, webrtcTransport)
		response.Data = webrtcTransport.FillJson()

	case mediasoupdata.MethodRouterCreatePlainTransport:

	case mediasoupdata.MethodRouterCreatePipeTransport:
		var options mediasoupdata.PipeTransportOptions
		_ = json.Unmarshal(request.Data, &options)
		pipeTransport, err := newPipeTransport(pipeTransportParam{
			options: options,
			transportParam: transportParam{
				Options: mediasoupdata.TransportOptions{
					SctpOptions: options.SctpOptions,
				},
				Id:                                     request.Internal.TransportId,
				OnTransportNewProducer:                 r.OnTransportNewProducer,
				OnTransportProducerClosed:              r.OnTransportProducerClosed,
				OnTransportProducerRtpPacketReceived:   r.OnTransportProducerRtpPacketReceived,
				OnTransportNewConsumer:                 r.OnTransportNewConsumer,
				OnTransportConsumerClosed:              r.OnTransportConsumerClosed,
				OnTransportConsumerKeyFrameRequested:   r.OnTransportConsumerKeyFrameRequested,
				OnTransportNeedWorstRemoteFractionLost: r.OnTransportNeedWorstRemoteFractionLost,
			},
		})
		if err != nil {
			r.logger.Error().Err(err).Msg("createPipeTransport")
			response.Err = mserror.ErrCreatePipeTransport
			return
		}
		r.mapTransports.Store(request.Internal.TransportId, pipeTransport)
		response.Data = pipeTransport.FillJson()

	case mediasoupdata.MethodRouterCreateDirectTransport:
		var options mediasoupdata.DirectTransportOptions
		_ = json.Unmarshal(request.Data, &options)
		directTransport, err := newDirectTransport(directTransportParam{
			options: options,
			transportParam: transportParam{
				Options: mediasoupdata.TransportOptions{
					DirectTransportOptions: options,
				},
				Id:                                     request.Internal.TransportId,
				OnTransportNewProducer:                 r.OnTransportNewProducer,
				OnTransportProducerClosed:              r.OnTransportProducerClosed,
				OnTransportProducerRtpPacketReceived:   r.OnTransportProducerRtpPacketReceived,
				OnTransportNewConsumer:                 r.OnTransportNewConsumer,
				OnTransportConsumerClosed:              r.OnTransportConsumerClosed,
				OnTransportConsumerKeyFrameRequested:   r.OnTransportConsumerKeyFrameRequested,
				OnTransportNeedWorstRemoteFractionLost: r.OnTransportNeedWorstRemoteFractionLost,
			},
		})
		if err != nil {
			r.logger.Error().Err(err).Msg("createDirectTransport")
			response.Err = mserror.ErrCreateDirectTransport
			return
		}
		r.mapTransports.Store(request.Internal.TransportId, directTransport)
		response.Data = directTransport.FillJson()

	case mediasoupdata.MethodRouterCreateActiveSpeakerObserver:
		var options mediasoupdata.ActiveSpeakerObserverOptions
		_ = json.Unmarshal(request.Data, &options)
		audioLevelObserver, err := newActiveSpeakerObserver(ActiveSpeakerObserverParam{
			Id:      request.Internal.RtpObserverId,
			Options: options,
		})
		if err != nil {
			r.logger.Error().Err(err).Msg("newActiveSpeakerObserver")
			response.Err = mserror.ErrCreateActiveSpeakerObserver
			return
		}
		r.mapRtpObservers.Store(request.Internal.RtpObserverId, audioLevelObserver)

	case mediasoupdata.MethodRouterCreateAudioLevelObserver:
		var options mediasoupdata.AudioLevelObserverOptions
		_ = json.Unmarshal(request.Data, &options)
		audioLevelObserver, err := newAudioLevelObserver(AudioLevelObserverParam{
			Id:      request.Internal.RtpObserverId,
			Options: options,
		})
		if err != nil {
			r.logger.Error().Err(err).Msg("newAudioLevelObserver")
			response.Err = mserror.ErrCreateAudioLevelObserver
			return
		}
		r.mapRtpObservers.Store(request.Internal.RtpObserverId, audioLevelObserver)

	case mediasoupdata.MethodRouterDump:
		response.Data = r.FillJson()

	case mediasoupdata.MethodRouterClose:
		r.Close()
	default:
		//t, ok := r.mapTransports.Load(request.Internal.TransportId)
		//if !ok {
		//	response.Err = mserror.ErrTransportNotFound
		//	return
		//}
		//transport := t.(ITransport)
		//transport.HandleRequest(request, response)
		r.logger.Error().Str("method", request.Method).Msg("router handle request method not found")
		response.Err = mserror.ErrInvalidMethod
		return
	}
}

func (r *Router) Close() {
	r.mapTransports.Range(func(key, value interface{}) bool {
		transport := value.(ITransport)
		transport.Close()
		return true
	})
	r.mapProducers.Range(func(key, value interface{}) bool {
		producer := value.(*Producer)
		producer.Close()
		return true
	})
	workerchannel.UnregisterHandler(r.id)
	r.logger.Warn().Msg("router stop")
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
	r.logger.Debug().Msgf("dumpData:%+v", dumpData)
	return data
}

func (r *Router) OnTransportNewProducer(producer *Producer) error {
	if _, ok := r.mapProducers.Load(producer.id); ok {
		return mserror.ErrProducerExist
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
		r.logger.Error().Msg("mapProducerConsumers get consumers failed")
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
		return mserror.ErrProducerNotFound
	}
	r.mapProducerConsumers.Store(producerId, consumer.GetId(), consumer)
	r.mapConsumerProducer.Store(consumer.GetId(), producerId)
	r.logger.Debug().Str("producerId", producerId).
		Str("consumerId", consumer.GetId()).
		Msg("OnTransportNewConsumer store mapProducerConsumers")

	return nil
}

func (r *Router) OnTransportConsumerClosed(producerId, consumerId string) {
	// clear mapConsumerProducer
	r.mapConsumerProducer.Delete(consumerId)
	// clear mapProducerConsumers
	v, ok := r.mapProducerConsumers.Load(producerId, consumerId)
	if !ok {
		r.logger.Error().Str("producerId", producerId).
			Str("consumerId", consumerId).
			Msg("consumer not found in mapProducerConsumers")

	} else {
		v.(IConsumer).Close()
		r.mapProducerConsumers.Delete(producerId, consumerId)
	}
}

func (r *Router) OnTransportProducerRtpPacketReceived(producer *Producer, packet *rtpparser.Packet) {
	value, ok := r.mapProducerConsumers.Get(producer.id)
	if !ok {
		r.logger.Trace().Msg("no consumers to router RTP")
		return
	}
	consumersMap, ok := value.(map[interface{}]interface{})
	if !ok {
		r.logger.Error().Msg("mapProducerConsumers get consumers failed")
		return
	}
	for _, v := range consumersMap {
		consumer := v.(IConsumer)
		mid := consumer.GetRtpParameters().Mid
		if mid != "" {
			if err := packet.UpdateMid(mid); err != nil {
				r.logger.Warn().Err(err).Msg("UpdateMid in OnTransportProducerRtpPacketReceived failed")
			}
		}
		consumer.SendRtpPacket(packet)
	}
}

func (r *Router) OnTransportConsumerKeyFrameRequested(consumerId string, mappedSsrc uint32) {
	v, ok := r.mapConsumerProducer.Load(consumerId)
	if !ok {
		r.logger.Error().Str("consumerId", consumerId).Msg("OnTransportConsumerKeyFrameRequested producer not found")
		return
	}
	producerId := v.(string)
	v, ok = r.mapProducers.Load(producerId)
	if !ok {
		r.logger.Error().Str("producerId", producerId).Msg("OnTransportConsumerKeyFrameRequested producerId not found")
		return
	}
	v.(*Producer).RequestKeyFrame(mappedSsrc)
}

func (r *Router) OnTransportNeedWorstRemoteFractionLost(producerId string, worstRemoteFractionLost *uint8) {
	value, ok := r.mapProducerConsumers.Get(producerId)
	if !ok {
		r.logger.Trace().Msg("no consumers to router RTP")
		return
	}
	consumersMap, ok := value.(map[interface{}]interface{})
	if !ok {
		r.logger.Error().Msg("mapProducerConsumers get consumers failed")
		return
	}
	for _, v := range consumersMap {
		v.(IConsumer).NeedWorstRemoteFractionLost(worstRemoteFractionLost)
	}
}
