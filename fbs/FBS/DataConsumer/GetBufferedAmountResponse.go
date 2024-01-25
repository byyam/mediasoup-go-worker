// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DataConsumer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type GetBufferedAmountResponseT struct {
	BufferedAmount uint32 `json:"buffered_amount"`
}

func (t *GetBufferedAmountResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	GetBufferedAmountResponseStart(builder)
	GetBufferedAmountResponseAddBufferedAmount(builder, t.BufferedAmount)
	return GetBufferedAmountResponseEnd(builder)
}

func (rcv *GetBufferedAmountResponse) UnPackTo(t *GetBufferedAmountResponseT) {
	t.BufferedAmount = rcv.BufferedAmount()
}

func (rcv *GetBufferedAmountResponse) UnPack() *GetBufferedAmountResponseT {
	if rcv == nil {
		return nil
	}
	t := &GetBufferedAmountResponseT{}
	rcv.UnPackTo(t)
	return t
}

type GetBufferedAmountResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsGetBufferedAmountResponse(buf []byte, offset flatbuffers.UOffsetT) *GetBufferedAmountResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &GetBufferedAmountResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishGetBufferedAmountResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsGetBufferedAmountResponse(buf []byte, offset flatbuffers.UOffsetT) *GetBufferedAmountResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &GetBufferedAmountResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedGetBufferedAmountResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *GetBufferedAmountResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *GetBufferedAmountResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *GetBufferedAmountResponse) BufferedAmount() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *GetBufferedAmountResponse) MutateBufferedAmount(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func GetBufferedAmountResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func GetBufferedAmountResponseAddBufferedAmount(builder *flatbuffers.Builder, bufferedAmount uint32) {
	builder.PrependUint32Slot(0, bufferedAmount, 0)
}
func GetBufferedAmountResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
