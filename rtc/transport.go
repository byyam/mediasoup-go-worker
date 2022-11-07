package rtc

import (
	"encoding/json"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pion/rtcp"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpprobation"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type ITransport interface {
	Connected()
	Close()
	FillJson() json.RawMessage
	HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData)
	ReceiveRtpPacket(packet *rtpparser.Packet)
	ReceiveRtcpPacket(header *rtcp.Header, packets []rtcp.Packet)
}

type Transport struct {
	id             string
	direct         bool
	maxMessageSize uint32
	logger         zerolog.Logger

	mapProducers       sync.Map //map[string]*Producer
	mapConsumers       sync.Map
	mapSsrcConsumer    sync.Map
	mapRtxSsrcConsumer sync.Map
	rtpListener        *RtpListener

	// handler
	onTransportNewProducerHandler                 atomic.Value
	onTransportProducerClosedHandler              func(producerId string)
	onTransportProducerRtpPacketReceivedHandler   func(*Producer, *rtpparser.Packet)
	onTransportNewConsumerHandler                 func(consumer IConsumer, producerId string) error
	onTransportConsumerClosedHandler              func(producerId, consumerId string)
	onTransportConsumerKeyFrameRequestedHandler   func(consumerId string, mappedSsrc uint32)
	onTransportNeedWorstRemoteFractionLostHandler func(producerId string, worstRemoteFractionLost *uint8)

	// transport base call sons
	sendRtpPacketFunc          func(packet *rtpparser.Packet)
	sendRtcpPacketFunc         func(packet rtcp.Packet)
	sendRtcpCompoundPacketFunc func(packets []rtcp.Packet)
	notifyCloseFunc            func()

	// close
	closeOnce sync.Once
}

func (t *Transport) Close() {
	t.closeOnce.Do(func() {
		t.logger.Info().Msg("closed")
		workerchannel.UnregisterHandler(t.id)
	})
}

func (t *Transport) FillJson() json.RawMessage {
	var producerIds []string
	t.mapProducers.Range(func(key, value interface{}) bool {
		producerIds = append(producerIds, key.(string))
		return true
	})
	dumpData := mediasoupdata.TransportDump{
		Id:                      t.id,
		Direct:                  false,
		ProducerIds:             producerIds,
		ConsumerIds:             nil,
		MapSsrcConsumerId:       nil,
		MapRtxSsrcConsumerId:    nil,
		DataProducerIds:         nil,
		DataConsumerIds:         nil,
		RecvRtpHeaderExtensions: nil,
		RtpListener:             nil,
		SctpParameters:          mediasoupdata.SctpParameters{},
		SctpState:               "",
		SctpListener:            nil,
		TraceEventTypes:         "",
		PlainTransportDump:      nil,
		WebRtcTransportDump:     nil,
	}
	data, _ := json.Marshal(&dumpData)
	t.logger.Debug().Msgf("dumpData:%+v", dumpData)
	return data
}

func (t *Transport) FillJsonStats() json.RawMessage {
	jsonData := mediasoupdata.TransportStat{
		Type:                        "",
		TransportId:                 "",
		Timestamp:                   0,
		SctpState:                   "",
		BytesReceived:               0,
		RecvBitrate:                 0,
		BytesSent:                   0,
		SendBitrate:                 0,
		RtpBytesReceived:            0,
		RtpRecvBitrate:              0,
		RtpBytesSent:                0,
		RtpSendBitrate:              0,
		RtxBytesReceived:            0,
		RtxRecvBitrate:              0,
		RtxBytesSent:                0,
		RtxSendBitrate:              0,
		ProbationBytesSent:          0,
		ProbationSendBitrate:        0,
		AvailableOutgoingBitrate:    0,
		AvailableIncomingBitrate:    0,
		MaxIncomingBitrate:          0,
		RtpPacketLossReceived:       0,
		RtpPacketLossSent:           0,
		WebRtcTransportSpecificStat: nil,
	}
	data, _ := json.Marshal(&([]mediasoupdata.TransportStat{jsonData}))
	t.logger.Debug().Msgf("getStats:%+v", jsonData)
	return data
}

type transportParam struct {
	Id                                     string
	Direct                                 bool
	MaxMessageSize                         uint32
	OnTransportNewProducer                 func(producer *Producer) error
	OnTransportProducerClosed              func(producerId string)
	OnTransportProducerRtpPacketReceived   func(producer *Producer, packet *rtpparser.Packet)
	OnTransportNewConsumer                 func(consumer IConsumer, producerId string) error
	OnTransportConsumerClosed              func(consumerId, producerId string)
	OnTransportConsumerKeyFrameRequested   func(consumerId string, mappedSsrc uint32)
	OnTransportNeedWorstRemoteFractionLost func(producerId string, worstRemoteFractionLost *uint8)
	// call webrtcTransport
	SendRtpPacketFunc          func(packet *rtpparser.Packet)
	SendRtcpPacketFunc         func(packet rtcp.Packet)
	SendRtcpCompoundPacketFunc func(packets []rtcp.Packet)
	NotifyCloseFunc            func()
}

