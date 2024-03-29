// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Common

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type StringUint8T struct {
	Key string `json:"key"`
	Value byte `json:"value"`
}

func (t *StringUint8T) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	keyOffset := flatbuffers.UOffsetT(0)
	if t.Key != "" {
		keyOffset = builder.CreateString(t.Key)
	}
	StringUint8Start(builder)
	StringUint8AddKey(builder, keyOffset)
	StringUint8AddValue(builder, t.Value)
	return StringUint8End(builder)
}

func (rcv *StringUint8) UnPackTo(t *StringUint8T) {
	t.Key = string(rcv.Key())
	t.Value = rcv.Value()
}

func (rcv *StringUint8) UnPack() *StringUint8T {
	if rcv == nil {
		return nil
	}
	t := &StringUint8T{}
	rcv.UnPackTo(t)
	return t
}

type StringUint8 struct {
	_tab flatbuffers.Table
}

func GetRootAsStringUint8(buf []byte, offset flatbuffers.UOffsetT) *StringUint8 {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &StringUint8{}
	x.Init(buf, n+offset)
	return x
}

func FinishStringUint8Buffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsStringUint8(buf []byte, offset flatbuffers.UOffsetT) *StringUint8 {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &StringUint8{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedStringUint8Buffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *StringUint8) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *StringUint8) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *StringUint8) Key() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *StringUint8) Value() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StringUint8) MutateValue(n byte) bool {
	return rcv._tab.MutateByteSlot(6, n)
}

func StringUint8Start(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func StringUint8AddKey(builder *flatbuffers.Builder, key flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(key), 0)
}
func StringUint8AddValue(builder *flatbuffers.Builder, value byte) {
	builder.PrependByteSlot(1, value, 0)
}
func StringUint8End(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
