// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WebRtcTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "FBS/Transport"
)

type WebRtcTransportOptions struct {
	_tab flatbuffers.Table
}

func GetRootAsWebRtcTransportOptions(buf []byte, offset flatbuffers.UOffsetT) *WebRtcTransportOptions {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &WebRtcTransportOptions{}
	x.Init(buf, n+offset)
	return x
}

func FinishWebRtcTransportOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsWebRtcTransportOptions(buf []byte, offset flatbuffers.UOffsetT) *WebRtcTransportOptions {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &WebRtcTransportOptions{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedWebRtcTransportOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *WebRtcTransportOptions) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *WebRtcTransportOptions) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *WebRtcTransportOptions) Base(obj *FBS__Transport.Options) *FBS__Transport.Options {
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

func (rcv *WebRtcTransportOptions) ListenType() Listen {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return Listen(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *WebRtcTransportOptions) MutateListenType(n Listen) bool {
	return rcv._tab.MutateByteSlot(6, byte(n))
}

func (rcv *WebRtcTransportOptions) Listen(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func (rcv *WebRtcTransportOptions) EnableUdp() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return true
}

func (rcv *WebRtcTransportOptions) MutateEnableUdp(n bool) bool {
	return rcv._tab.MutateBoolSlot(10, n)
}

func (rcv *WebRtcTransportOptions) EnableTcp() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return true
}

func (rcv *WebRtcTransportOptions) MutateEnableTcp(n bool) bool {
	return rcv._tab.MutateBoolSlot(12, n)
}

func (rcv *WebRtcTransportOptions) PreferUdp() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *WebRtcTransportOptions) MutatePreferUdp(n bool) bool {
	return rcv._tab.MutateBoolSlot(14, n)
}

func (rcv *WebRtcTransportOptions) PreferTcp() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *WebRtcTransportOptions) MutatePreferTcp(n bool) bool {
	return rcv._tab.MutateBoolSlot(16, n)
}

func WebRtcTransportOptionsStart(builder *flatbuffers.Builder) {
	builder.StartObject(7)
}
func WebRtcTransportOptionsAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func WebRtcTransportOptionsAddListenType(builder *flatbuffers.Builder, listenType Listen) {
	builder.PrependByteSlot(1, byte(listenType), 0)
}
func WebRtcTransportOptionsAddListen(builder *flatbuffers.Builder, listen flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(listen), 0)
}
func WebRtcTransportOptionsAddEnableUdp(builder *flatbuffers.Builder, enableUdp bool) {
	builder.PrependBoolSlot(3, enableUdp, true)
}
func WebRtcTransportOptionsAddEnableTcp(builder *flatbuffers.Builder, enableTcp bool) {
	builder.PrependBoolSlot(4, enableTcp, true)
}
func WebRtcTransportOptionsAddPreferUdp(builder *flatbuffers.Builder, preferUdp bool) {
	builder.PrependBoolSlot(5, preferUdp, false)
}
func WebRtcTransportOptionsAddPreferTcp(builder *flatbuffers.Builder, preferTcp bool) {
	builder.PrependBoolSlot(6, preferTcp, false)
}
func WebRtcTransportOptionsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
