// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Consumer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SetPreferredLayersRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsSetPreferredLayersRequest(buf []byte, offset flatbuffers.UOffsetT) *SetPreferredLayersRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SetPreferredLayersRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishSetPreferredLayersRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSetPreferredLayersRequest(buf []byte, offset flatbuffers.UOffsetT) *SetPreferredLayersRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SetPreferredLayersRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSetPreferredLayersRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SetPreferredLayersRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SetPreferredLayersRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SetPreferredLayersRequest) PreferredLayers(obj *ConsumerLayers) *ConsumerLayers {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(ConsumerLayers)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func SetPreferredLayersRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func SetPreferredLayersRequestAddPreferredLayers(builder *flatbuffers.Builder, preferredLayers flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(preferredLayers), 0)
}
func SetPreferredLayersRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
