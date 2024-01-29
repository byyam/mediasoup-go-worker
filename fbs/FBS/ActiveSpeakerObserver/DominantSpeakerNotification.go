// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package ActiveSpeakerObserver

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type DominantSpeakerNotificationT struct {
	ProducerId string `json:"producer_id"`
}

func (t *DominantSpeakerNotificationT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	producerIdOffset := flatbuffers.UOffsetT(0)
	if t.ProducerId != "" {
		producerIdOffset = builder.CreateString(t.ProducerId)
	}
	DominantSpeakerNotificationStart(builder)
	DominantSpeakerNotificationAddProducerId(builder, producerIdOffset)
	return DominantSpeakerNotificationEnd(builder)
}

func (rcv *DominantSpeakerNotification) UnPackTo(t *DominantSpeakerNotificationT) {
	t.ProducerId = string(rcv.ProducerId())
}

func (rcv *DominantSpeakerNotification) UnPack() *DominantSpeakerNotificationT {
	if rcv == nil {
		return nil
	}
	t := &DominantSpeakerNotificationT{}
	rcv.UnPackTo(t)
	return t
}

type DominantSpeakerNotification struct {
	_tab flatbuffers.Table
}

func GetRootAsDominantSpeakerNotification(buf []byte, offset flatbuffers.UOffsetT) *DominantSpeakerNotification {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DominantSpeakerNotification{}
	x.Init(buf, n+offset)
	return x
}

func FinishDominantSpeakerNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsDominantSpeakerNotification(buf []byte, offset flatbuffers.UOffsetT) *DominantSpeakerNotification {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &DominantSpeakerNotification{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedDominantSpeakerNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *DominantSpeakerNotification) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DominantSpeakerNotification) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *DominantSpeakerNotification) ProducerId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func DominantSpeakerNotificationStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func DominantSpeakerNotificationAddProducerId(builder *flatbuffers.Builder, producerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(producerId), 0)
}
func DominantSpeakerNotificationEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
