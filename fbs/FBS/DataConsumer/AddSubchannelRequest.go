// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DataConsumer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type AddSubchannelRequestT struct {
	Subchannel uint16 `json:"subchannel"`
}

func (t *AddSubchannelRequestT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	AddSubchannelRequestStart(builder)
	AddSubchannelRequestAddSubchannel(builder, t.Subchannel)
	return AddSubchannelRequestEnd(builder)
}

func (rcv *AddSubchannelRequest) UnPackTo(t *AddSubchannelRequestT) {
	t.Subchannel = rcv.Subchannel()
}

func (rcv *AddSubchannelRequest) UnPack() *AddSubchannelRequestT {
	if rcv == nil {
		return nil
	}
	t := &AddSubchannelRequestT{}
	rcv.UnPackTo(t)
	return t
}

type AddSubchannelRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsAddSubchannelRequest(buf []byte, offset flatbuffers.UOffsetT) *AddSubchannelRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &AddSubchannelRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishAddSubchannelRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsAddSubchannelRequest(buf []byte, offset flatbuffers.UOffsetT) *AddSubchannelRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &AddSubchannelRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedAddSubchannelRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *AddSubchannelRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *AddSubchannelRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *AddSubchannelRequest) Subchannel() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AddSubchannelRequest) MutateSubchannel(n uint16) bool {
	return rcv._tab.MutateUint16Slot(4, n)
}

func AddSubchannelRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func AddSubchannelRequestAddSubchannel(builder *flatbuffers.Builder, subchannel uint16) {
	builder.PrependUint16Slot(0, subchannel, 0)
}
func AddSubchannelRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
