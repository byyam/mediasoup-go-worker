// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Response

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"strconv"

	FBS__Consumer "github.com/byyam/mediasoup-go-worker/fbs/FBS/Consumer"
	FBS__DataConsumer "github.com/byyam/mediasoup-go-worker/fbs/FBS/DataConsumer"
	FBS__DataProducer "github.com/byyam/mediasoup-go-worker/fbs/FBS/DataProducer"
	FBS__DirectTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/DirectTransport"
	FBS__PipeTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/PipeTransport"
	FBS__PlainTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/PlainTransport"
	FBS__Producer "github.com/byyam/mediasoup-go-worker/fbs/FBS/Producer"
	FBS__Router "github.com/byyam/mediasoup-go-worker/fbs/FBS/Router"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcServer "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcServer"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	FBS__Worker "github.com/byyam/mediasoup-go-worker/fbs/FBS/Worker"
)

type Body byte

const (
	BodyNONE                                   Body = 0
	BodyWorker_DumpResponse                    Body = 1
	BodyWorker_ResourceUsageResponse           Body = 2
	BodyWebRtcServer_DumpResponse              Body = 3
	BodyRouter_DumpResponse                    Body = 4
	BodyTransport_ProduceResponse              Body = 5
	BodyTransport_ConsumeResponse              Body = 6
	BodyTransport_RestartIceResponse           Body = 7
	BodyPlainTransport_ConnectResponse         Body = 8
	BodyPlainTransport_DumpResponse            Body = 9
	BodyPlainTransport_GetStatsResponse        Body = 10
	BodyPipeTransport_ConnectResponse          Body = 11
	BodyPipeTransport_DumpResponse             Body = 12
	BodyPipeTransport_GetStatsResponse         Body = 13
	BodyDirectTransport_DumpResponse           Body = 14
	BodyDirectTransport_GetStatsResponse       Body = 15
	BodyWebRtcTransport_ConnectResponse        Body = 16
	BodyWebRtcTransport_DumpResponse           Body = 17
	BodyWebRtcTransport_GetStatsResponse       Body = 18
	BodyProducer_DumpResponse                  Body = 19
	BodyProducer_GetStatsResponse              Body = 20
	BodyConsumer_DumpResponse                  Body = 21
	BodyConsumer_GetStatsResponse              Body = 22
	BodyConsumer_SetPreferredLayersResponse    Body = 23
	BodyConsumer_SetPriorityResponse           Body = 24
	BodyDataProducer_DumpResponse              Body = 25
	BodyDataProducer_GetStatsResponse          Body = 26
	BodyDataConsumer_GetBufferedAmountResponse Body = 27
	BodyDataConsumer_DumpResponse              Body = 28
	BodyDataConsumer_GetStatsResponse          Body = 29
	BodyDataConsumer_SetSubchannelsResponse    Body = 30
	BodyDataConsumer_AddSubchannelResponse     Body = 31
	BodyDataConsumer_RemoveSubchannelResponse  Body = 32
)

