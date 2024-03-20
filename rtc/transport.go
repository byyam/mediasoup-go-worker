package rtc

import (
	"encoding/json"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kr/pretty"
	"github.com/pion/rtcp"
	"github.com/rs/zerolog"

	FBS__DataConsumer "github.com/byyam/mediasoup-go-worker/fbs/FBS/DataConsumer"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Response "github.com/byyam/mediasoup-go-worker/fbs/FBS/Response"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtctime"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpprobation"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type ITransport interface {
	Connected()
	Close()
	GetJson(data *FBS__Transport.DumpT)
	GetBaseStats() *FBS__Transport.StatsT
	FillJson() json.RawMessage
	HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData)
	ReceiveRtpPacket(packet *rtpparser.Packet)
	ReceiveRtcpPacket(header *rtcp.Header, packets []rtcp.Packet)
}

type Transport struct {
	id         string
	optionsFBS *FBS__Transport.OptionsT
	logger     zerolog.Logger

	mapProducers              sync.Map //map[string]*Producer
	mapConsumers              sync.Map
	mapDataProducers          sync.Map
	mapDataConsumers          sync.Map
	mapSsrcConsumer           sync.Map
	mapRtxSsrcConsumer        sync.Map
	rtpListener               *RtpListener
	sctpAssociation           *SctpAssociation
	recvRtpHeaderExtensionIds RtpHeaderExtensionIds

	// handler
	onTransportNewProducerHandler                 atomic.Value
	onTransportProducerClosedHandler              func(producerId string)
	onTransportProducerRtpPacketReceivedHandler   func(*Producer, *rtpparser.Packet)
	onTransportNewConsumerHandler                 func(consumer IConsumer, producerId string) error
	onTransportConsumerClosedHandler              func(producerId, consumerId string)
	onTransportConsumerKeyFrameRequestedHandler   func(consumerId string, mappedSsrc uint32)
	onTransportNeedWorstRemoteFractionLostHandler func(producerId string, worstRemoteFractionLost *uint8)
	onTransportDataProducerClosedHandler          func(producerId string)

	// transport base call sons
	sendRtpPacketFunc          func(packet *rtpparser.Packet)
	sendRtcpPacketFunc         func(packet rtcp.Packet)
	sendRtcpCompoundPacketFunc func(packets []rtcp.Packet)
	NotifyCloseFunc            func()

	// close
	closeOnce sync.Once
}

func (t *Transport) Close() {
	t.closeOnce.Do(func() {
		t.logger.Info().Msg("closed")
		workerchannel.UnregisterHandler(t.id)
	})
}

func (t *Transport) GetJson(data *FBS__Transport.DumpT) {
	var producerIds []string
	t.mapProducers.Range(func(key, value interface{}) bool {
		producerIds = append(producerIds, key.(string))
		return true
	})

	data.Id = t.id
	data.Direct = t.optionsFBS.Direct
	data.ProducerIds = producerIds
	if t.sctpAssociation != nil {
		data.SctpParameters = t.sctpAssociation.GetSctpAssociationParam()
	}
	data.RecvRtpHeaderExtensions = &FBS__Transport.RecvRtpHeaderExtensionsT{}
	data.RtpListener = &FBS__Transport.RtpListenerT{}
}

func (t *Transport) FillJson() json.RawMessage {
	dumpData := &mediasoupdata.TransportDump{}

	data, _ := json.Marshal(dumpData)
	t.logger.Debug().Str("data", string(data)).Msg("dumpData")
	return data
}

func (t *Transport) GetBaseStats() *FBS__Transport.StatsT {
	stats := &FBS__Transport.StatsT{
		TransportId:              t.id,
		Timestamp:                uint64(rtctime.GetTimeMs()),
		SctpState:                nil,
		BytesReceived:            0,
		RecvBitrate:              0,
		BytesSent:                0,
		SendBitrate:              0,
		RtpBytesReceived:         0,
		RtpRecvBitrate:           0,
		RtpBytesSent:             0,
		RtpSendBitrate:           0,
		RtxBytesReceived:         0,
		RtxRecvBitrate:           0,
		RtxBytesSent:             0,
		RtxSendBitrate:           0,
		ProbationBytesSent:       0,
		ProbationSendBitrate:     0,
		AvailableOutgoingBitrate: nil,
		AvailableIncomingBitrate: nil,
		MaxIncomingBitrate:       nil,
		MaxOutgoingBitrate:       nil,
		MinOutgoingBitrate:       nil,
		RtpPacketLossReceived:    nil,
		RtpPacketLossSent:        nil,
	}
	return stats
}

