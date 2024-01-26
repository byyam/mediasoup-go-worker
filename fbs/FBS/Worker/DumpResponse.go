// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Worker

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__LibUring "github.com/byyam/mediasoup-go-worker/fbs/FBS/LibUring"
)

type DumpResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsDumpResponse(buf []byte, offset flatbuffers.UOffsetT) *DumpResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DumpResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishDumpResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsDumpResponse(buf []byte, offset flatbuffers.UOffsetT) *DumpResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &DumpResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedDumpResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *DumpResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DumpResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *DumpResponse) Pid() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *DumpResponse) MutatePid(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *DumpResponse) WebRtcServerIds(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *DumpResponse) WebRtcServerIdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) RouterIds(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *DumpResponse) RouterIdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) ChannelMessageHandlers(obj *ChannelMessageHandlers) *ChannelMessageHandlers {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(ChannelMessageHandlers)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *DumpResponse) Liburing(obj *FBS__LibUring.Dump) *FBS__LibUring.Dump {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__LibUring.Dump)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func DumpResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func DumpResponseAddPid(builder *flatbuffers.Builder, pid uint32) {
	builder.PrependUint32Slot(0, pid, 0)
}
func DumpResponseAddWebRtcServerIds(builder *flatbuffers.Builder, webRtcServerIds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(webRtcServerIds), 0)
}
func DumpResponseStartWebRtcServerIdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddRouterIds(builder *flatbuffers.Builder, routerIds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(routerIds), 0)
}
func DumpResponseStartRouterIdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddChannelMessageHandlers(builder *flatbuffers.Builder, channelMessageHandlers flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(channelMessageHandlers), 0)
}
func DumpResponseAddLiburing(builder *flatbuffers.Builder, liburing flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(liburing), 0)
}
func DumpResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
