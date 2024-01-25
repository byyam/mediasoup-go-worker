// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpParameters

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RtxT struct {
	Ssrc uint32 `json:"ssrc"`
}

func (t *RtxT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	RtxStart(builder)
	RtxAddSsrc(builder, t.Ssrc)
	return RtxEnd(builder)
}

func (rcv *Rtx) UnPackTo(t *RtxT) {
	t.Ssrc = rcv.Ssrc()
}

func (rcv *Rtx) UnPack() *RtxT {
	if rcv == nil {
		return nil
	}
	t := &RtxT{}
	rcv.UnPackTo(t)
	return t
}

type Rtx struct {
	_tab flatbuffers.Table
}

func GetRootAsRtx(buf []byte, offset flatbuffers.UOffsetT) *Rtx {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Rtx{}
	x.Init(buf, n+offset)
	return x
}

func FinishRtxBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRtx(buf []byte, offset flatbuffers.UOffsetT) *Rtx {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Rtx{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRtxBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Rtx) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Rtx) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Rtx) Ssrc() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Rtx) MutateSsrc(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func RtxStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func RtxAddSsrc(builder *flatbuffers.Builder, ssrc uint32) {
	builder.PrependUint32Slot(0, ssrc, 0)
}
func RtxEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
