// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DirectTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "FBS/Transport"
)

type DirectTransportOptions struct {
	_tab flatbuffers.Table
}

func GetRootAsDirectTransportOptions(buf []byte, offset flatbuffers.UOffsetT) *DirectTransportOptions {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DirectTransportOptions{}
	x.Init(buf, n+offset)
	return x
}

func FinishDirectTransportOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsDirectTransportOptions(buf []byte, offset flatbuffers.UOffsetT) *DirectTransportOptions {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &DirectTransportOptions{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedDirectTransportOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *DirectTransportOptions) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DirectTransportOptions) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *DirectTransportOptions) Base(obj *FBS__Transport.Options) *FBS__Transport.Options {
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

func DirectTransportOptionsStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func DirectTransportOptionsAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func DirectTransportOptionsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
