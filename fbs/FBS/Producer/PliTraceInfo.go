// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Producer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type PliTraceInfoT struct {
	Ssrc uint32 `json:"ssrc"`
}

func (t *PliTraceInfoT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	PliTraceInfoStart(builder)
	PliTraceInfoAddSsrc(builder, t.Ssrc)
	return PliTraceInfoEnd(builder)
}

func (rcv *PliTraceInfo) UnPackTo(t *PliTraceInfoT) {
	t.Ssrc = rcv.Ssrc()
}

func (rcv *PliTraceInfo) UnPack() *PliTraceInfoT {
	if rcv == nil {
		return nil
	}
	t := &PliTraceInfoT{}
	rcv.UnPackTo(t)
	return t
}

type PliTraceInfo struct {
	_tab flatbuffers.Table
}

func GetRootAsPliTraceInfo(buf []byte, offset flatbuffers.UOffsetT) *PliTraceInfo {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PliTraceInfo{}
	x.Init(buf, n+offset)
	return x
}

func FinishPliTraceInfoBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsPliTraceInfo(buf []byte, offset flatbuffers.UOffsetT) *PliTraceInfo {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &PliTraceInfo{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedPliTraceInfoBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *PliTraceInfo) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PliTraceInfo) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PliTraceInfo) Ssrc() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *PliTraceInfo) MutateSsrc(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func PliTraceInfoStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func PliTraceInfoAddSsrc(builder *flatbuffers.Builder, ssrc uint32) {
	builder.PrependUint32Slot(0, ssrc, 0)
}
func PliTraceInfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
