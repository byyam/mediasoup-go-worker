// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpParameters

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RtpCodecParameters struct {
	_tab flatbuffers.Table
}

func GetRootAsRtpCodecParameters(buf []byte, offset flatbuffers.UOffsetT) *RtpCodecParameters {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RtpCodecParameters{}
	x.Init(buf, n+offset)
	return x
}

func FinishRtpCodecParametersBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRtpCodecParameters(buf []byte, offset flatbuffers.UOffsetT) *RtpCodecParameters {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RtpCodecParameters{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRtpCodecParametersBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *RtpCodecParameters) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RtpCodecParameters) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RtpCodecParameters) MimeType() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RtpCodecParameters) PayloadType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RtpCodecParameters) MutatePayloadType(n byte) bool {
	return rcv._tab.MutateByteSlot(6, n)
}

func (rcv *RtpCodecParameters) ClockRate() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RtpCodecParameters) MutateClockRate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func (rcv *RtpCodecParameters) Channels() *byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		v := rcv._tab.GetByte(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RtpCodecParameters) MutateChannels(n byte) bool {
	return rcv._tab.MutateByteSlot(10, n)
}

func (rcv *RtpCodecParameters) Parameters(obj *Parameter, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *RtpCodecParameters) ParametersLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RtpCodecParameters) RtcpFeedback(obj *RtcpFeedback, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *RtpCodecParameters) RtcpFeedbackLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RtpCodecParametersStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func RtpCodecParametersAddMimeType(builder *flatbuffers.Builder, mimeType flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(mimeType), 0)
}
func RtpCodecParametersAddPayloadType(builder *flatbuffers.Builder, payloadType byte) {
	builder.PrependByteSlot(1, payloadType, 0)
}
func RtpCodecParametersAddClockRate(builder *flatbuffers.Builder, clockRate uint32) {
	builder.PrependUint32Slot(2, clockRate, 0)
}
func RtpCodecParametersAddChannels(builder *flatbuffers.Builder, channels byte) {
	builder.PrependByte(channels)
	builder.Slot(3)
}
func RtpCodecParametersAddParameters(builder *flatbuffers.Builder, parameters flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(parameters), 0)
}
func RtpCodecParametersStartParametersVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RtpCodecParametersAddRtcpFeedback(builder *flatbuffers.Builder, rtcpFeedback flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(rtcpFeedback), 0)
}
func RtpCodecParametersStartRtcpFeedbackVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RtpCodecParametersEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
