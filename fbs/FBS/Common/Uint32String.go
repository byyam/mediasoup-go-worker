// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Common

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Uint32StringT struct {
	Key uint32 `json:"key"`
	Value string `json:"value"`
}

func (t *Uint32StringT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	valueOffset := flatbuffers.UOffsetT(0)
	if t.Value != "" {
		valueOffset = builder.CreateString(t.Value)
	}
	Uint32StringStart(builder)
	Uint32StringAddKey(builder, t.Key)
	Uint32StringAddValue(builder, valueOffset)
	return Uint32StringEnd(builder)
}

func (rcv *Uint32String) UnPackTo(t *Uint32StringT) {
	t.Key = rcv.Key()
	t.Value = string(rcv.Value())
}

func (rcv *Uint32String) UnPack() *Uint32StringT {
	if rcv == nil {
		return nil
	}
	t := &Uint32StringT{}
	rcv.UnPackTo(t)
	return t
}

type Uint32String struct {
	_tab flatbuffers.Table
}

func GetRootAsUint32String(buf []byte, offset flatbuffers.UOffsetT) *Uint32String {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Uint32String{}
	x.Init(buf, n+offset)
	return x
}

func FinishUint32StringBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsUint32String(buf []byte, offset flatbuffers.UOffsetT) *Uint32String {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Uint32String{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedUint32StringBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Uint32String) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Uint32String) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Uint32String) Key() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Uint32String) MutateKey(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *Uint32String) Value() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func Uint32StringStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func Uint32StringAddKey(builder *flatbuffers.Builder, key uint32) {
	builder.PrependUint32Slot(0, key, 0)
}
func Uint32StringAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(value), 0)
}
func Uint32StringEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
