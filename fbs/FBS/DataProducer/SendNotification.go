// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DataProducer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SendNotificationT struct {
	Ppid uint32 `json:"ppid"`
	Data []byte `json:"data"`
	Subchannels []uint16 `json:"subchannels"`
	RequiredSubchannel *uint16 `json:"required_subchannel"`
}

func (t *SendNotificationT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	dataOffset := flatbuffers.UOffsetT(0)
	if t.Data != nil {
		dataOffset = builder.CreateByteString(t.Data)
	}
	subchannelsOffset := flatbuffers.UOffsetT(0)
	if t.Subchannels != nil {
		subchannelsLength := len(t.Subchannels)
		SendNotificationStartSubchannelsVector(builder, subchannelsLength)
		for j := subchannelsLength - 1; j >= 0; j-- {
			builder.PrependUint16(t.Subchannels[j])
		}
		subchannelsOffset = builder.EndVector(subchannelsLength)
	}
	SendNotificationStart(builder)
	SendNotificationAddPpid(builder, t.Ppid)
	SendNotificationAddData(builder, dataOffset)
	SendNotificationAddSubchannels(builder, subchannelsOffset)
	if t.RequiredSubchannel != nil {
		SendNotificationAddRequiredSubchannel(builder, *t.RequiredSubchannel)
	}
	return SendNotificationEnd(builder)
}

func (rcv *SendNotification) UnPackTo(t *SendNotificationT) {
	t.Ppid = rcv.Ppid()
	t.Data = rcv.DataBytes()
	subchannelsLength := rcv.SubchannelsLength()
	t.Subchannels = make([]uint16, subchannelsLength)
	for j := 0; j < subchannelsLength; j++ {
		t.Subchannels[j] = rcv.Subchannels(j)
	}
	t.RequiredSubchannel = rcv.RequiredSubchannel()
}

func (rcv *SendNotification) UnPack() *SendNotificationT {
	if rcv == nil {
		return nil
	}
	t := &SendNotificationT{}
	rcv.UnPackTo(t)
	return t
}

type SendNotification struct {
	_tab flatbuffers.Table
}

func GetRootAsSendNotification(buf []byte, offset flatbuffers.UOffsetT) *SendNotification {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SendNotification{}
	x.Init(buf, n+offset)
	return x
}

func FinishSendNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSendNotification(buf []byte, offset flatbuffers.UOffsetT) *SendNotification {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SendNotification{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSendNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SendNotification) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SendNotification) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SendNotification) Ppid() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SendNotification) MutatePpid(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *SendNotification) Data(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *SendNotification) DataLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *SendNotification) DataBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *SendNotification) MutateData(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *SendNotification) Subchannels(j int) uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint16(a + flatbuffers.UOffsetT(j*2))
	}
	return 0
}

func (rcv *SendNotification) SubchannelsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *SendNotification) MutateSubchannels(j int, n uint16) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint16(a+flatbuffers.UOffsetT(j*2), n)
	}
	return false
}

func (rcv *SendNotification) RequiredSubchannel() *uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		v := rcv._tab.GetUint16(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *SendNotification) MutateRequiredSubchannel(n uint16) bool {
	return rcv._tab.MutateUint16Slot(10, n)
}

func SendNotificationStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func SendNotificationAddPpid(builder *flatbuffers.Builder, ppid uint32) {
	builder.PrependUint32Slot(0, ppid, 0)
}
func SendNotificationAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(data), 0)
}
func SendNotificationStartDataVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func SendNotificationAddSubchannels(builder *flatbuffers.Builder, subchannels flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(subchannels), 0)
}
func SendNotificationStartSubchannelsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(2, numElems, 2)
}
func SendNotificationAddRequiredSubchannel(builder *flatbuffers.Builder, requiredSubchannel uint16) {
	builder.PrependUint16(requiredSubchannel)
	builder.Slot(3)
}
func SendNotificationEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
