// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Consumer

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__RtpPacket "FBS/RtpPacket"
)

type KeyFrameTraceInfo struct {
	_tab flatbuffers.Table
}

func GetRootAsKeyFrameTraceInfo(buf []byte, offset flatbuffers.UOffsetT) *KeyFrameTraceInfo {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &KeyFrameTraceInfo{}
	x.Init(buf, n+offset)
	return x
}

func FinishKeyFrameTraceInfoBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsKeyFrameTraceInfo(buf []byte, offset flatbuffers.UOffsetT) *KeyFrameTraceInfo {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &KeyFrameTraceInfo{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedKeyFrameTraceInfoBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *KeyFrameTraceInfo) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *KeyFrameTraceInfo) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *KeyFrameTraceInfo) RtpPacket(obj *FBS__RtpPacket.Dump) *FBS__RtpPacket.Dump {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__RtpPacket.Dump)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *KeyFrameTraceInfo) IsRtx() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *KeyFrameTraceInfo) MutateIsRtx(n bool) bool {
	return rcv._tab.MutateBoolSlot(6, n)
}

func KeyFrameTraceInfoStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func KeyFrameTraceInfoAddRtpPacket(builder *flatbuffers.Builder, rtpPacket flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(rtpPacket), 0)
}
func KeyFrameTraceInfoAddIsRtx(builder *flatbuffers.Builder, isRtx bool) {
	builder.PrependBoolSlot(1, isRtx, false)
}
func KeyFrameTraceInfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
