// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Consumer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SetPriorityRequestT struct {
	Priority byte `json:"priority"`
}

func (t *SetPriorityRequestT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	SetPriorityRequestStart(builder)
	SetPriorityRequestAddPriority(builder, t.Priority)
	return SetPriorityRequestEnd(builder)
}

func (rcv *SetPriorityRequest) UnPackTo(t *SetPriorityRequestT) {
	t.Priority = rcv.Priority()
}

func (rcv *SetPriorityRequest) UnPack() *SetPriorityRequestT {
	if rcv == nil {
		return nil
	}
	t := &SetPriorityRequestT{}
	rcv.UnPackTo(t)
	return t
}

type SetPriorityRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsSetPriorityRequest(buf []byte, offset flatbuffers.UOffsetT) *SetPriorityRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SetPriorityRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishSetPriorityRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSetPriorityRequest(buf []byte, offset flatbuffers.UOffsetT) *SetPriorityRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SetPriorityRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSetPriorityRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SetPriorityRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SetPriorityRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SetPriorityRequest) Priority() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SetPriorityRequest) MutatePriority(n byte) bool {
	return rcv._tab.MutateByteSlot(4, n)
}

func SetPriorityRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func SetPriorityRequestAddPriority(builder *flatbuffers.Builder, priority byte) {
	builder.PrependByteSlot(0, priority, 0)
}
func SetPriorityRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
