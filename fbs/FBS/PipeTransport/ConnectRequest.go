// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package PipeTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__SrtpParameters "FBS/SrtpParameters"
)

type ConnectRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsConnectRequest(buf []byte, offset flatbuffers.UOffsetT) *ConnectRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ConnectRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishConnectRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConnectRequest(buf []byte, offset flatbuffers.UOffsetT) *ConnectRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ConnectRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConnectRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ConnectRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ConnectRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ConnectRequest) Ip() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ConnectRequest) Port() *uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		v := rcv._tab.GetUint16(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *ConnectRequest) MutatePort(n uint16) bool {
	return rcv._tab.MutateUint16Slot(6, n)
}

func (rcv *ConnectRequest) SrtpParameters(obj *FBS__SrtpParameters.SrtpParameters) *FBS__SrtpParameters.SrtpParameters {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__SrtpParameters.SrtpParameters)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func ConnectRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func ConnectRequestAddIp(builder *flatbuffers.Builder, ip flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(ip), 0)
}
func ConnectRequestAddPort(builder *flatbuffers.Builder, port uint16) {
	builder.PrependUint16(port)
	builder.Slot(1)
}
func ConnectRequestAddSrtpParameters(builder *flatbuffers.Builder, srtpParameters flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(srtpParameters), 0)
}
func ConnectRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
