// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpParameters

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type CodecMappingT struct {
	PayloadType byte `json:"payload_type"`
	MappedPayloadType byte `json:"mapped_payload_type"`
}

func (t *CodecMappingT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	CodecMappingStart(builder)
	CodecMappingAddPayloadType(builder, t.PayloadType)
	CodecMappingAddMappedPayloadType(builder, t.MappedPayloadType)
	return CodecMappingEnd(builder)
}

func (rcv *CodecMapping) UnPackTo(t *CodecMappingT) {
	t.PayloadType = rcv.PayloadType()
	t.MappedPayloadType = rcv.MappedPayloadType()
}

func (rcv *CodecMapping) UnPack() *CodecMappingT {
	if rcv == nil {
		return nil
	}
	t := &CodecMappingT{}
	rcv.UnPackTo(t)
	return t
}

type CodecMapping struct {
	_tab flatbuffers.Table
}

func GetRootAsCodecMapping(buf []byte, offset flatbuffers.UOffsetT) *CodecMapping {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &CodecMapping{}
	x.Init(buf, n+offset)
	return x
}

func FinishCodecMappingBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsCodecMapping(buf []byte, offset flatbuffers.UOffsetT) *CodecMapping {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &CodecMapping{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedCodecMappingBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *CodecMapping) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CodecMapping) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *CodecMapping) PayloadType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *CodecMapping) MutatePayloadType(n byte) bool {
	return rcv._tab.MutateByteSlot(4, n)
}

func (rcv *CodecMapping) MappedPayloadType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *CodecMapping) MutateMappedPayloadType(n byte) bool {
	return rcv._tab.MutateByteSlot(6, n)
}

func CodecMappingStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func CodecMappingAddPayloadType(builder *flatbuffers.Builder, payloadType byte) {
	builder.PrependByteSlot(0, payloadType, 0)
}
func CodecMappingAddMappedPayloadType(builder *flatbuffers.Builder, mappedPayloadType byte) {
	builder.PrependByteSlot(1, mappedPayloadType, 0)
}
func CodecMappingEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
