// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DataConsumer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RemoveSubchannelResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsRemoveSubchannelResponse(buf []byte, offset flatbuffers.UOffsetT) *RemoveSubchannelResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RemoveSubchannelResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishRemoveSubchannelResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRemoveSubchannelResponse(buf []byte, offset flatbuffers.UOffsetT) *RemoveSubchannelResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RemoveSubchannelResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRemoveSubchannelResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *RemoveSubchannelResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RemoveSubchannelResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RemoveSubchannelResponse) Subchannels(j int) uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint16(a + flatbuffers.UOffsetT(j*2))
	}
	return 0
}

func (rcv *RemoveSubchannelResponse) SubchannelsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RemoveSubchannelResponse) MutateSubchannels(j int, n uint16) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint16(a+flatbuffers.UOffsetT(j*2), n)
	}
	return false
}

func RemoveSubchannelResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func RemoveSubchannelResponseAddSubchannels(builder *flatbuffers.Builder, subchannels flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(subchannels), 0)
}
func RemoveSubchannelResponseStartSubchannelsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(2, numElems, 2)
}
func RemoveSubchannelResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
