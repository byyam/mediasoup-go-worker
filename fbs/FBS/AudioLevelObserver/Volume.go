// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package AudioLevelObserver

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Volume struct {
	_tab flatbuffers.Table
}

func GetRootAsVolume(buf []byte, offset flatbuffers.UOffsetT) *Volume {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Volume{}
	x.Init(buf, n+offset)
	return x
}

func FinishVolumeBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsVolume(buf []byte, offset flatbuffers.UOffsetT) *Volume {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Volume{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedVolumeBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Volume) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Volume) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Volume) ProducerId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Volume) Volume() int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Volume) MutateVolume(n int8) bool {
	return rcv._tab.MutateInt8Slot(6, n)
}

func VolumeStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func VolumeAddProducerId(builder *flatbuffers.Builder, producerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(producerId), 0)
}
func VolumeAddVolume(builder *flatbuffers.Builder, volume int8) {
	builder.PrependInt8Slot(1, volume, 0)
}
func VolumeEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
