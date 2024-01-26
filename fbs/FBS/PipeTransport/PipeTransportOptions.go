// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package PipeTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
)

type PipeTransportOptions struct {
	_tab flatbuffers.Table
}

func GetRootAsPipeTransportOptions(buf []byte, offset flatbuffers.UOffsetT) *PipeTransportOptions {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PipeTransportOptions{}
	x.Init(buf, n+offset)
	return x
}

func FinishPipeTransportOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsPipeTransportOptions(buf []byte, offset flatbuffers.UOffsetT) *PipeTransportOptions {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &PipeTransportOptions{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedPipeTransportOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *PipeTransportOptions) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PipeTransportOptions) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PipeTransportOptions) Base(obj *FBS__Transport.Options) *FBS__Transport.Options {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__Transport.Options)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *PipeTransportOptions) ListenInfo(obj *FBS__Transport.ListenInfo) *FBS__Transport.ListenInfo {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__Transport.ListenInfo)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *PipeTransportOptions) EnableRtx() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *PipeTransportOptions) MutateEnableRtx(n bool) bool {
	return rcv._tab.MutateBoolSlot(8, n)
}

func (rcv *PipeTransportOptions) EnableSrtp() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *PipeTransportOptions) MutateEnableSrtp(n bool) bool {
	return rcv._tab.MutateBoolSlot(10, n)
}

func PipeTransportOptionsStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func PipeTransportOptionsAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func PipeTransportOptionsAddListenInfo(builder *flatbuffers.Builder, listenInfo flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(listenInfo), 0)
}
func PipeTransportOptionsAddEnableRtx(builder *flatbuffers.Builder, enableRtx bool) {
	builder.PrependBoolSlot(2, enableRtx, false)
}
func PipeTransportOptionsAddEnableSrtp(builder *flatbuffers.Builder, enableSrtp bool) {
	builder.PrependBoolSlot(3, enableSrtp, false)
}
func PipeTransportOptionsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
