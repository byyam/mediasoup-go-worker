package request_fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
)

type BodyT struct {
	Type  Request.Body
	Value interface{}
}

func (t *BodyT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	//switch t.Type {
	//case Request.BodyWorker_UpdateSettingsRequest:
	//	return t.Value.(*FBS__Worker.UpdateSettingsRequestT).Pack(builder)
	//case Request.BodyWorker_CreateWebRtcServerRequest:
	//	return t.Value.(*FBS__Worker.CreateWebRtcServerRequestT).Pack(builder)
	//case Request.BodyWorker_CloseWebRtcServerRequest:
	//	return t.Value.(*FBS__Worker.CloseWebRtcServerRequestT).Pack(builder)
	//case Request.BodyWorker_CreateRouterRequest:
	//	return t.Value.(*FBS__Worker.CreateRouterRequestT).Pack(builder)
	//case Request.BodyWorker_CloseRouterRequest:
	//	return t.Value.(*FBS__Worker.CloseRouterRequestT).Pack(builder)
	//case Request.BodyRouter_CreateWebRtcTransportRequest:
	//	return t.Value.(*FBS__Router.CreateWebRtcTransportRequestT).Pack(builder)
	//case Request.BodyRouter_CreatePlainTransportRequest:
	//	return t.Value.(*FBS__Router.CreatePlainTransportRequestT).Pack(builder)
	//case Request.BodyRouter_CreatePipeTransportRequest:
	//	return t.Value.(*FBS__Router.CreatePipeTransportRequestT).Pack(builder)
	//case Request.BodyRouter_CreateDirectTransportRequest:
	//	return t.Value.(*FBS__Router.CreateDirectTransportRequestT).Pack(builder)
	//case Request.BodyRouter_CreateActiveSpeakerObserverRequest:
	//	return t.Value.(*FBS__Router.CreateActiveSpeakerObserverRequestT).Pack(builder)
	//case Request.BodyRouter_CreateAudioLevelObserverRequest:
	//	return t.Value.(*FBS__Router.CreateAudioLevelObserverRequestT).Pack(builder)
	//case Request.BodyRouter_CloseTransportRequest:
	//	return t.Value.(*FBS__Router.CloseTransportRequestT).Pack(builder)
	//case Request.BodyRouter_CloseRtpObserverRequest:
	//	return t.Value.(*FBS__Router.CloseRtpObserverRequestT).Pack(builder)
	//case Request.BodyTransport_SetMaxIncomingBitrateRequest:
	//	return t.Value.(*FBS__Transport.SetMaxIncomingBitrateRequestT).Pack(builder)
	//case Request.BodyTransport_SetMaxOutgoingBitrateRequest:
	//	return t.Value.(*FBS__Transport.SetMaxOutgoingBitrateRequestT).Pack(builder)
	//case Request.BodyTransport_SetMinOutgoingBitrateRequest:
	//	return t.Value.(*FBS__Transport.SetMinOutgoingBitrateRequestT).Pack(builder)
	//case Request.BodyTransport_ProduceRequest:
	//	return t.Value.(*FBS__Transport.ProduceRequestT).Pack(builder)
	//case Request.BodyTransport_ConsumeRequest:
	//	return t.Value.(*FBS__Transport.ConsumeRequestT).Pack(builder)
	//case Request.BodyTransport_ProduceDataRequest:
	//	return t.Value.(*FBS__Transport.ProduceDataRequestT).Pack(builder)
	//case Request.BodyTransport_ConsumeDataRequest:
	//	return t.Value.(*FBS__Transport.ConsumeDataRequestT).Pack(builder)
	//case Request.BodyTransport_EnableTraceEventRequest:
	//	return t.Value.(*FBS__Transport.EnableTraceEventRequestT).Pack(builder)
	//case Request.BodyTransport_CloseProducerRequest:
	//	return t.Value.(*FBS__Transport.CloseProducerRequestT).Pack(builder)
	//case Request.BodyTransport_CloseConsumerRequest:
	//	return t.Value.(*FBS__Transport.CloseConsumerRequestT).Pack(builder)
	//case Request.BodyTransport_CloseDataProducerRequest:
	//	return t.Value.(*FBS__Transport.CloseDataProducerRequestT).Pack(builder)
	//case Request.BodyTransport_CloseDataConsumerRequest:
	//	return t.Value.(*FBS__Transport.CloseDataConsumerRequestT).Pack(builder)
	//case Request.BodyPlainTransport_ConnectRequest:
	//	return t.Value.(*FBS__PlainTransport.ConnectRequestT).Pack(builder)
	//case Request.BodyPipeTransport_ConnectRequest:
	//	return t.Value.(*FBS__PipeTransport.ConnectRequestT).Pack(builder)
	//case Request.BodyWebRtcTransport_ConnectRequest:
	//	return t.Value.(*FBS__WebRtcTransport.ConnectRequestT).Pack(builder)
	//case Request.BodyProducer_EnableTraceEventRequest:
	//	return t.Value.(*FBS__Producer.EnableTraceEventRequestT).Pack(builder)
	//case Request.BodyConsumer_SetPreferredLayersRequest:
	//	return t.Value.(*FBS__Consumer.SetPreferredLayersRequestT).Pack(builder)
	//case Request.BodyConsumer_SetPriorityRequest:
	//	return t.Value.(*FBS__Consumer.SetPriorityRequestT).Pack(builder)
	//case Request.BodyConsumer_EnableTraceEventRequest:
	//	return t.Value.(*FBS__Consumer.EnableTraceEventRequestT).Pack(builder)
	//case Request.BodyDataConsumer_SetBufferedAmountLowThresholdRequest:
	//	return t.Value.(*FBS__DataConsumer.SetBufferedAmountLowThresholdRequestT).Pack(builder)
	//case Request.BodyDataConsumer_SendRequest:
	//	return t.Value.(*FBS__DataConsumer.SendRequestT).Pack(builder)
	//case Request.BodyDataConsumer_SetSubchannelsRequest:
	//	return t.Value.(*FBS__DataConsumer.SetSubchannelsRequestT).Pack(builder)
	//case Request.BodyDataConsumer_AddSubchannelRequest:
	//	return t.Value.(*FBS__DataConsumer.AddSubchannelRequestT).Pack(builder)
	//case Request.BodyDataConsumer_RemoveSubchannelRequest:
	//	return t.Value.(*FBS__DataConsumer.RemoveSubchannelRequestT).Pack(builder)
	//case Request.BodyRtpObserver_AddProducerRequest:
	//	return t.Value.(*FBS__RtpObserver.AddProducerRequestT).Pack(builder)
	//case Request.BodyRtpObserver_RemoveProducerRequest:
	//	return t.Value.(*FBS__RtpObserver.RemoveProducerRequestT).Pack(builder)
	//}
	return 0
}

