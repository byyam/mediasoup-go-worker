// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package PlainTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__SrtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/SrtpParameters"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
)

type PlainTransportOptions struct {
	_tab flatbuffers.Table
}

func GetRootAsPlainTransportOptions(buf []byte, offset flatbuffers.UOffsetT) *PlainTransportOptions {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PlainTransportOptions{}
	x.Init(buf, n+offset)
	return x
}

func FinishPlainTransportOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsPlainTransportOptions(buf []byte, offset flatbuffers.UOffsetT) *PlainTransportOptions {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &PlainTransportOptions{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedPlainTransportOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *PlainTransportOptions) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PlainTransportOptions) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PlainTransportOptions) Base(obj *FBS__Transport.Options) *FBS__Transport.Options {
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

func (rcv *PlainTransportOptions) ListenInfo(obj *FBS__Transport.ListenInfo) *FBS__Transport.ListenInfo {
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

func (rcv *PlainTransportOptions) RtcpListenInfo(obj *FBS__Transport.ListenInfo) *FBS__Transport.ListenInfo {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
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

func (rcv *PlainTransportOptions) RtcpMux() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *PlainTransportOptions) MutateRtcpMux(n bool) bool {
	return rcv._tab.MutateBoolSlot(10, n)
}

func (rcv *PlainTransportOptions) Comedia() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *PlainTransportOptions) MutateComedia(n bool) bool {
	return rcv._tab.MutateBoolSlot(12, n)
}

func (rcv *PlainTransportOptions) EnableSrtp() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *PlainTransportOptions) MutateEnableSrtp(n bool) bool {
	return rcv._tab.MutateBoolSlot(14, n)
}

func (rcv *PlainTransportOptions) SrtpCryptoSuite() *FBS__SrtpParameters.SrtpCryptoSuite {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		v := FBS__SrtpParameters.SrtpCryptoSuite(rcv._tab.GetByte(o + rcv._tab.Pos))
		return &v
	}
	return nil
}

func (rcv *PlainTransportOptions) MutateSrtpCryptoSuite(n FBS__SrtpParameters.SrtpCryptoSuite) bool {
	return rcv._tab.MutateByteSlot(16, byte(n))
}

func PlainTransportOptionsStart(builder *flatbuffers.Builder) {
	builder.StartObject(7)
}
func PlainTransportOptionsAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func PlainTransportOptionsAddListenInfo(builder *flatbuffers.Builder, listenInfo flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(listenInfo), 0)
}
func PlainTransportOptionsAddRtcpListenInfo(builder *flatbuffers.Builder, rtcpListenInfo flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(rtcpListenInfo), 0)
}
func PlainTransportOptionsAddRtcpMux(builder *flatbuffers.Builder, rtcpMux bool) {
	builder.PrependBoolSlot(3, rtcpMux, false)
}
func PlainTransportOptionsAddComedia(builder *flatbuffers.Builder, comedia bool) {
	builder.PrependBoolSlot(4, comedia, false)
}
func PlainTransportOptionsAddEnableSrtp(builder *flatbuffers.Builder, enableSrtp bool) {
	builder.PrependBoolSlot(5, enableSrtp, false)
}
func PlainTransportOptionsAddSrtpCryptoSuite(builder *flatbuffers.Builder, srtpCryptoSuite FBS__SrtpParameters.SrtpCryptoSuite) {
	builder.PrependByte(byte(srtpCryptoSuite))
	builder.Slot(6)
}
func PlainTransportOptionsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
