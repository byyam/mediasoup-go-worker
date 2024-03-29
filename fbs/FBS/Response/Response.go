// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Response

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ResponseT struct {
	Id uint32 `json:"id"`
	Accepted bool `json:"accepted"`
	Body *BodyT `json:"body"`
	Error string `json:"error"`
	Reason string `json:"reason"`
}

func (t *ResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	bodyOffset := t.Body.Pack(builder)

	errorOffset := flatbuffers.UOffsetT(0)
	if t.Error != "" {
		errorOffset = builder.CreateString(t.Error)
	}
	reasonOffset := flatbuffers.UOffsetT(0)
	if t.Reason != "" {
		reasonOffset = builder.CreateString(t.Reason)
	}
	ResponseStart(builder)
	ResponseAddId(builder, t.Id)
	ResponseAddAccepted(builder, t.Accepted)
	if t.Body != nil {
		ResponseAddBodyType(builder, t.Body.Type)
	}
	ResponseAddBody(builder, bodyOffset)
	ResponseAddError(builder, errorOffset)
	ResponseAddReason(builder, reasonOffset)
	return ResponseEnd(builder)
}

func (rcv *Response) UnPackTo(t *ResponseT) {
	t.Id = rcv.Id()
	t.Accepted = rcv.Accepted()
	bodyTable := flatbuffers.Table{}
	if rcv.Body(&bodyTable) {
		t.Body = rcv.BodyType().UnPack(bodyTable)
	}
	t.Error = string(rcv.Error())
	t.Reason = string(rcv.Reason())
}

func (rcv *Response) UnPack() *ResponseT {
	if rcv == nil {
		return nil
	}
	t := &ResponseT{}
	rcv.UnPackTo(t)
	return t
}

type Response struct {
	_tab flatbuffers.Table
}

func GetRootAsResponse(buf []byte, offset flatbuffers.UOffsetT) *Response {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Response{}
	x.Init(buf, n+offset)
	return x
}

func FinishResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsResponse(buf []byte, offset flatbuffers.UOffsetT) *Response {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Response{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Response) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Response) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Response) Id() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Response) MutateId(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *Response) Accepted() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Response) MutateAccepted(n bool) bool {
	return rcv._tab.MutateBoolSlot(6, n)
}

func (rcv *Response) BodyType() Body {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return Body(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *Response) MutateBodyType(n Body) bool {
	return rcv._tab.MutateByteSlot(8, byte(n))
}

func (rcv *Response) Body(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func (rcv *Response) Error() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Response) Reason() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func ResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func ResponseAddId(builder *flatbuffers.Builder, id uint32) {
	builder.PrependUint32Slot(0, id, 0)
}
func ResponseAddAccepted(builder *flatbuffers.Builder, accepted bool) {
	builder.PrependBoolSlot(1, accepted, false)
}
func ResponseAddBodyType(builder *flatbuffers.Builder, bodyType Body) {
	builder.PrependByteSlot(2, byte(bodyType), 0)
}
func ResponseAddBody(builder *flatbuffers.Builder, body flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(body), 0)
}
func ResponseAddError(builder *flatbuffers.Builder, error flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(error), 0)
}
func ResponseAddReason(builder *flatbuffers.Builder, reason flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(reason), 0)
}
func ResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