func BodyUnPack(rcv Request.Body, table flatbuffers.Table) *BodyT {
	//switch rcv {
	//case Request.BodyWorker_UpdateSettingsRequest:
	//	var x FBS__Worker.UpdateSettingsRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyWorker_UpdateSettingsRequest, Value: x.UnPack()}
	//case Request.BodyWorker_CreateWebRtcServerRequest:
	//	var x FBS__Worker.CreateWebRtcServerRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyWorker_CreateWebRtcServerRequest, Value: x.UnPack()}
	//case Request.BodyWorker_CloseWebRtcServerRequest:
	//	var x FBS__Worker.CloseWebRtcServerRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyWorker_CloseWebRtcServerRequest, Value: x.UnPack()}
	//case Request.BodyWorker_CreateRouterRequest:
	//	var x FBS__Worker.CreateRouterRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyWorker_CreateRouterRequest, Value: x.UnPack()}
	//case Request.BodyWorker_CloseRouterRequest:
	//	var x FBS__Worker.CloseRouterRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyWorker_CloseRouterRequest, Value: x.UnPack()}
	//case Request.BodyRouter_CreateWebRtcTransportRequest:
	//	var x FBS__Router.CreateWebRtcTransportRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRouter_CreateWebRtcTransportRequest, Value: x.UnPack()}
	//case Request.BodyRouter_CreatePlainTransportRequest:
	//	var x FBS__Router.CreatePlainTransportRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRouter_CreatePlainTransportRequest, Value: x.UnPack()}
	//case Request.BodyRouter_CreatePipeTransportRequest:
	//	var x FBS__Router.CreatePipeTransportRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRouter_CreatePipeTransportRequest, Value: x.UnPack()}
	//case Request.BodyRouter_CreateDirectTransportRequest:
	//	var x FBS__Router.CreateDirectTransportRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRouter_CreateDirectTransportRequest, Value: x.UnPack()}
	//case Request.BodyRouter_CreateActiveSpeakerObserverRequest:
	//	var x FBS__Router.CreateActiveSpeakerObserverRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRouter_CreateActiveSpeakerObserverRequest, Value: x.UnPack()}
	//case Request.BodyRouter_CreateAudioLevelObserverRequest:
	//	var x FBS__Router.CreateAudioLevelObserverRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRouter_CreateAudioLevelObserverRequest, Value: x.UnPack()}
	//case Request.BodyRouter_CloseTransportRequest:
	//	var x FBS__Router.CloseTransportRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRouter_CloseTransportRequest, Value: x.UnPack()}
	//case Request.BodyRouter_CloseRtpObserverRequest:
	//	var x FBS__Router.CloseRtpObserverRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRouter_CloseRtpObserverRequest, Value: x.UnPack()}
	//case Request.BodyTransport_SetMaxIncomingBitrateRequest:
	//	var x FBS__Transport.SetMaxIncomingBitrateRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_SetMaxIncomingBitrateRequest, Value: x.UnPack()}
	//case Request.BodyTransport_SetMaxOutgoingBitrateRequest:
	//	var x FBS__Transport.SetMaxOutgoingBitrateRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_SetMaxOutgoingBitrateRequest, Value: x.UnPack()}
	//case Request.BodyTransport_SetMinOutgoingBitrateRequest:
	//	var x FBS__Transport.SetMinOutgoingBitrateRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_SetMinOutgoingBitrateRequest, Value: x.UnPack()}
	//case Request.BodyTransport_ProduceRequest:
	//	var x FBS__Transport.ProduceRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_ProduceRequest, Value: x.UnPack()}
	//case Request.BodyTransport_ConsumeRequest:
	//	var x FBS__Transport.ConsumeRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_ConsumeRequest, Value: x.UnPack()}
	//case Request.BodyTransport_ProduceDataRequest:
	//	var x FBS__Transport.ProduceDataRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_ProduceDataRequest, Value: x.UnPack()}
	//case Request.BodyTransport_ConsumeDataRequest:
	//	var x FBS__Transport.ConsumeDataRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_ConsumeDataRequest, Value: x.UnPack()}
	//case Request.BodyTransport_EnableTraceEventRequest:
	//	var x FBS__Transport.EnableTraceEventRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_EnableTraceEventRequest, Value: x.UnPack()}
	//case Request.BodyTransport_CloseProducerRequest:
	//	var x FBS__Transport.CloseProducerRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_CloseProducerRequest, Value: x.UnPack()}
	//case Request.BodyTransport_CloseConsumerRequest:
	//	var x FBS__Transport.CloseConsumerRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_CloseConsumerRequest, Value: x.UnPack()}
	//case Request.BodyTransport_CloseDataProducerRequest:
	//	var x FBS__Transport.CloseDataProducerRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_CloseDataProducerRequest, Value: x.UnPack()}
	//case Request.BodyTransport_CloseDataConsumerRequest:
	//	var x FBS__Transport.CloseDataConsumerRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyTransport_CloseDataConsumerRequest, Value: x.UnPack()}
	//case Request.BodyPlainTransport_ConnectRequest:
	//	var x FBS__PlainTransport.ConnectRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyPlainTransport_ConnectRequest, Value: x.UnPack()}
	//case Request.BodyPipeTransport_ConnectRequest:
	//	var x FBS__PipeTransport.ConnectRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyPipeTransport_ConnectRequest, Value: x.UnPack()}
	//case Request.BodyWebRtcTransport_ConnectRequest:
	//	var x FBS__WebRtcTransport.ConnectRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyWebRtcTransport_ConnectRequest, Value: x.UnPack()}
	//case Request.BodyProducer_EnableTraceEventRequest:
	//	var x FBS__Producer.EnableTraceEventRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyProducer_EnableTraceEventRequest, Value: x.UnPack()}
	//case Request.BodyConsumer_SetPreferredLayersRequest:
	//	var x FBS__Consumer.SetPreferredLayersRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyConsumer_SetPreferredLayersRequest, Value: x.UnPack()}
	//case Request.BodyConsumer_SetPriorityRequest:
	//	var x FBS__Consumer.SetPriorityRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyConsumer_SetPriorityRequest, Value: x.UnPack()}
	//case Request.BodyConsumer_EnableTraceEventRequest:
	//	var x FBS__Consumer.EnableTraceEventRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyConsumer_EnableTraceEventRequest, Value: x.UnPack()}
	//case Request.BodyDataConsumer_SetBufferedAmountLowThresholdRequest:
	//	var x FBS__DataConsumer.SetBufferedAmountLowThresholdRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyDataConsumer_SetBufferedAmountLowThresholdRequest, Value: x.UnPack()}
	//case Request.BodyDataConsumer_SendRequest:
	//	var x FBS__DataConsumer.SendRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyDataConsumer_SendRequest, Value: x.UnPack()}
	//case Request.BodyDataConsumer_SetSubchannelsRequest:
	//	var x FBS__DataConsumer.SetSubchannelsRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyDataConsumer_SetSubchannelsRequest, Value: x.UnPack()}
	//case Request.BodyDataConsumer_AddSubchannelRequest:
	//	var x FBS__DataConsumer.AddSubchannelRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyDataConsumer_AddSubchannelRequest, Value: x.UnPack()}
	//case Request.BodyDataConsumer_RemoveSubchannelRequest:
	//	var x FBS__DataConsumer.RemoveSubchannelRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyDataConsumer_RemoveSubchannelRequest, Value: x.UnPack()}
	//case Request.BodyRtpObserver_AddProducerRequest:
	//	var x FBS__RtpObserver.AddProducerRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRtpObserver_AddProducerRequest, Value: x.UnPack()}
	//case Request.BodyRtpObserver_RemoveProducerRequest:
	//	var x FBS__RtpObserver.RemoveProducerRequest
	//	x.Init(table.Bytes, table.Pos)
	//	return &BodyT{Type: Request.BodyRtpObserver_RemoveProducerRequest, Value: x.UnPack()}
	//}
	return nil
}
