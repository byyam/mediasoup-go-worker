// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Router

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__PlainTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/PlainTransport"
)

type CreatePlainTransportRequestT struct {
	TransportId string `json:"transport_id"`
	Options *FBS__PlainTransport.PlainTransportOptionsT `json:"options"`
}

func (t *CreatePlainTransportRequestT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	transportIdOffset := flatbuffers.UOffsetT(0)
	if t.TransportId != "" {
		transportIdOffset = builder.CreateString(t.TransportId)
	}
	optionsOffset := t.Options.Pack(builder)
	CreatePlainTransportRequestStart(builder)
	CreatePlainTransportRequestAddTransportId(builder, transportIdOffset)
	CreatePlainTransportRequestAddOptions(builder, optionsOffset)
	return CreatePlainTransportRequestEnd(builder)
}

func (rcv *CreatePlainTransportRequest) UnPackTo(t *CreatePlainTransportRequestT) {
	t.TransportId = string(rcv.TransportId())
	t.Options = rcv.Options(nil).UnPack()
}

func (rcv *CreatePlainTransportRequest) UnPack() *CreatePlainTransportRequestT {
	if rcv == nil {
		return nil
	}
	t := &CreatePlainTransportRequestT{}
	rcv.UnPackTo(t)
	return t
}

type CreatePlainTransportRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsCreatePlainTransportRequest(buf []byte, offset flatbuffers.UOffsetT) *CreatePlainTransportRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &CreatePlainTransportRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishCreatePlainTransportRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsCreatePlainTransportRequest(buf []byte, offset flatbuffers.UOffsetT) *CreatePlainTransportRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &CreatePlainTransportRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedCreatePlainTransportRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *CreatePlainTransportRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CreatePlainTransportRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *CreatePlainTransportRequest) TransportId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *CreatePlainTransportRequest) Options(obj *FBS__PlainTransport.PlainTransportOptions) *FBS__PlainTransport.PlainTransportOptions {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__PlainTransport.PlainTransportOptions)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func CreatePlainTransportRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func CreatePlainTransportRequestAddTransportId(builder *flatbuffers.Builder, transportId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(transportId), 0)
}
func CreatePlainTransportRequestAddOptions(builder *flatbuffers.Builder, options flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(options), 0)
}
func CreatePlainTransportRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