func (t transportParam) valid() bool {
	if t.Id == "" {
		return false
	}
	if t.OnTransportNewProducer == nil || t.OnTransportProducerRtpPacketReceived == nil {
		return false
	}
	if t.SendRtpPacketFunc == nil || t.SendRtcpPacketFunc == nil || t.SendRtcpCompoundPacketFunc == nil || t.NotifyCloseFunc == nil {
		return false
	}
	return true
}

func newTransport(param transportParam) (ITransport, error) {
	if !param.valid() {
		return nil, mserror.ErrInvalidParam
	}
	transport := &Transport{
		id:             param.Id,
		direct:         param.Direct,
		maxMessageSize: param.MaxMessageSize,
		logger:         zerowrapper.NewScope("transport", param.Id),
		rtpListener:    newRtpListener(),
	}
	transport.onTransportNewProducerHandler.Store(param.OnTransportNewProducer)
	transport.onTransportProducerClosedHandler = param.OnTransportProducerClosed
	transport.onTransportProducerRtpPacketReceivedHandler = param.OnTransportProducerRtpPacketReceived
	transport.onTransportNewConsumerHandler = param.OnTransportNewConsumer
	transport.onTransportConsumerClosedHandler = param.OnTransportConsumerClosed
	transport.onTransportConsumerKeyFrameRequestedHandler = param.OnTransportConsumerKeyFrameRequested
	transport.onTransportNeedWorstRemoteFractionLostHandler = param.OnTransportNeedWorstRemoteFractionLost
	transport.sendRtpPacketFunc = param.SendRtpPacketFunc
	transport.sendRtcpPacketFunc = param.SendRtcpPacketFunc
	transport.sendRtcpCompoundPacketFunc = param.SendRtcpCompoundPacketFunc
	transport.notifyCloseFunc = param.NotifyCloseFunc
	go transport.OnTimer()

	return transport, nil
}

func (t *Transport) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		t.logger.Info().Str("request", request.String()).Str("response", response.String()).Msg("handle channel request done")
	}()

	switch request.Method {

	case mediasoupdata.MethodTransportDump:
		response.Data = t.FillJson()

	case mediasoupdata.MethodTransportClose:
		t.notifyCloseFunc() // call son close, tiger this close

	case mediasoupdata.MethodTransportProduce:
		var options mediasoupdata.ProducerOptions
		_ = json.Unmarshal(request.Data, &options)
		data, err := t.Produce(request.Internal.ProducerId, options)
		response.Data, _ = json.Marshal(data)
		response.Err = err

	case mediasoupdata.MethodTransportConsume:
		var options mediasoupdata.ConsumerOptions
		_ = json.Unmarshal(request.Data, &options)
		data, err := t.Consume(request.Internal.ProducerId, request.Internal.ConsumerId, options)
		response.Data, _ = json.Marshal(data)
		response.Err = err

	case mediasoupdata.MethodTransportProduceData:
		var options mediasoupdata.DataProducerOptions
		_ = json.Unmarshal(request.Data, &options)
		dataProducer, err := t.DataProduce(request.Internal.DataProducerId, options)
		if err != nil {
			response.Err = err
			return
		}
		response.Data = dataProducer.FillJson()

	case mediasoupdata.MethodTransportConsumeData:

	case mediasoupdata.MethodTransportSetMaxIncomingBitrate:

	case mediasoupdata.MethodTransportSetMaxOutgoingBitrate:

	case mediasoupdata.MethodTransportEnableTraceEvent:

	case mediasoupdata.MethodTransportGetStats:
		response.Data = t.FillJsonStats()

	// producer
	case mediasoupdata.MethodProducerDump, mediasoupdata.MethodProducerGetStats, mediasoupdata.MethodProducerPause,
		mediasoupdata.MethodProducerResume, mediasoupdata.MethodProducerEnableTraceEvent:
		value, ok := t.mapProducers.Load(request.Internal.ProducerId)
		if !ok {
			response.Err = mserror.ErrProducerNotFound
			return
		}
		producer := value.(*Producer)
		producer.HandleRequest(request, response)

	case mediasoupdata.MethodProducerClose:
		value, ok := t.mapProducers.Load(request.Internal.ProducerId)
		if !ok {
			response.Err = mserror.ErrProducerNotFound
			return
		}
		producer := value.(*Producer)
		producer.Close()
		t.mapProducers.Delete(request.Internal.ProducerId)
		t.onTransportProducerClosedHandler(producer.id)

	// consumer
	case mediasoupdata.MethodConsumerDump, mediasoupdata.MethodConsumerGetStats, mediasoupdata.MethodConsumerPause,
		mediasoupdata.MethodConsumerResume, mediasoupdata.MethodConsumerSetPreferredLayers, mediasoupdata.MethodConsumerSetPriority,
		mediasoupdata.MethodConsumerRequestKeyFrame, mediasoupdata.MethodConsumerEnableTraceEvent:
		value, ok := t.mapConsumers.Load(request.Internal.ConsumerId)
		if !ok {
			response.Err = mserror.ErrConsumerNotFound
			return
		}
		consumer := value.(IConsumer)
		consumer.HandleRequest(request, response)

	case mediasoupdata.MethodConsumerClose:
		value, ok := t.mapConsumers.Load(request.Internal.ConsumerId)
		if !ok {
			response.Err = mserror.ErrConsumerNotFound
			return
		}
		consumer := value.(IConsumer)
		consumer.Close()
		t.mapConsumers.Delete(request.Internal.ConsumerId)
		for _, ssrc := range consumer.GetMediaSsrcs() {
			t.mapSsrcConsumer.Delete(ssrc)
		}
		t.onTransportConsumerClosedHandler(request.Internal.ProducerId, consumer.GetId())

	default:
		t.logger.Error().Str("method", request.Method).Msg("transport handle request method not found")
		return
	}
}

