// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Router

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__PipeTransport "FBS/PipeTransport"
)

type CreatePipeTransportRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsCreatePipeTransportRequest(buf []byte, offset flatbuffers.UOffsetT) *CreatePipeTransportRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &CreatePipeTransportRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishCreatePipeTransportRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsCreatePipeTransportRequest(buf []byte, offset flatbuffers.UOffsetT) *CreatePipeTransportRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &CreatePipeTransportRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedCreatePipeTransportRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *CreatePipeTransportRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CreatePipeTransportRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *CreatePipeTransportRequest) TransportId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *CreatePipeTransportRequest) Options(obj *FBS__PipeTransport.PipeTransportOptions) *FBS__PipeTransport.PipeTransportOptions {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__PipeTransport.PipeTransportOptions)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func CreatePipeTransportRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func CreatePipeTransportRequestAddTransportId(builder *flatbuffers.Builder, transportId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(transportId), 0)
}
func CreatePipeTransportRequestAddOptions(builder *flatbuffers.Builder, options flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(options), 0)
}
func CreatePipeTransportRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
