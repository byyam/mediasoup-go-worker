package rtc

import (
	"encoding/json"
	"sync"
	"sync/atomic"

	"github.com/pion/rtp"

	"github.com/byyam/mediasoup-go-worker/common"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type ITransport interface {
	Close()
	FillJson() json.RawMessage
	HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData)
	ReceiveRtpPacket(packet *rtp.Packet)
}

type Transport struct {
	id     string
	logger utils.Logger

	mapProducers sync.Map //map[string]*Producer
	rtpListener  *RtpListener

	// handler
	onTransportNewProducerHandler               atomic.Value
	onTransportProducerRtpPacketReceivedHandler atomic.Value // (producer *Producer, packet *rtp.Packet)
	onTransportNewConsumerHandler               atomic.Value
	sendRtpPacketFunc                           func(packet *rtp.Packet)
}

func (t *Transport) Close() {
}

func (t *Transport) FillJson() json.RawMessage {
	dumpData := mediasoupdata.TransportDump{
		Id:                      t.id,
		Direct:                  false,
		ProducerIds:             nil,
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
	OnTransportNewProducer               func(producer *Producer) error
	OnTransportProducerRtpPacketReceived func(producer *Producer, packet *rtp.Packet)
	OnTransportNewConsumer               func(consumer IConsumer, producerId string) error
	SendRtpPacketFunc                    func(packet *rtp.Packet) // call webrtcTransport
}

func (t transportParam) valid() bool {
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
		logger:      utils.NewLogger("transport"),
		rtpListener: newRtpListener(),
	}
	transport.onTransportNewProducerHandler.Store(param.OnTransportNewProducer)
	transport.onTransportProducerRtpPacketReceivedHandler.Store(param.OnTransportProducerRtpPacketReceived)
	transport.onTransportNewConsumerHandler.Store(param.OnTransportNewConsumer)
	transport.sendRtpPacketFunc = param.SendRtpPacketFunc
	return transport, nil
}
func (t *Transport) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	t.logger.Debug("method=%s,internal=%+v", request.Method, request.InternalData)

	switch request.Method {

	case mediasoupdata.MethodTransportDump:

	case mediasoupdata.MethodTransportClose:

	case mediasoupdata.MethodTransportProduce:
		var options mediasoupdata.ProducerOptions
		_ = json.Unmarshal(request.Data, &options)
		data, err := t.Produce(request.InternalData.ProducerId, options)
		if err == nil {
			response.Data, _ = json.Marshal(data)
		}
		response.Err = err

	case mediasoupdata.MethodTransportConsume:
		var options mediasoupdata.ConsumerOptions
		_ = json.Unmarshal(request.Data, &options)
		data, err := t.Consume(request.InternalData.ProducerId, request.InternalData.ConsumerId, options)
		if err == nil {
			response.Data, _ = json.Marshal(data)
		}
		response.Err = err

	case mediasoupdata.MethodTransportProduceData:

	case mediasoupdata.MethodTransportConsumeData:

	case mediasoupdata.MethodTransportSetMaxIncomingBitrate:

	case mediasoupdata.MethodTransportSetMaxOutgoingBitrate:

	case mediasoupdata.MethodTransportEnableTraceEvent:

	case mediasoupdata.MethodTransportGetStats:

	default:
		t.logger.Error("unknown method:%s", request.Method)
	}
	t.logger.Debug("method:%s, response:%s", request.Method, response)
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
				id:            consumerId,
				producerId:    producerId,
				rtpParameters: options.RtpParameters,
			},
			OnConsumerSendRtpPacket: t.OnConsumerSendRtpPacket,
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
	if handler, ok := t.onTransportNewConsumerHandler.Load().(func(consumer IConsumer, producerId string) error); ok && handler != nil {
		if err := handler(consumer, producerId); err != nil {
			return nil, err
		}
	}

	t.logger.Debug("Consumer created [producerId:%s][consumerId:%s],type:%s", producerId, consumerId, options.Type)
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
	t.logger.Debug("OnConsumerSendRtpPacket:%+v", packet.Header)
	t.sendRtpPacketFunc(packet)
}