func (t *Transport) Consume(producerId, consumerId string, options mediasoupdata.ConsumerOptions) (*mediasoupdata.ConsumerData, error) {
	if producerId == "" || consumerId == "" {
		return nil, mserror.ErrInvalidParam
	}

	var consumer IConsumer
	var err error
	switch options.Type {
	case mediasoupdata.ConsumerType_Simple:
		consumer, err = newSimpleConsumer(simpleConsumerParam{
			consumerParam: consumerParam{
				id:                     consumerId,
				producerId:             producerId,
				kind:                   options.Kind,
				rtpParameters:          options.RtpParameters,
				consumableRtpEncodings: options.ConsumableRtpEncodings,
			},
			OnConsumerSendRtpPacket:       t.OnConsumerSendRtpPacket,
			OnConsumerKeyFrameRequested:   t.OnConsumerKeyFrameRequested,
			OnConsumerRetransmitRtpPacket: t.OnConsumerRetransmitRtpPacket,
		})

	case mediasoupdata.ConsumerType_Simulcast: // todo...
	case mediasoupdata.ConsumerType_Svc:
	default:
		t.logger.Error().Str("type", string(options.Type)).Msg("unsupported consumer type")
		return nil, mserror.ErrInvalidParam
	}

	if err != nil {
		t.logger.Error().Msgf("create consumer[%s] failed:%v", options.Type, err)
		return nil, err
	}
	if err := t.onTransportNewConsumerHandler(consumer, producerId); err != nil {
		return nil, err
	}

	t.mapConsumers.Store(consumerId, consumer)
	for _, ssrc := range consumer.GetMediaSsrcs() {
		t.mapSsrcConsumer.Store(ssrc, consumer)
	}
	t.logger.Debug().Msgf("Consumer created [producerId:%s][consumerId:%s],type:%s,kind:%s,ssrc:%v", producerId, consumerId, options.Type, options.Kind, consumer.GetMediaSsrcs())
	return &mediasoupdata.ConsumerData{
		Paused:         false,
		ProducerPaused: false,
		Score:          mediasoupdata.ConsumerScore{},
	}, nil
}

func (t *Transport) Produce(id string, options mediasoupdata.ProducerOptions) (*mediasoupdata.ProducerData, error) {
	if id == "" {
		return nil, mserror.ErrInvalidParam
	}
	if _, ok := t.mapProducers.Load(id); ok {
		return nil, mserror.ErrDuplicatedId
	}
	producer, err := newProducer(producerParam{
		id:                                    id,
		options:                               options,
		OnProducerRtpPacketReceived:           t.OnProducerRtpPacketReceived,
		OnProducerSendRtcpPacket:              t.OnProducerSendRtcpPacket,
		OnProducerNeedWorstRemoteFractionLost: t.onTransportNeedWorstRemoteFractionLostHandler,
	})
	if err != nil {
		t.logger.Err(err).Msg("produce failed")
		return nil, err
	}
	if err = t.rtpListener.AddProducer(producer); err != nil {
		return nil, err
	}
	if handler, ok := t.onTransportNewProducerHandler.Load().(func(*Producer) error); ok && handler != nil {
		if err := handler(producer); err != nil {
			return nil, err
		}
	}
	t.mapProducers.Store(id, producer)
	t.logger.Debug().Msgf("Producer created [producerId:%s],type:%s", id, producer.Type)
	// todo

	return &mediasoupdata.ProducerData{Type: producer.Type}, nil
}