type transportParam struct {
	Id                                     string
	Options                                mediasoupdata.TransportOptions
	OptionsFBS                             *FBS__Transport.OptionsT
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
		id:          param.Id,
		optionsFBS:  param.OptionsFBS,
		logger:      zerowrapper.NewScope("transport", param.Id),
		rtpListener: newRtpListener(),
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
	transport.NotifyCloseFunc = param.NotifyCloseFunc
	go transport.OnTimer()

	transport.logger.Debug().Msgf("newTransport options:%# v", pretty.Formatter(transport.optionsFBS))

	var err error
	if transport.optionsFBS.EnableSctp {
		transport.sctpAssociation, err = newSctpAssociation(transport.optionsFBS)
		if err != nil {
			transport.logger.Err(err).Msg("newSctpAssociation failed")
			return nil, err
		}
	}

	return transport, nil
}

func (t *Transport) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	t.logger.Debug().Str("request", request.String()).Msg("handle channel request")

	switch request.MethodType {

	case FBS__Request.MethodTRANSPORT_DUMP:
		response.Data = t.FillJson()

	case FBS__Request.MethodTRANSPORT_PRODUCE:
		requestT := request.Request.Body.Value.(*FBS__Transport.ProduceRequestT)
		if _, ok := t.mapProducers.Load(requestT.ProducerId); ok {
			response.Err = mserror.ErrProducerExist
			return
		}
		data, err := t.Produce(requestT.ProducerId, requestT)
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyTransport_ProduceResponse,
			Value: data,
		}
		response.RspBody = rspBody
		response.Err = err

	case FBS__Request.MethodTRANSPORT_CONSUME:
		var options mediasoupdata.ConsumerOptions
		_ = json.Unmarshal(request.Data, &options)
		data, err := t.Consume(request.Internal.ProducerId, request.Internal.ConsumerId, options)
		response.Data, _ = json.Marshal(data)
		response.Err = err

	case FBS__Request.MethodTRANSPORT_PRODUCE_DATA:
		requestT := request.Request.Body.Value.(*FBS__Transport.ProduceDataRequestT)
		if _, ok := t.mapDataProducers.Load(requestT.DataProducerId); ok {
			response.Err = mserror.ErrDataProducerExist
			return
		}
		dataProducer, err := t.DataProduce(requestT.DataProducerId, requestT)
		if err != nil {
			response.Err = err
			return
		}
		// set rsp
		dataDump := dataProducer.FillJson()
		_ = mediasoupdata.Clone(&response.Data, dataDump)
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyDataProducer_DumpResponse,
			Value: dataDump,
		}
		response.RspBody = rspBody

	case FBS__Request.MethodTRANSPORT_CONSUME_DATA:
		// todo
		// set rsp
		dataDump := &FBS__DataConsumer.DumpResponseT{}
		_ = mediasoupdata.Clone(&response.Data, dataDump)
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyDataConsumer_DumpResponse,
			Value: dataDump,
		}
		response.RspBody = rspBody

	case FBS__Request.MethodTRANSPORT_SET_MAX_INCOMING_BITRATE:

	case FBS__Request.MethodTRANSPORT_SET_MAX_OUTGOING_BITRATE:

	case FBS__Request.MethodTRANSPORT_ENABLE_TRACE_EVENT:

	case FBS__Request.MethodTRANSPORT_CLOSE_PRODUCER:
		requestT := request.Request.Body.Value.(*FBS__Transport.CloseProducerRequestT)
		value, ok := t.mapProducers.Load(requestT.ProducerId)
		if !ok {
			response.Err = mserror.ErrProducerNotFound
			return
		}
		producer := value.(*Producer)
		producer.Close()
		t.mapProducers.Delete(requestT.ProducerId)
		t.onTransportProducerClosedHandler(producer.id)

	case FBS__Request.MethodTRANSPORT_CLOSE_CONSUMER:
		requestT := request.Request.Body.Value.(*FBS__Transport.CloseConsumerRequestT)
		value, ok := t.mapConsumers.Load(requestT.ConsumerId)
		if !ok {
			response.Err = mserror.ErrConsumerNotFound
			return
		}
		consumer := value.(IConsumer)
		consumer.Close()
		t.mapConsumers.Delete(requestT.ConsumerId)
		for _, ssrc := range consumer.GetMediaSsrcs() {
			t.mapSsrcConsumer.Delete(ssrc)
		}
		t.onTransportConsumerClosedHandler(request.Internal.ProducerId, consumer.GetId())

	case FBS__Request.MethodTRANSPORT_CLOSE_DATAPRODUCER:
		requestT := request.Request.Body.Value.(*FBS__Transport.CloseDataProducerRequestT)
		value, ok := t.mapDataProducers.Load(requestT.DataProducerId)
		if !ok {
			response.Err = mserror.ErrDataProducerNotFound
			return
		}
		producer := value.(*DataProducer)
		producer.Close()
		t.mapDataProducers.Delete(requestT.DataProducerId)
		// todo: register
		t.onTransportDataProducerClosedHandler(producer.id)
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

	case mediasoupdata.ConsumerType_Simulcast:
		consumer, err = newSimulcastConsumer(simulcastConsumerParam{
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

	case mediasoupdata.ConsumerType_Svc:
		// todo
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

func (t *Transport) Produce(id string, request *FBS__Transport.ProduceRequestT) (*FBS__Transport.ProduceResponseT, error) {
	if id == "" {
		return nil, mserror.ErrInvalidParam
	}
	if _, ok := t.mapProducers.Load(id); ok {
		return nil, mserror.ErrDuplicatedId
	}
	producer, err := newProducer(producerParam{
		id:                                    id,
		optionsFBS:                            request,
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
	// Take the transport related RTP header extensions of the Producer and
	// add them to the Transport.
	// NOTE: Producer::GetRtpHeaderExtensionIds() returns the original
	// header extension ids of the Producer (and not their mapped values).
	t.recvRtpHeaderExtensionIds = producer.RtpHeaderExtensionIds
	t.logger.Info().Str("recvRtpHeaderExtensionIds", t.recvRtpHeaderExtensionIds.String()).
		Str("kind", producer.Kind.String()).Str("type", producer.Type.String()).Msgf("Producer created [producerId:%s]", id)

	// todo

	return &FBS__Transport.ProduceResponseT{
		Type: producer.Type,
	}, nil
}

func (t *Transport) DataProduce(id string, request *FBS__Transport.ProduceDataRequestT) (*DataProducer, error) {
	if request.DataProducerId == "" {
		return nil, mserror.ErrInvalidParam
	}
	if _, ok := t.mapDataProducers.Load(id); ok {
		return nil, mserror.ErrDuplicatedId
	}
	maxMessageSize := uint32(0)
	if t.sctpAssociation != nil {
		maxMessageSize = t.sctpAssociation.GetSctpAssociationParam().MaxMessageSize
	}
	dataProducer, err := newDataProducer(id, maxMessageSize, request)
	if err != nil {
		t.logger.Err(err).Msg("data produce failed")
		return nil, err
	}
	t.mapDataProducers.Store(id, dataProducer)
	t.logger.Debug().Msgf("DataProducer created [producerId:%s],type:%s,optionsFBS:%+v", id, request.Type, t.optionsFBS)
	return dataProducer, nil
}

func (t *Transport) ReceiveRtpPacket(packet *rtpparser.Packet) {
	// Apply the Transport RTP header extension ids so the RTP listener can use them.
	packet.SetMidExtensionId(t.recvRtpHeaderExtensionIds.Mid)
	packet.SetRidExtensionId(t.recvRtpHeaderExtensionIds.Rid)
	packet.SetRepairedRidExtensionId(t.recvRtpHeaderExtensionIds.RRid)
	packet.SetAbsSendTimeExtensionId(t.recvRtpHeaderExtensionIds.AbsSendTime)
	packet.SetTransportWideCc01ExtensionId(t.recvRtpHeaderExtensionIds.TransportWideCc01)
	// get producer from ssrc, to producer
	producer := t.rtpListener.GetProducer(packet)
	if producer == nil {
		t.logger.Warn().Str("packet", packet.String()).Str("mid", packet.GetMid()).Str("rid", packet.GetRid()).Msg("producer not found")
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
