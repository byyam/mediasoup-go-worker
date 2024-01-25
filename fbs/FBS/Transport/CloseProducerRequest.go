// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type CloseProducerRequestT struct {
	ProducerId string `json:"producer_id"`
}

func (t *CloseProducerRequestT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	producerIdOffset := flatbuffers.UOffsetT(0)
	if t.ProducerId != "" {
		producerIdOffset = builder.CreateString(t.ProducerId)
	}
	CloseProducerRequestStart(builder)
	CloseProducerRequestAddProducerId(builder, producerIdOffset)
	return CloseProducerRequestEnd(builder)
}

func (rcv *CloseProducerRequest) UnPackTo(t *CloseProducerRequestT) {
	t.ProducerId = string(rcv.ProducerId())
}

func (rcv *CloseProducerRequest) UnPack() *CloseProducerRequestT {
	if rcv == nil {
		return nil
	}
	t := &CloseProducerRequestT{}
	rcv.UnPackTo(t)
	return t
}

type CloseProducerRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsCloseProducerRequest(buf []byte, offset flatbuffers.UOffsetT) *CloseProducerRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &CloseProducerRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishCloseProducerRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsCloseProducerRequest(buf []byte, offset flatbuffers.UOffsetT) *CloseProducerRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &CloseProducerRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedCloseProducerRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *CloseProducerRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CloseProducerRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *CloseProducerRequest) ProducerId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func CloseProducerRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func CloseProducerRequestAddProducerId(builder *flatbuffers.Builder, producerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(producerId), 0)
}
func CloseProducerRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
