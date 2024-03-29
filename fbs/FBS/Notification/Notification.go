// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Notification

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type NotificationT struct {
	HandlerId string `json:"handler_id"`
	Event Event `json:"event"`
	Body *BodyT `json:"body"`
}

func (t *NotificationT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	handlerIdOffset := flatbuffers.UOffsetT(0)
	if t.HandlerId != "" {
		handlerIdOffset = builder.CreateString(t.HandlerId)
	}
	bodyOffset := t.Body.Pack(builder)

	NotificationStart(builder)
	NotificationAddHandlerId(builder, handlerIdOffset)
	NotificationAddEvent(builder, t.Event)
	if t.Body != nil {
		NotificationAddBodyType(builder, t.Body.Type)
	}
	NotificationAddBody(builder, bodyOffset)
	return NotificationEnd(builder)
}

func (rcv *Notification) UnPackTo(t *NotificationT) {
	t.HandlerId = string(rcv.HandlerId())
	t.Event = rcv.Event()
	bodyTable := flatbuffers.Table{}
	if rcv.Body(&bodyTable) {
		t.Body = rcv.BodyType().UnPack(bodyTable)
	}
}

func (rcv *Notification) UnPack() *NotificationT {
	if rcv == nil {
		return nil
	}
	t := &NotificationT{}
	rcv.UnPackTo(t)
	return t
}

type Notification struct {
	_tab flatbuffers.Table
}

func GetRootAsNotification(buf []byte, offset flatbuffers.UOffsetT) *Notification {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Notification{}
	x.Init(buf, n+offset)
	return x
}

func FinishNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsNotification(buf []byte, offset flatbuffers.UOffsetT) *Notification {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Notification{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Notification) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Notification) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Notification) HandlerId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Notification) Event() Event {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return Event(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *Notification) MutateEvent(n Event) bool {
	return rcv._tab.MutateByteSlot(6, byte(n))
}

func (rcv *Notification) BodyType() Body {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return Body(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *Notification) MutateBodyType(n Body) bool {
	return rcv._tab.MutateByteSlot(8, byte(n))
}

func (rcv *Notification) Body(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func NotificationStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func NotificationAddHandlerId(builder *flatbuffers.Builder, handlerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(handlerId), 0)
}
func NotificationAddEvent(builder *flatbuffers.Builder, event Event) {
	builder.PrependByteSlot(1, byte(event), 0)
}
func NotificationAddBodyType(builder *flatbuffers.Builder, bodyType Body) {
	builder.PrependByteSlot(2, byte(bodyType), 0)
}
func NotificationAddBody(builder *flatbuffers.Builder, body flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(body), 0)
}
func NotificationEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