var EnumNamesBody = map[Body]string{
	BodyNONE:                                   "NONE",
	BodyWorker_DumpResponse:                    "Worker_DumpResponse",
	BodyWorker_ResourceUsageResponse:           "Worker_ResourceUsageResponse",
	BodyWebRtcServer_DumpResponse:              "WebRtcServer_DumpResponse",
	BodyRouter_DumpResponse:                    "Router_DumpResponse",
	BodyTransport_ProduceResponse:              "Transport_ProduceResponse",
	BodyTransport_ConsumeResponse:              "Transport_ConsumeResponse",
	BodyTransport_RestartIceResponse:           "Transport_RestartIceResponse",
	BodyPlainTransport_ConnectResponse:         "PlainTransport_ConnectResponse",
	BodyPlainTransport_DumpResponse:            "PlainTransport_DumpResponse",
	BodyPlainTransport_GetStatsResponse:        "PlainTransport_GetStatsResponse",
	BodyPipeTransport_ConnectResponse:          "PipeTransport_ConnectResponse",
	BodyPipeTransport_DumpResponse:             "PipeTransport_DumpResponse",
	BodyPipeTransport_GetStatsResponse:         "PipeTransport_GetStatsResponse",
	BodyDirectTransport_DumpResponse:           "DirectTransport_DumpResponse",
	BodyDirectTransport_GetStatsResponse:       "DirectTransport_GetStatsResponse",
	BodyWebRtcTransport_ConnectResponse:        "WebRtcTransport_ConnectResponse",
	BodyWebRtcTransport_DumpResponse:           "WebRtcTransport_DumpResponse",
	BodyWebRtcTransport_GetStatsResponse:       "WebRtcTransport_GetStatsResponse",
	BodyProducer_DumpResponse:                  "Producer_DumpResponse",
	BodyProducer_GetStatsResponse:              "Producer_GetStatsResponse",
	BodyConsumer_DumpResponse:                  "Consumer_DumpResponse",
	BodyConsumer_GetStatsResponse:              "Consumer_GetStatsResponse",
	BodyConsumer_SetPreferredLayersResponse:    "Consumer_SetPreferredLayersResponse",
	BodyConsumer_SetPriorityResponse:           "Consumer_SetPriorityResponse",
	BodyDataProducer_DumpResponse:              "DataProducer_DumpResponse",
	BodyDataProducer_GetStatsResponse:          "DataProducer_GetStatsResponse",
	BodyDataConsumer_GetBufferedAmountResponse: "DataConsumer_GetBufferedAmountResponse",
	BodyDataConsumer_DumpResponse:              "DataConsumer_DumpResponse",
	BodyDataConsumer_GetStatsResponse:          "DataConsumer_GetStatsResponse",
	BodyDataConsumer_SetSubchannelsResponse:    "DataConsumer_SetSubchannelsResponse",
	BodyDataConsumer_AddSubchannelResponse:     "DataConsumer_AddSubchannelResponse",
	BodyDataConsumer_RemoveSubchannelResponse:  "DataConsumer_RemoveSubchannelResponse",
}

var EnumValuesBody = map[string]Body{
	"NONE":                                   BodyNONE,
	"Worker_DumpResponse":                    BodyWorker_DumpResponse,
	"Worker_ResourceUsageResponse":           BodyWorker_ResourceUsageResponse,
	"WebRtcServer_DumpResponse":              BodyWebRtcServer_DumpResponse,
	"Router_DumpResponse":                    BodyRouter_DumpResponse,
	"Transport_ProduceResponse":              BodyTransport_ProduceResponse,
	"Transport_ConsumeResponse":              BodyTransport_ConsumeResponse,
	"Transport_RestartIceResponse":           BodyTransport_RestartIceResponse,
	"PlainTransport_ConnectResponse":         BodyPlainTransport_ConnectResponse,
	"PlainTransport_DumpResponse":            BodyPlainTransport_DumpResponse,
	"PlainTransport_GetStatsResponse":        BodyPlainTransport_GetStatsResponse,
	"PipeTransport_ConnectResponse":          BodyPipeTransport_ConnectResponse,
	"PipeTransport_DumpResponse":             BodyPipeTransport_DumpResponse,
	"PipeTransport_GetStatsResponse":         BodyPipeTransport_GetStatsResponse,
	"DirectTransport_DumpResponse":           BodyDirectTransport_DumpResponse,
	"DirectTransport_GetStatsResponse":       BodyDirectTransport_GetStatsResponse,
	"WebRtcTransport_ConnectResponse":        BodyWebRtcTransport_ConnectResponse,
	"WebRtcTransport_DumpResponse":           BodyWebRtcTransport_DumpResponse,
	"WebRtcTransport_GetStatsResponse":       BodyWebRtcTransport_GetStatsResponse,
	"Producer_DumpResponse":                  BodyProducer_DumpResponse,
	"Producer_GetStatsResponse":              BodyProducer_GetStatsResponse,
	"Consumer_DumpResponse":                  BodyConsumer_DumpResponse,
	"Consumer_GetStatsResponse":              BodyConsumer_GetStatsResponse,
	"Consumer_SetPreferredLayersResponse":    BodyConsumer_SetPreferredLayersResponse,
	"Consumer_SetPriorityResponse":           BodyConsumer_SetPriorityResponse,
	"DataProducer_DumpResponse":              BodyDataProducer_DumpResponse,
	"DataProducer_GetStatsResponse":          BodyDataProducer_GetStatsResponse,
	"DataConsumer_GetBufferedAmountResponse": BodyDataConsumer_GetBufferedAmountResponse,
	"DataConsumer_DumpResponse":              BodyDataConsumer_DumpResponse,
	"DataConsumer_GetStatsResponse":          BodyDataConsumer_GetStatsResponse,
	"DataConsumer_SetSubchannelsResponse":    BodyDataConsumer_SetSubchannelsResponse,
	"DataConsumer_AddSubchannelResponse":     BodyDataConsumer_AddSubchannelResponse,
	"DataConsumer_RemoveSubchannelResponse":  BodyDataConsumer_RemoveSubchannelResponse,
}

