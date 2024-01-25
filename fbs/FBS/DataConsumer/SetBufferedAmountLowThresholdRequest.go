// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DataConsumer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SetBufferedAmountLowThresholdRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsSetBufferedAmountLowThresholdRequest(buf []byte, offset flatbuffers.UOffsetT) *SetBufferedAmountLowThresholdRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SetBufferedAmountLowThresholdRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishSetBufferedAmountLowThresholdRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSetBufferedAmountLowThresholdRequest(buf []byte, offset flatbuffers.UOffsetT) *SetBufferedAmountLowThresholdRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SetBufferedAmountLowThresholdRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSetBufferedAmountLowThresholdRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SetBufferedAmountLowThresholdRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SetBufferedAmountLowThresholdRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SetBufferedAmountLowThresholdRequest) Threshold() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SetBufferedAmountLowThresholdRequest) MutateThreshold(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func SetBufferedAmountLowThresholdRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func SetBufferedAmountLowThresholdRequestAddThreshold(builder *flatbuffers.Builder, threshold uint32) {
	builder.PrependUint32Slot(0, threshold, 0)
}
func SetBufferedAmountLowThresholdRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
