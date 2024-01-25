// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpPacket

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Dump struct {
	_tab flatbuffers.Table
}

func GetRootAsDump(buf []byte, offset flatbuffers.UOffsetT) *Dump {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Dump{}
	x.Init(buf, n+offset)
	return x
}

func FinishDumpBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsDump(buf []byte, offset flatbuffers.UOffsetT) *Dump {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Dump{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedDumpBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Dump) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Dump) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Dump) PayloadType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Dump) MutatePayloadType(n byte) bool {
	return rcv._tab.MutateByteSlot(4, n)
}

func (rcv *Dump) SequenceNumber() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Dump) MutateSequenceNumber(n uint16) bool {
	return rcv._tab.MutateUint16Slot(6, n)
}

func (rcv *Dump) Timestamp() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Dump) MutateTimestamp(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func (rcv *Dump) Marker() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Dump) MutateMarker(n bool) bool {
	return rcv._tab.MutateBoolSlot(10, n)
}

func (rcv *Dump) Ssrc() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Dump) MutateSsrc(n uint32) bool {
	return rcv._tab.MutateUint32Slot(12, n)
}

func (rcv *Dump) IsKeyFrame() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Dump) MutateIsKeyFrame(n bool) bool {
	return rcv._tab.MutateBoolSlot(14, n)
}

func (rcv *Dump) Size() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Dump) MutateSize(n uint64) bool {
	return rcv._tab.MutateUint64Slot(16, n)
}

func (rcv *Dump) PayloadSize() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Dump) MutatePayloadSize(n uint64) bool {
	return rcv._tab.MutateUint64Slot(18, n)
}

func (rcv *Dump) SpatialLayer() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Dump) MutateSpatialLayer(n byte) bool {
	return rcv._tab.MutateByteSlot(20, n)
}

func (rcv *Dump) TemporalLayer() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Dump) MutateTemporalLayer(n byte) bool {
	return rcv._tab.MutateByteSlot(22, n)
}

func (rcv *Dump) Mid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Dump) Rid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Dump) Rrid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Dump) WideSequenceNumber() *uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		v := rcv._tab.GetUint16(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *Dump) MutateWideSequenceNumber(n uint16) bool {
	return rcv._tab.MutateUint16Slot(30, n)
}

func DumpStart(builder *flatbuffers.Builder) {
	builder.StartObject(14)
}
func DumpAddPayloadType(builder *flatbuffers.Builder, payloadType byte) {
	builder.PrependByteSlot(0, payloadType, 0)
}
func DumpAddSequenceNumber(builder *flatbuffers.Builder, sequenceNumber uint16) {
	builder.PrependUint16Slot(1, sequenceNumber, 0)
}
func DumpAddTimestamp(builder *flatbuffers.Builder, timestamp uint32) {
	builder.PrependUint32Slot(2, timestamp, 0)
}
func DumpAddMarker(builder *flatbuffers.Builder, marker bool) {
	builder.PrependBoolSlot(3, marker, false)
}
func DumpAddSsrc(builder *flatbuffers.Builder, ssrc uint32) {
	builder.PrependUint32Slot(4, ssrc, 0)
}
func DumpAddIsKeyFrame(builder *flatbuffers.Builder, isKeyFrame bool) {
	builder.PrependBoolSlot(5, isKeyFrame, false)
}
func DumpAddSize(builder *flatbuffers.Builder, size uint64) {
	builder.PrependUint64Slot(6, size, 0)
}
func DumpAddPayloadSize(builder *flatbuffers.Builder, payloadSize uint64) {
	builder.PrependUint64Slot(7, payloadSize, 0)
}
func DumpAddSpatialLayer(builder *flatbuffers.Builder, spatialLayer byte) {
	builder.PrependByteSlot(8, spatialLayer, 0)
}
func DumpAddTemporalLayer(builder *flatbuffers.Builder, temporalLayer byte) {
	builder.PrependByteSlot(9, temporalLayer, 0)
}
func DumpAddMid(builder *flatbuffers.Builder, mid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(10, flatbuffers.UOffsetT(mid), 0)
}
func DumpAddRid(builder *flatbuffers.Builder, rid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(11, flatbuffers.UOffsetT(rid), 0)
}
func DumpAddRrid(builder *flatbuffers.Builder, rrid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(12, flatbuffers.UOffsetT(rrid), 0)
}
func DumpAddWideSequenceNumber(builder *flatbuffers.Builder, wideSequenceNumber uint16) {
	builder.PrependUint16(wideSequenceNumber)
	builder.Slot(13)
}
func DumpEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