func (v Body) String() string {
	if s, ok := EnumNamesBody[v]; ok {
		return s
	}
	return "Body(" + strconv.FormatInt(int64(v), 10) + ")"
}

type BodyT struct {
	Type Body
	Value interface{}
}

func (t *BodyT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	switch t.Type {
	case BodyWorker_DumpResponse:
		return t.Value.(*FBS__Worker.DumpResponseT).Pack(builder)
	case BodyWorker_ResourceUsageResponse:
		return t.Value.(*FBS__Worker.ResourceUsageResponseT).Pack(builder)
	case BodyWebRtcServer_DumpResponse:
		return t.Value.(*FBS__WebRtcServer.DumpResponseT).Pack(builder)
	case BodyRouter_DumpResponse:
		return t.Value.(*FBS__Router.DumpResponseT).Pack(builder)
	case BodyTransport_ProduceResponse:
		return t.Value.(*FBS__Transport.ProduceResponseT).Pack(builder)
	case BodyTransport_ConsumeResponse:
		return t.Value.(*FBS__Transport.ConsumeResponseT).Pack(builder)
	case BodyTransport_RestartIceResponse:
		return t.Value.(*FBS__Transport.RestartIceResponseT).Pack(builder)
	case BodyPlainTransport_ConnectResponse:
		return t.Value.(*FBS__PlainTransport.ConnectResponseT).Pack(builder)
	case BodyPlainTransport_DumpResponse:
		return t.Value.(*FBS__PlainTransport.DumpResponseT).Pack(builder)
	case BodyPlainTransport_GetStatsResponse:
		return t.Value.(*FBS__PlainTransport.GetStatsResponseT).Pack(builder)
	case BodyPipeTransport_ConnectResponse:
		return t.Value.(*FBS__PipeTransport.ConnectResponseT).Pack(builder)
	case BodyPipeTransport_DumpResponse:
		return t.Value.(*FBS__PipeTransport.DumpResponseT).Pack(builder)
	case BodyPipeTransport_GetStatsResponse:
		return t.Value.(*FBS__PipeTransport.GetStatsResponseT).Pack(builder)
	case BodyDirectTransport_DumpResponse:
		return t.Value.(*FBS__DirectTransport.DumpResponseT).Pack(builder)
	case BodyDirectTransport_GetStatsResponse:
		return t.Value.(*FBS__DirectTransport.GetStatsResponseT).Pack(builder)
	case BodyWebRtcTransport_ConnectResponse:
		return t.Value.(*FBS__WebRtcTransport.ConnectResponseT).Pack(builder)
	case BodyWebRtcTransport_DumpResponse:
		return t.Value.(*FBS__WebRtcTransport.DumpResponseT).Pack(builder)
	case BodyWebRtcTransport_GetStatsResponse:
		return t.Value.(*FBS__WebRtcTransport.GetStatsResponseT).Pack(builder)
	case BodyProducer_DumpResponse:
		return t.Value.(*FBS__Producer.DumpResponseT).Pack(builder)
	case BodyProducer_GetStatsResponse:
		return t.Value.(*FBS__Producer.GetStatsResponseT).Pack(builder)
	case BodyConsumer_DumpResponse:
		return t.Value.(*FBS__Consumer.DumpResponseT).Pack(builder)
	case BodyConsumer_GetStatsResponse:
		return t.Value.(*FBS__Consumer.GetStatsResponseT).Pack(builder)
	case BodyConsumer_SetPreferredLayersResponse:
		return t.Value.(*FBS__Consumer.SetPreferredLayersResponseT).Pack(builder)
	case BodyConsumer_SetPriorityResponse:
		return t.Value.(*FBS__Consumer.SetPriorityResponseT).Pack(builder)
	case BodyDataProducer_DumpResponse:
		return t.Value.(*FBS__DataProducer.DumpResponseT).Pack(builder)
	case BodyDataProducer_GetStatsResponse:
		return t.Value.(*FBS__DataProducer.GetStatsResponseT).Pack(builder)
	case BodyDataConsumer_GetBufferedAmountResponse:
		return t.Value.(*FBS__DataConsumer.GetBufferedAmountResponseT).Pack(builder)
	case BodyDataConsumer_DumpResponse:
		return t.Value.(*FBS__DataConsumer.DumpResponseT).Pack(builder)
	case BodyDataConsumer_GetStatsResponse:
		return t.Value.(*FBS__DataConsumer.GetStatsResponseT).Pack(builder)
	case BodyDataConsumer_SetSubchannelsResponse:
		return t.Value.(*FBS__DataConsumer.SetSubchannelsResponseT).Pack(builder)
	case BodyDataConsumer_AddSubchannelResponse:
		return t.Value.(*FBS__DataConsumer.AddSubchannelResponseT).Pack(builder)
	case BodyDataConsumer_RemoveSubchannelResponse:
		return t.Value.(*FBS__DataConsumer.RemoveSubchannelResponseT).Pack(builder)
	}
	return 0
}

