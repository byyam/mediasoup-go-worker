// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DataConsumer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SetSubchannelsRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsSetSubchannelsRequest(buf []byte, offset flatbuffers.UOffsetT) *SetSubchannelsRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SetSubchannelsRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishSetSubchannelsRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSetSubchannelsRequest(buf []byte, offset flatbuffers.UOffsetT) *SetSubchannelsRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SetSubchannelsRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSetSubchannelsRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SetSubchannelsRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SetSubchannelsRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SetSubchannelsRequest) Subchannels(j int) uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint16(a + flatbuffers.UOffsetT(j*2))
	}
	return 0
}

func (rcv *SetSubchannelsRequest) SubchannelsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *SetSubchannelsRequest) MutateSubchannels(j int, n uint16) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint16(a+flatbuffers.UOffsetT(j*2), n)
	}
	return false
}

func SetSubchannelsRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func SetSubchannelsRequestAddSubchannels(builder *flatbuffers.Builder, subchannels flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(subchannels), 0)
}
func SetSubchannelsRequestStartSubchannelsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(2, numElems, 2)
}
func SetSubchannelsRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
