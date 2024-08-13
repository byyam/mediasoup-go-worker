package rtc

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog"

	FBS__DirectTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/DirectTransport"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Response "github.com/byyam/mediasoup-go-worker/fbs/FBS/Response"
	FBS__Router "github.com/byyam/mediasoup-go-worker/fbs/FBS/Router"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
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
	if id == "" {
		return nil
	}
	r := &Router{
		id:                   id,
		logger:               zerowrapper.NewScope("router", id),
		mapProducerConsumers: hashmap.NewHashMap(),
	}
	workerchannel.RegisterHandler(id, r.HandleRequest)
	r.logger.Info().Str("id", id).Msg("new router start")
	return r
}

func (r *Router) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	r.logger.Debug().Str("request", request.String()).Err(response.Err).Msg("handle channel request")

	switch request.MethodType {
	case FBS__Request.MethodROUTER_CREATE_WEBRTCTRANSPORT:
		requestT := request.Request.Body.Value.(*FBS__Router.CreateWebRtcTransportRequestT)
		webrtcTransport, err := newWebrtcTransport(webrtcTransportParam{
			optionsFBS: requestT.Options,
			transportParam: transportParam{
				OptionsFBS:                             requestT.Options.Base,
				Id:                                     requestT.TransportId,
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
		r.mapTransports.Store(requestT.TransportId, webrtcTransport)
		response.Data = webrtcTransport.FillJson()
		// set rsp
		dataDump := &FBS__WebRtcTransport.DumpResponseT{}
		_ = mediasoupdata.Clone(&response.Data, dataDump)
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyWebRtcTransport_DumpResponse,
			Value: dataDump,
		}
		response.RspBody = rspBody

	case FBS__Request.MethodROUTER_CREATE_PLAINTRANSPORT:

	case FBS__Request.MethodROUTER_CREATE_PIPETRANSPORT:
		requestT := request.Request.Body.Value.(*FBS__Router.CreatePipeTransportRequestT)
		var options mediasoupdata.PipeTransportOptions // todo
		_ = json.Unmarshal(request.Data, &options)
		pipeTransport, err := newPipeTransport(pipeTransportParam{
			options: options,
			transportParam: transportParam{
				Options: mediasoupdata.TransportOptions{
					SctpOptions: options.SctpOptions,
				},
				Id:                                     requestT.TransportId,
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
		r.mapTransports.Store(requestT.TransportId, pipeTransport)
		response.Data = pipeTransport.FillJson()

	case FBS__Request.MethodROUTER_CREATE_DIRECTTRANSPORT:
		requestT := request.Request.Body.Value.(*FBS__Router.CreateDirectTransportRequestT)
		directTransport, err := newDirectTransport(directTransportParam{
			optionsFBS: requestT.Options,
			transportParam: transportParam{
				OptionsFBS:                             requestT.Options.Base,
				Id:                                     requestT.TransportId,
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
			r.logger.Error().Err(err).Msgf("createDirectTransport options:%+v", requestT.Options)
			response.Err = mserror.ErrCreateDirectTransport
			return
		}
		r.mapTransports.Store(requestT.TransportId, directTransport)
		response.Data = directTransport.FillJson()
		// set rsp
		dataDump := &FBS__Transport.DumpT{}
		_ = mediasoupdata.Clone(&response.Data, dataDump)
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyDirectTransport_DumpResponse,
			Value: &FBS__DirectTransport.DumpResponseT{Base: dataDump},
		}
		response.RspBody = rspBody

	case FBS__Request.MethodROUTER_CREATE_ACTIVESPEAKEROBSERVER:
		requestT := request.Request.Body.Value.(*FBS__Router.CreateActiveSpeakerObserverRequestT)
		var options mediasoupdata.ActiveSpeakerObserverOptions
		_ = json.Unmarshal(request.Data, &options)
		audioLevelObserver, err := newActiveSpeakerObserver(ActiveSpeakerObserverParam{
			Id:      requestT.RtpObserverId,
			Options: options,
		})
		if err != nil {
			r.logger.Error().Err(err).Msg("newActiveSpeakerObserver")
			response.Err = mserror.ErrCreateActiveSpeakerObserver
			return
		}
		r.mapRtpObservers.Store(requestT.RtpObserverId, audioLevelObserver)

	case FBS__Request.MethodROUTER_CREATE_AUDIOLEVELOBSERVER:
		requestT := request.Request.Body.Value.(*FBS__Router.CreateAudioLevelObserverRequestT)
		var options mediasoupdata.AudioLevelObserverOptions
		_ = json.Unmarshal(request.Data, &options)
		audioLevelObserver, err := newAudioLevelObserver(AudioLevelObserverParam{
			Id:      requestT.RtpObserverId,
			Options: options,
		})
		if err != nil {
			r.logger.Error().Err(err).Msg("newAudioLevelObserver")
			response.Err = mserror.ErrCreateAudioLevelObserver
			return
		}
		r.mapRtpObservers.Store(requestT.RtpObserverId, audioLevelObserver)

	case FBS__Request.MethodROUTER_DUMP:
		response.Data = r.FillJson()

	case FBS__Request.MethodROUTER_CLOSE_TRANSPORT:
		requestT := request.Request.Body.Value.(*FBS__Router.CloseTransportRequestT)
		v, ok := r.mapTransports.Load(requestT.TransportId)
		if !ok {
			response.Err = mserror.ErrTransportNotFound
			return
		}
		transport := v.(ITransport)
		transport.NotifyClose() // call son close, tiger this close

	default:
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