func (rcv Body) UnPack(table flatbuffers.Table) *BodyT {
	switch rcv {
	case BodyWorker_DumpResponse:
		var x FBS__Worker.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyWorker_DumpResponse, Value: x.UnPack()}
	case BodyWorker_ResourceUsageResponse:
		var x FBS__Worker.ResourceUsageResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyWorker_ResourceUsageResponse, Value: x.UnPack()}
	case BodyWebRtcServer_DumpResponse:
		var x FBS__WebRtcServer.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyWebRtcServer_DumpResponse, Value: x.UnPack()}
	case BodyRouter_DumpResponse:
		var x FBS__Router.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyRouter_DumpResponse, Value: x.UnPack()}
	case BodyTransport_ProduceResponse:
		var x FBS__Transport.ProduceResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyTransport_ProduceResponse, Value: x.UnPack()}
	case BodyTransport_ConsumeResponse:
		var x FBS__Transport.ConsumeResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyTransport_ConsumeResponse, Value: x.UnPack()}
	case BodyTransport_RestartIceResponse:
		var x FBS__Transport.RestartIceResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyTransport_RestartIceResponse, Value: x.UnPack()}
	case BodyPlainTransport_ConnectResponse:
		var x FBS__PlainTransport.ConnectResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyPlainTransport_ConnectResponse, Value: x.UnPack()}
	case BodyPlainTransport_DumpResponse:
		var x FBS__PlainTransport.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyPlainTransport_DumpResponse, Value: x.UnPack()}
	case BodyPlainTransport_GetStatsResponse:
		var x FBS__PlainTransport.GetStatsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyPlainTransport_GetStatsResponse, Value: x.UnPack()}
	case BodyPipeTransport_ConnectResponse:
		var x FBS__PipeTransport.ConnectResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyPipeTransport_ConnectResponse, Value: x.UnPack()}
	case BodyPipeTransport_DumpResponse:
		var x FBS__PipeTransport.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyPipeTransport_DumpResponse, Value: x.UnPack()}
	case BodyPipeTransport_GetStatsResponse:
		var x FBS__PipeTransport.GetStatsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyPipeTransport_GetStatsResponse, Value: x.UnPack()}
	case BodyDirectTransport_DumpResponse:
		var x FBS__DirectTransport.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDirectTransport_DumpResponse, Value: x.UnPack()}
	case BodyDirectTransport_GetStatsResponse:
		var x FBS__DirectTransport.GetStatsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDirectTransport_GetStatsResponse, Value: x.UnPack()}
	case BodyWebRtcTransport_ConnectResponse:
		var x FBS__WebRtcTransport.ConnectResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyWebRtcTransport_ConnectResponse, Value: x.UnPack()}
	case BodyWebRtcTransport_DumpResponse:
		var x FBS__WebRtcTransport.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyWebRtcTransport_DumpResponse, Value: x.UnPack()}
	case BodyWebRtcTransport_GetStatsResponse:
		var x FBS__WebRtcTransport.GetStatsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyWebRtcTransport_GetStatsResponse, Value: x.UnPack()}
	case BodyProducer_DumpResponse:
		var x FBS__Producer.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyProducer_DumpResponse, Value: x.UnPack()}
	case BodyProducer_GetStatsResponse:
		var x FBS__Producer.GetStatsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyProducer_GetStatsResponse, Value: x.UnPack()}
	case BodyConsumer_DumpResponse:
		var x FBS__Consumer.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyConsumer_DumpResponse, Value: x.UnPack()}
	case BodyConsumer_GetStatsResponse:
		var x FBS__Consumer.GetStatsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyConsumer_GetStatsResponse, Value: x.UnPack()}
	case BodyConsumer_SetPreferredLayersResponse:
		var x FBS__Consumer.SetPreferredLayersResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyConsumer_SetPreferredLayersResponse, Value: x.UnPack()}
	case BodyConsumer_SetPriorityResponse:
		var x FBS__Consumer.SetPriorityResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyConsumer_SetPriorityResponse, Value: x.UnPack()}
	case BodyDataProducer_DumpResponse:
		var x FBS__DataProducer.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDataProducer_DumpResponse, Value: x.UnPack()}
	case BodyDataProducer_GetStatsResponse:
		var x FBS__DataProducer.GetStatsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDataProducer_GetStatsResponse, Value: x.UnPack()}
	case BodyDataConsumer_GetBufferedAmountResponse:
		var x FBS__DataConsumer.GetBufferedAmountResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDataConsumer_GetBufferedAmountResponse, Value: x.UnPack()}
	case BodyDataConsumer_DumpResponse:
		var x FBS__DataConsumer.DumpResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDataConsumer_DumpResponse, Value: x.UnPack()}
	case BodyDataConsumer_GetStatsResponse:
		var x FBS__DataConsumer.GetStatsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDataConsumer_GetStatsResponse, Value: x.UnPack()}
	case BodyDataConsumer_SetSubchannelsResponse:
		var x FBS__DataConsumer.SetSubchannelsResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDataConsumer_SetSubchannelsResponse, Value: x.UnPack()}
	case BodyDataConsumer_AddSubchannelResponse:
		var x FBS__DataConsumer.AddSubchannelResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDataConsumer_AddSubchannelResponse, Value: x.UnPack()}
	case BodyDataConsumer_RemoveSubchannelResponse:
		var x FBS__DataConsumer.RemoveSubchannelResponse
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: BodyDataConsumer_RemoveSubchannelResponse, Value: x.UnPack()}
	}
	return nil
}