func (t *Transport) DataProduce(id string, options mediasoupdata.DataProducerOptions) (*DataProducer, error) {
	if id == "" {
		return nil, mserror.ErrInvalidParam
	}
	dataProducer, err := newDataProducer(id, t.maxMessageSize, options)
	if err != nil {
		t.logger.Err(err).Msg("data produce failed")
		return nil, err
	}
	// todo: store in map
	t.logger.Debug().Msgf("DataProducer created [producerId:%s],type:%s", id, dataProducer.options.Type)
	return dataProducer, nil
}

func (t *Transport) ReceiveRtpPacket(packet *rtpparser.Packet) {
	// get producer from ssrc, to producer
	producer := t.rtpListener.GetProducer(packet)
	if producer == nil {
		monitor.RtpRecvCount(monitor.TraceSsrcNotFound)
		return
	}

	producer.ReceiveRtpPacket(packet)
}

func (t *Transport) OnProducerRtpPacketReceived(producer *Producer, packet *rtpparser.Packet) {
	t.onTransportProducerRtpPacketReceivedHandler(producer, packet)
}

func (t *Transport) OnProducerSendRtcpPacket(packet rtcp.Packet) {
	t.sendRtcpPacketFunc(packet)
}

func (t *Transport) OnConsumerSendRtpPacket(consumer IConsumer, packet *rtpparser.Packet) {
	t.logger.Trace().Msgf("OnConsumerSendRtpPacket:%+v", packet.Header)
	t.sendRtpPacketFunc(packet)
}

func (t *Transport) OnConsumerRetransmitRtpPacket(packet *rtpparser.Packet) {
	// todo: tcc
	t.sendRtpPacketFunc(packet)
}

func (t *Transport) ReceiveRtcpPacket(header *rtcp.Header, packets []rtcp.Packet) {
	t.logger.Info().Msgf("ReceiveRtcpPacket[%d]:\n%+v\nheader:%+v", len(packets), packets, header)
	for _, packet := range packets {
		t.HandleRtcpPacket(header, packet)
	}
}

func (t *Transport) HandleRtcpPacket(header *rtcp.Header, packet rtcp.Packet) {
	t.logger.Debug().Msgf("HandleRtcpPacket:%v", packet.DestinationSSRC())
	switch packet.(type) {
	case *rtcp.SenderReport:
		pkg := packet.(*rtcp.SenderReport)
		for _, sr := range pkg.Reports {
			t.logger.Debug().Msgf("handle SR:%s,report:%+v", pkg.String(), sr)
			producer := t.rtpListener.GetProducerBySSRC(sr.SSRC)
			if producer == nil {
				t.logger.Warn().Msgf("no Producer found for received Sender Report [ssrc:%d]", sr.SSRC)
				continue
			}
			producer.ReceiveRtcpSenderReport(&sr)
		}
	case *rtcp.ReceiverReport:
		pkg := packet.(*rtcp.ReceiverReport)
		for _, rr := range pkg.Reports {
			t.logger.Debug().Msgf("handle RR:%s,report:%+v", pkg.String(), rr)
			// Special case for the RTP probator.
			if rr.SSRC == rtpprobation.RtpProbationSsrc {
				continue
			}
			consumer, ok := t.mapSsrcConsumer.Load(rr.SSRC)
			if !ok {
				// Special case for the RTP probator.
				if rr.SSRC == rtpprobation.RtpProbationSsrc {
					continue
				}
				// Special case for (unused) RTCP-RR from the RTX stream.
				_, ok := t.mapRtxSsrcConsumer.Load(rr.SSRC)
				if ok {
					continue
				}
				t.logger.Warn().Msgf("no Consumer found for received Receiver Report [ssrc %d]", rr.SSRC)
				continue
			}
			consumer.(IConsumer).ReceiveRtcpReceiverReport(&rr)
			// todo
		}
	case *rtcp.SourceDescription:
		pkg := packet.(*rtcp.SourceDescription)
		t.logger.Debug().Msgf("%s", pkg.String())
	case *rtcp.Goodbye:
		pkg := packet.(*rtcp.Goodbye)
		t.logger.Debug().Msgf("ignoring received RTCP BYE %s", pkg.String())
	case *rtcp.FullIntraRequest:
		pkg := packet.(*rtcp.FullIntraRequest)
		t.ReceiveKeyFrameRequest(header.Count, pkg.MediaSSRC)
		monitor.KeyframeCount(pkg.MediaSSRC, monitor.KeyframeRecvFIR)
	case *rtcp.PictureLossIndication:
		pkg := packet.(*rtcp.PictureLossIndication)
		t.ReceiveKeyFrameRequest(header.Count, pkg.MediaSSRC)
		monitor.KeyframeCount(pkg.MediaSSRC, monitor.KeyframeRecvPLI)
	case *rtcp.ReceiverEstimatedMaximumBitrate:
		pkg := packet.(*rtcp.ReceiverEstimatedMaximumBitrate)
		t.logger.Debug().Msgf("%s", pkg.String())
	case *rtcp.TransportLayerNack:
		pkg := packet.(*rtcp.TransportLayerNack)
		t.logger.Debug().Msgf("TransportLayerNack:%+v", pkg)
		consumer, ok := t.mapSsrcConsumer.Load(pkg.MediaSSRC)
		if !ok {
			t.logger.Warn().Msgf("no Consumer found for received NACK Feedback packet [sender ssrc:%d, media ssrc:%d]", pkg.SenderSSRC, pkg.MediaSSRC)
			return
		}
		consumer.(IConsumer).ReceiveNack(pkg)

	case *rtcp.TransportLayerCC:
		pkg := packet.(*rtcp.TransportLayerCC)
		t.logger.Info().Msgf("TransportLayerCC:%+v", pkg)
	default:
		monitor.RtcpRecvCount(monitor.TraceUnknownRtcpType)
		t.logger.Warn().Msgf("unhandled RTCP type received %+v", header)
	}
}

