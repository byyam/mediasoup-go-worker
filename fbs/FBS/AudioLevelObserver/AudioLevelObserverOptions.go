// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package AudioLevelObserver

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type AudioLevelObserverOptionsT struct {
	MaxEntries uint16 `json:"max_entries"`
	Threshold int8 `json:"threshold"`
	Interval uint16 `json:"interval"`
}

func (t *AudioLevelObserverOptionsT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	AudioLevelObserverOptionsStart(builder)
	AudioLevelObserverOptionsAddMaxEntries(builder, t.MaxEntries)
	AudioLevelObserverOptionsAddThreshold(builder, t.Threshold)
	AudioLevelObserverOptionsAddInterval(builder, t.Interval)
	return AudioLevelObserverOptionsEnd(builder)
}

func (rcv *AudioLevelObserverOptions) UnPackTo(t *AudioLevelObserverOptionsT) {
	t.MaxEntries = rcv.MaxEntries()
	t.Threshold = rcv.Threshold()
	t.Interval = rcv.Interval()
}

func (rcv *AudioLevelObserverOptions) UnPack() *AudioLevelObserverOptionsT {
	if rcv == nil {
		return nil
	}
	t := &AudioLevelObserverOptionsT{}
	rcv.UnPackTo(t)
	return t
}

type AudioLevelObserverOptions struct {
	_tab flatbuffers.Table
}

func GetRootAsAudioLevelObserverOptions(buf []byte, offset flatbuffers.UOffsetT) *AudioLevelObserverOptions {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &AudioLevelObserverOptions{}
	x.Init(buf, n+offset)
	return x
}

func FinishAudioLevelObserverOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsAudioLevelObserverOptions(buf []byte, offset flatbuffers.UOffsetT) *AudioLevelObserverOptions {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &AudioLevelObserverOptions{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedAudioLevelObserverOptionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *AudioLevelObserverOptions) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *AudioLevelObserverOptions) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *AudioLevelObserverOptions) MaxEntries() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AudioLevelObserverOptions) MutateMaxEntries(n uint16) bool {
	return rcv._tab.MutateUint16Slot(4, n)
}

func (rcv *AudioLevelObserverOptions) Threshold() int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AudioLevelObserverOptions) MutateThreshold(n int8) bool {
	return rcv._tab.MutateInt8Slot(6, n)
}

func (rcv *AudioLevelObserverOptions) Interval() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AudioLevelObserverOptions) MutateInterval(n uint16) bool {
	return rcv._tab.MutateUint16Slot(8, n)
}

func AudioLevelObserverOptionsStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func AudioLevelObserverOptionsAddMaxEntries(builder *flatbuffers.Builder, maxEntries uint16) {
	builder.PrependUint16Slot(0, maxEntries, 0)
}
func AudioLevelObserverOptionsAddThreshold(builder *flatbuffers.Builder, threshold int8) {
	builder.PrependInt8Slot(1, threshold, 0)
}
func AudioLevelObserverOptionsAddInterval(builder *flatbuffers.Builder, interval uint16) {
	builder.PrependUint16Slot(2, interval, 0)
}
func AudioLevelObserverOptionsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
