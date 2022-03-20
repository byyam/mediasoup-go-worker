package rtc

import (
	"encoding/json"
	"sync"
	"sync/atomic"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/pion/rtcp"

	"github.com/pion/rtp"

	"github.com/byyam/mediasoup-go-worker/common"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type ITransport interface {
	Close()
	FillJson() json.RawMessage
	HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData)
	ReceiveRtpPacket(packet *rtp.Packet)
	ReceiveRtcpPacket(header *rtcp.Header, packets []rtcp.Packet)
}

type Transport struct {
	id     string
	logger utils.Logger

	mapProducers    sync.Map //map[string]*Producer
	mapConsumers    sync.Map
	mapSsrcConsumer sync.Map
	rtpListener     *RtpListener

	// handler
	onTransportNewProducerHandler               atomic.Value
	onTransportProducerClosedHandler            func(producerId string)
	onTransportProducerRtpPacketReceivedHandler atomic.Value // (producer *Producer, packet *rtp.Packet)
	onTransportNewConsumerHandler               func(consumer IConsumer, producerId string) error
	onTransportConsumerClosedHandler            func(producerId, consumerId string)
	onTransportConsumerKeyFrameRequestedHandler func(consumerId string, mappedSsrc uint32)

	// transport base call sons
	sendRtpPacketFunc func(packet *rtp.Packet)
}

func (t *Transport) Close() {
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
	t.logger.Debug("dumpData:%+v", dumpData)
	return data
}

type transportParam struct {
	Id                                   string
	OnTransportNewProducer               func(producer *Producer) error
	OnTransportProducerClosed            func(producerId string)
	OnTransportProducerRtpPacketReceived func(producer *Producer, packet *rtp.Packet)
	OnTransportNewConsumer               func(consumer IConsumer, producerId string) error
	OnTransportConsumerClosed            func(consumerId, producerId string)
	OnTransportConsumerKeyFrameRequested func(consumerId string, mappedSsrc uint32)
	SendRtpPacketFunc                    func(packet *rtp.Packet) // call webrtcTransport
}

func (t transportParam) valid() bool {
	if t.Id == "" {
		return false
	}
	if t.OnTransportNewProducer == nil || t.OnTransportProducerRtpPacketReceived == nil {
		return false
	}
	if t.SendRtpPacketFunc == nil {
		return false
	}
	return true
}

func newTransport(param transportParam) (ITransport, error) {
	if !param.valid() {
		return nil, common.ErrInvalidParam
	}
	transport := &Transport{
		id:          param.Id,
		logger:      utils.NewLogger("transport"),
		rtpListener: newRtpListener(),
	}
	transport.onTransportNewProducerHandler.Store(param.OnTransportNewProducer)
	transport.onTransportProducerClosedHandler = param.OnTransportProducerClosed
	transport.onTransportProducerRtpPacketReceivedHandler.Store(param.OnTransportProducerRtpPacketReceived)
	transport.onTransportNewConsumerHandler = param.OnTransportNewConsumer
	transport.onTransportConsumerClosedHandler = param.OnTransportConsumerClosed
	transport.onTransportConsumerKeyFrameRequestedHandler = param.OnTransportConsumerKeyFrameRequested
	transport.sendRtpPacketFunc = param.SendRtpPacketFunc
	return transport, nil
}
func (t *Transport) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		t.logger.Debug("method=%s,internal=%+v,response:%s", request.Method, request.Internal, response)
	}()

	switch request.Method {

	case mediasoupdata.MethodTransportDump:
		response.Data = t.FillJson()

	case mediasoupdata.MethodTransportClose:

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

	case mediasoupdata.MethodTransportConsumeData:

	case mediasoupdata.MethodTransportSetMaxIncomingBitrate:

	case mediasoupdata.MethodTransportSetMaxOutgoingBitrate:

	case mediasoupdata.MethodTransportEnableTraceEvent:

	case mediasoupdata.MethodTransportGetStats:

	// producer
	case mediasoupdata.MethodProducerDump, mediasoupdata.MethodProducerGetStats, mediasoupdata.MethodProducerPause,
		mediasoupdata.MethodProducerResume, mediasoupdata.MethodProducerEnableTraceEvent:
		value, ok := t.mapProducers.Load(request.Internal.ProducerId)
		if !ok {
			response.Err = common.ErrProducerNotFound
			return
		}
		producer := value.(*Producer)
		producer.HandleRequest(request, response)

	case mediasoupdata.MethodProducerClose:
		value, ok := t.mapProducers.Load(request.Internal.ProducerId)
		if !ok {
			response.Err = common.ErrProducerNotFound
			return
		}
		producer := value.(*Producer)
		producer.Close()
		t.mapProducers.Delete(request.Internal.ProducerId)
		t.onTransportProducerClosedHandler(producer.id)

	case mediasoupdata.MethodConsumerClose:
		value, ok := t.mapConsumers.Load(request.Internal.ConsumerId)
		if !ok {
			response.Err = common.ErrConsumerNotFound
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
		t.logger.Error("unknown method:%s", request.Method)
	}
}

