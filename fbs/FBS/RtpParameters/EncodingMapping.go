// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpParameters

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type EncodingMappingT struct {
	Rid string `json:"rid"`
	Ssrc *uint32 `json:"ssrc"`
	ScalabilityMode string `json:"scalability_mode"`
	MappedSsrc uint32 `json:"mapped_ssrc"`
}

func (t *EncodingMappingT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	ridOffset := flatbuffers.UOffsetT(0)
	if t.Rid != "" {
		ridOffset = builder.CreateString(t.Rid)
	}
	scalabilityModeOffset := flatbuffers.UOffsetT(0)
	if t.ScalabilityMode != "" {
		scalabilityModeOffset = builder.CreateString(t.ScalabilityMode)
	}
	EncodingMappingStart(builder)
	EncodingMappingAddRid(builder, ridOffset)
	if t.Ssrc != nil {
		EncodingMappingAddSsrc(builder, *t.Ssrc)
	}
	EncodingMappingAddScalabilityMode(builder, scalabilityModeOffset)
	EncodingMappingAddMappedSsrc(builder, t.MappedSsrc)
	return EncodingMappingEnd(builder)
}

func (rcv *EncodingMapping) UnPackTo(t *EncodingMappingT) {
	t.Rid = string(rcv.Rid())
	t.Ssrc = rcv.Ssrc()
	t.ScalabilityMode = string(rcv.ScalabilityMode())
	t.MappedSsrc = rcv.MappedSsrc()
}

func (rcv *EncodingMapping) UnPack() *EncodingMappingT {
	if rcv == nil {
		return nil
	}
	t := &EncodingMappingT{}
	rcv.UnPackTo(t)
	return t
}

type EncodingMapping struct {
	_tab flatbuffers.Table
}

func GetRootAsEncodingMapping(buf []byte, offset flatbuffers.UOffsetT) *EncodingMapping {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &EncodingMapping{}
	x.Init(buf, n+offset)
	return x
}

func FinishEncodingMappingBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsEncodingMapping(buf []byte, offset flatbuffers.UOffsetT) *EncodingMapping {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &EncodingMapping{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedEncodingMappingBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *EncodingMapping) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *EncodingMapping) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *EncodingMapping) Rid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *EncodingMapping) Ssrc() *uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		v := rcv._tab.GetUint32(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *EncodingMapping) MutateSsrc(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *EncodingMapping) ScalabilityMode() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *EncodingMapping) MappedSsrc() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *EncodingMapping) MutateMappedSsrc(n uint32) bool {
	return rcv._tab.MutateUint32Slot(10, n)
}

func EncodingMappingStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func EncodingMappingAddRid(builder *flatbuffers.Builder, rid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(rid), 0)
}
func EncodingMappingAddSsrc(builder *flatbuffers.Builder, ssrc uint32) {
	builder.PrependUint32(ssrc)
	builder.Slot(1)
}
func EncodingMappingAddScalabilityMode(builder *flatbuffers.Builder, scalabilityMode flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(scalabilityMode), 0)
}
func EncodingMappingAddMappedSsrc(builder *flatbuffers.Builder, mappedSsrc uint32) {
	builder.PrependUint32Slot(3, mappedSsrc, 0)
}
func EncodingMappingEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
