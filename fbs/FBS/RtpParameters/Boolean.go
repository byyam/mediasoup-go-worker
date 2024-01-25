// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpParameters

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Boolean struct {
	_tab flatbuffers.Table
}

func GetRootAsBoolean(buf []byte, offset flatbuffers.UOffsetT) *Boolean {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Boolean{}
	x.Init(buf, n+offset)
	return x
}

func FinishBooleanBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsBoolean(buf []byte, offset flatbuffers.UOffsetT) *Boolean {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Boolean{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedBooleanBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Boolean) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Boolean) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Boolean) Value() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Boolean) MutateValue(n byte) bool {
	return rcv._tab.MutateByteSlot(4, n)
}

func BooleanStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func BooleanAddValue(builder *flatbuffers.Builder, value byte) {
	builder.PrependByteSlot(0, value, 0)
}
func BooleanEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