func (t *Transport) Consume(producerId, consumerId string, options mediasoupdata.ConsumerOptions) (*mediasoupdata.ConsumerData, error) {
	if producerId == "" || consumerId == "" {
		return nil, common.ErrInvalidParam
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
			OnConsumerSendRtpPacket:     t.OnConsumerSendRtpPacket,
			OnConsumerKeyFrameRequested: t.OnConsumerKeyFrameRequested,
		})

	case mediasoupdata.ConsumerType_Simulcast: // todo...
	case mediasoupdata.ConsumerType_Svc:
	default:
		return nil, common.ErrInvalidParam
	}

	if err != nil {
		t.logger.Error("create consumer[%s] failed:%v", options.Type, err)
		return nil, err
	}
	if err := t.onTransportNewConsumerHandler(consumer, producerId); err != nil {
		return nil, err
	}

	t.mapConsumers.Store(consumerId, consumer)
	for _, ssrc := range consumer.GetMediaSsrcs() {
		t.mapSsrcConsumer.Store(ssrc, consumer)
	}
	t.logger.Debug("Consumer created [producerId:%s][consumerId:%s],type:%s,ssrc:%v", producerId, consumerId, options.Type, consumer.GetMediaSsrcs())
	return &mediasoupdata.ConsumerData{
		Paused:         false,
		ProducerPaused: false,
		Score:          mediasoupdata.ConsumerScore{},
	}, nil
}

func (t *Transport) Produce(id string, options mediasoupdata.ProducerOptions) (*mediasoupdata.ProducerData, error) {
	if id == "" {
		return nil, common.ErrInvalidParam
	}
	if _, ok := t.mapProducers.Load(id); ok {
		return nil, common.ErrDuplicatedId
	}
	producer, err := newProducer(producerParam{
		id:                          id,
		options:                     options,
		OnProducerRtpPacketReceived: t.OnProducerRtpPacketReceived,
	})
	if err != nil {
		return nil, err
	}
	t.rtpListener.AddProducer(producer)
	if handler, ok := t.onTransportNewProducerHandler.Load().(func(*Producer) error); ok && handler != nil {
		if err := handler(producer); err != nil {
			return nil, err
		}
	}
	t.mapProducers.Store(id, producer)
	t.logger.Debug("Producer created [producerId:%s],type:%s", id, producer.Type)
	// todo

	return &mediasoupdata.ProducerData{Type: producer.Type}, nil
}

func (t *Transport) ReceiveRtpPacket(packet *rtp.Packet) {
	// get producer from ssrc, to producer
	producer := t.rtpListener.GetProducer(packet)
	if producer == nil {
		monitor.RtpRecvCount(monitor.ActionSsrcNotFound)
		return
	}

	producer.ReceiveRtpPacket(packet)
}

func (t *Transport) OnProducerRtpPacketReceived(producer *Producer, packet *rtp.Packet) {
	if handler, ok := t.onTransportProducerRtpPacketReceivedHandler.Load().(func(*Producer, *rtp.Packet)); ok && handler != nil {
		handler(producer, packet)
	}
}

func (t *Transport) OnConsumerSendRtpPacket(consumer IConsumer, packet *rtp.Packet) {
	t.logger.Trace("OnConsumerSendRtpPacket:%+v", packet.Header)
	t.sendRtpPacketFunc(packet)
}

func (t *Transport) ReceiveRtcpPacket(header *rtcp.Header, packets []rtcp.Packet) {
	c := rtcp.CompoundPacket(packets)
	t.logger.Info("ReceiveRtcpPacket:\n%+v\nheader:%+v\nCompoundPacket:%v", c.String(), header, c.Validate())
	for _, packet := range packets {
		t.HandleRtcpPacket(header, packet)
	}
}

func (t *Transport) HandleRtcpPacket(header *rtcp.Header, packet rtcp.Packet) {
	t.logger.Debug("HandleRtcpPacket:%v", packet.DestinationSSRC())
	switch packet.(type) {
	case *rtcp.SenderReport:
		pkg := packet.(*rtcp.SenderReport)
		t.logger.Debug("%s", pkg.String())
	case *rtcp.ReceiverReport:
		pkg := packet.(*rtcp.ReceiverReport)
		t.logger.Debug("%s", pkg.String())
	case *rtcp.SourceDescription:
		pkg := packet.(*rtcp.SourceDescription)
		t.logger.Debug("%s", pkg.String())
	case *rtcp.Goodbye:
		pkg := packet.(*rtcp.Goodbye)
		t.logger.Debug("%s", pkg.String())
	case *rtcp.FullIntraRequest:
		pkg := packet.(*rtcp.FullIntraRequest)
		t.ReceiveKeyFrameRequest(header.Count, pkg.MediaSSRC)
		t.logger.Debug("%s", pkg.String())
	case *rtcp.PictureLossIndication:
		pkg := packet.(*rtcp.PictureLossIndication)
		t.ReceiveKeyFrameRequest(header.Count, pkg.MediaSSRC)
		t.logger.Debug("%s", pkg.String())
	case *rtcp.ReceiverEstimatedMaximumBitrate:
		pkg := packet.(*rtcp.ReceiverEstimatedMaximumBitrate)
		t.logger.Debug("%s", pkg.String())
	default:
		t.logger.Warn("unknown rtcp header:%+v", header)
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
