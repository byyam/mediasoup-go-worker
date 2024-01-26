// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Consumer

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__RtpStream "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpStream"
)

type GetStatsResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsGetStatsResponse(buf []byte, offset flatbuffers.UOffsetT) *GetStatsResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &GetStatsResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishGetStatsResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsGetStatsResponse(buf []byte, offset flatbuffers.UOffsetT) *GetStatsResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &GetStatsResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedGetStatsResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *GetStatsResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *GetStatsResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *GetStatsResponse) Stats(obj *FBS__RtpStream.Stats, j int) bool {
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

func (rcv *GetStatsResponse) StatsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func GetStatsResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func GetStatsResponseAddStats(builder *flatbuffers.Builder, stats flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(stats), 0)
}
func GetStatsResponseStartStatsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func GetStatsResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
