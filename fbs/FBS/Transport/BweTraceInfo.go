// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type BweTraceInfoT struct {
	BweType BweType `json:"bwe_type"`
	DesiredBitrate uint32 `json:"desired_bitrate"`
	EffectiveDesiredBitrate uint32 `json:"effective_desired_bitrate"`
	MinBitrate uint32 `json:"min_bitrate"`
	MaxBitrate uint32 `json:"max_bitrate"`
	StartBitrate uint32 `json:"start_bitrate"`
	MaxPaddingBitrate uint32 `json:"max_padding_bitrate"`
	AvailableBitrate uint32 `json:"available_bitrate"`
}

func (t *BweTraceInfoT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	BweTraceInfoStart(builder)
	BweTraceInfoAddBweType(builder, t.BweType)
	BweTraceInfoAddDesiredBitrate(builder, t.DesiredBitrate)
	BweTraceInfoAddEffectiveDesiredBitrate(builder, t.EffectiveDesiredBitrate)
	BweTraceInfoAddMinBitrate(builder, t.MinBitrate)
	BweTraceInfoAddMaxBitrate(builder, t.MaxBitrate)
	BweTraceInfoAddStartBitrate(builder, t.StartBitrate)
	BweTraceInfoAddMaxPaddingBitrate(builder, t.MaxPaddingBitrate)
	BweTraceInfoAddAvailableBitrate(builder, t.AvailableBitrate)
	return BweTraceInfoEnd(builder)
}

func (rcv *BweTraceInfo) UnPackTo(t *BweTraceInfoT) {
	t.BweType = rcv.BweType()
	t.DesiredBitrate = rcv.DesiredBitrate()
	t.EffectiveDesiredBitrate = rcv.EffectiveDesiredBitrate()
	t.MinBitrate = rcv.MinBitrate()
	t.MaxBitrate = rcv.MaxBitrate()
	t.StartBitrate = rcv.StartBitrate()
	t.MaxPaddingBitrate = rcv.MaxPaddingBitrate()
	t.AvailableBitrate = rcv.AvailableBitrate()
}

func (rcv *BweTraceInfo) UnPack() *BweTraceInfoT {
	if rcv == nil {
		return nil
	}
	t := &BweTraceInfoT{}
	rcv.UnPackTo(t)
	return t
}

type BweTraceInfo struct {
	_tab flatbuffers.Table
}

func GetRootAsBweTraceInfo(buf []byte, offset flatbuffers.UOffsetT) *BweTraceInfo {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &BweTraceInfo{}
	x.Init(buf, n+offset)
	return x
}

func FinishBweTraceInfoBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsBweTraceInfo(buf []byte, offset flatbuffers.UOffsetT) *BweTraceInfo {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &BweTraceInfo{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedBweTraceInfoBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *BweTraceInfo) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *BweTraceInfo) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *BweTraceInfo) BweType() BweType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return BweType(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *BweTraceInfo) MutateBweType(n BweType) bool {
	return rcv._tab.MutateByteSlot(4, byte(n))
}

func (rcv *BweTraceInfo) DesiredBitrate() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *BweTraceInfo) MutateDesiredBitrate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *BweTraceInfo) EffectiveDesiredBitrate() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *BweTraceInfo) MutateEffectiveDesiredBitrate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func (rcv *BweTraceInfo) MinBitrate() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *BweTraceInfo) MutateMinBitrate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(10, n)
}

func (rcv *BweTraceInfo) MaxBitrate() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *BweTraceInfo) MutateMaxBitrate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(12, n)
}

func (rcv *BweTraceInfo) StartBitrate() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *BweTraceInfo) MutateStartBitrate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(14, n)
}

func (rcv *BweTraceInfo) MaxPaddingBitrate() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *BweTraceInfo) MutateMaxPaddingBitrate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(16, n)
}

func (rcv *BweTraceInfo) AvailableBitrate() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *BweTraceInfo) MutateAvailableBitrate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(18, n)
}

func BweTraceInfoStart(builder *flatbuffers.Builder) {
	builder.StartObject(8)
}
func BweTraceInfoAddBweType(builder *flatbuffers.Builder, bweType BweType) {
	builder.PrependByteSlot(0, byte(bweType), 0)
}
func BweTraceInfoAddDesiredBitrate(builder *flatbuffers.Builder, desiredBitrate uint32) {
	builder.PrependUint32Slot(1, desiredBitrate, 0)
}
func BweTraceInfoAddEffectiveDesiredBitrate(builder *flatbuffers.Builder, effectiveDesiredBitrate uint32) {
	builder.PrependUint32Slot(2, effectiveDesiredBitrate, 0)
}
func BweTraceInfoAddMinBitrate(builder *flatbuffers.Builder, minBitrate uint32) {
	builder.PrependUint32Slot(3, minBitrate, 0)
}
func BweTraceInfoAddMaxBitrate(builder *flatbuffers.Builder, maxBitrate uint32) {
	builder.PrependUint32Slot(4, maxBitrate, 0)
}
func BweTraceInfoAddStartBitrate(builder *flatbuffers.Builder, startBitrate uint32) {
	builder.PrependUint32Slot(5, startBitrate, 0)
}
func BweTraceInfoAddMaxPaddingBitrate(builder *flatbuffers.Builder, maxPaddingBitrate uint32) {
	builder.PrependUint32Slot(6, maxPaddingBitrate, 0)
}
func BweTraceInfoAddAvailableBitrate(builder *flatbuffers.Builder, availableBitrate uint32) {
	builder.PrependUint32Slot(7, availableBitrate, 0)
}
func BweTraceInfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
