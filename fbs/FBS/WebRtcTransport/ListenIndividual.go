// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WebRtcTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
)

type ListenIndividual struct {
	_tab flatbuffers.Table
}

func GetRootAsListenIndividual(buf []byte, offset flatbuffers.UOffsetT) *ListenIndividual {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ListenIndividual{}
	x.Init(buf, n+offset)
	return x
}

func FinishListenIndividualBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsListenIndividual(buf []byte, offset flatbuffers.UOffsetT) *ListenIndividual {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ListenIndividual{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedListenIndividualBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ListenIndividual) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ListenIndividual) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ListenIndividual) ListenInfos(obj *FBS__Transport.ListenInfo, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *ListenIndividual) ListenInfosLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ListenIndividualStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ListenIndividualAddListenInfos(builder *flatbuffers.Builder, listenInfos flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(listenInfos), 0)
}
func ListenIndividualStartListenInfosVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ListenIndividualEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
