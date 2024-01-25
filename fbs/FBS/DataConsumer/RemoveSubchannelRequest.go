// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DataConsumer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RemoveSubchannelRequestT struct {
	Subchannel uint16 `json:"subchannel"`
}

func (t *RemoveSubchannelRequestT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	RemoveSubchannelRequestStart(builder)
	RemoveSubchannelRequestAddSubchannel(builder, t.Subchannel)
	return RemoveSubchannelRequestEnd(builder)
}

func (rcv *RemoveSubchannelRequest) UnPackTo(t *RemoveSubchannelRequestT) {
	t.Subchannel = rcv.Subchannel()
}

func (rcv *RemoveSubchannelRequest) UnPack() *RemoveSubchannelRequestT {
	if rcv == nil {
		return nil
	}
	t := &RemoveSubchannelRequestT{}
	rcv.UnPackTo(t)
	return t
}

type RemoveSubchannelRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsRemoveSubchannelRequest(buf []byte, offset flatbuffers.UOffsetT) *RemoveSubchannelRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RemoveSubchannelRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishRemoveSubchannelRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRemoveSubchannelRequest(buf []byte, offset flatbuffers.UOffsetT) *RemoveSubchannelRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RemoveSubchannelRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRemoveSubchannelRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *RemoveSubchannelRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RemoveSubchannelRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RemoveSubchannelRequest) Subchannel() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RemoveSubchannelRequest) MutateSubchannel(n uint16) bool {
	return rcv._tab.MutateUint16Slot(4, n)
}

func RemoveSubchannelRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func RemoveSubchannelRequestAddSubchannel(builder *flatbuffers.Builder, subchannel uint16) {
	builder.PrependUint16Slot(0, subchannel, 0)
}
func RemoveSubchannelRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}