func (t *Transport) ReceiveKeyFrameRequest(feedbackFormat uint8, ssrc uint32) {
	v, ok := t.mapSsrcConsumer.Load(ssrc)
	if !ok {
		return
	}
	consumer := v.(IConsumer)
	consumer.ReceiveKeyFrameRequest(feedbackFormat, ssrc)
}

func (t *Transport) OnConsumerKeyFrameRequested(consumer IConsumer, mappedSsrc uint32) {
	t.onTransportConsumerKeyFrameRequestedHandler(consumer.GetId(), mappedSsrc)
}

func (t *Transport) OnTimer() {
	rtcpTimer := rand.Int31n(1000) + 500 // 0.5-1.5 s
	for {
		// todo: update interval
		time.Sleep(time.Millisecond * time.Duration(rtcpTimer))
		now := time.Now()
		t.SendRtcp(now)
		t.sendNacks()
	}
}

func (t *Transport) SendRtcp(now time.Time) {
	t.mapConsumers.Range(func(id, value interface{}) bool {
		consumer, ok := value.(IConsumer)
		if !ok || consumer == nil {
			return true
		}
		rtpStreams := consumer.GetRtpStreams()
		for _, rtpStream := range rtpStreams {
			packets := consumer.GetRtcp(rtpStream, now)
			if len(packets) == 0 {
				continue
			}
			t.sendRtcpCompoundPacketFunc(packets)
		}
		return true
	})
	t.mapProducers.Range(func(id, value interface{}) bool {
		producer, ok := value.(*Producer)
		if !ok || producer == nil {
			return true
		}
		// One more RR would exceed the MTU, send the compound packet now.
		packets := producer.GetRtcp(now)
		if len(packets) != 0 {
			t.sendRtcpCompoundPacketFunc(packets)
		}
		return true
	})
}

func (t *Transport) sendNacks() {
	t.mapProducers.Range(func(id, value interface{}) bool {
		producer, ok := value.(*Producer)
		if !ok || producer == nil {
			return true
		}
		producer.mapSsrcRtpStream.Range(func(key, value interface{}) bool {
			ssrc := key.(uint32)
			rtpStream, ok := value.(*RtpStreamRecv)
			if !ok || rtpStream == nil {
				return true
			}
			pairs, _ := rtpStream.nackGenerator.Pairs()
			if len(pairs) > 0 {
				packets := []rtcp.Packet{
					&rtcp.TransportLayerNack{
						MediaSSRC: ssrc,
						Nacks:     pairs,
					},
				}
				t.sendRtcpCompoundPacketFunc(packets)
			}
			return true
		})
		return true
	})
}

func (t *Transport) Connected() {

}
