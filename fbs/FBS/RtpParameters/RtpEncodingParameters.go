// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpParameters

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RtpEncodingParametersT struct {
	Ssrc *uint32 `json:"ssrc"`
	Rid string `json:"rid"`
	CodecPayloadType *byte `json:"codec_payload_type"`
	Rtx *RtxT `json:"rtx"`
	Dtx bool `json:"dtx"`
	ScalabilityMode string `json:"scalability_mode"`
	MaxBitrate *uint32 `json:"max_bitrate"`
}

func (t *RtpEncodingParametersT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	ridOffset := flatbuffers.UOffsetT(0)
	if t.Rid != "" {
		ridOffset = builder.CreateString(t.Rid)
	}
	rtxOffset := t.Rtx.Pack(builder)
	scalabilityModeOffset := flatbuffers.UOffsetT(0)
	if t.ScalabilityMode != "" {
		scalabilityModeOffset = builder.CreateString(t.ScalabilityMode)
	}
	RtpEncodingParametersStart(builder)
	if t.Ssrc != nil {
		RtpEncodingParametersAddSsrc(builder, *t.Ssrc)
	}
	RtpEncodingParametersAddRid(builder, ridOffset)
	if t.CodecPayloadType != nil {
		RtpEncodingParametersAddCodecPayloadType(builder, *t.CodecPayloadType)
	}
	RtpEncodingParametersAddRtx(builder, rtxOffset)
	RtpEncodingParametersAddDtx(builder, t.Dtx)
	RtpEncodingParametersAddScalabilityMode(builder, scalabilityModeOffset)
	if t.MaxBitrate != nil {
		RtpEncodingParametersAddMaxBitrate(builder, *t.MaxBitrate)
	}
	return RtpEncodingParametersEnd(builder)
}

func (rcv *RtpEncodingParameters) UnPackTo(t *RtpEncodingParametersT) {
	t.Ssrc = rcv.Ssrc()
	t.Rid = string(rcv.Rid())
	t.CodecPayloadType = rcv.CodecPayloadType()
	t.Rtx = rcv.Rtx(nil).UnPack()
	t.Dtx = rcv.Dtx()
	t.ScalabilityMode = string(rcv.ScalabilityMode())
	t.MaxBitrate = rcv.MaxBitrate()
}

func (rcv *RtpEncodingParameters) UnPack() *RtpEncodingParametersT {
	if rcv == nil {
		return nil
	}
	t := &RtpEncodingParametersT{}
	rcv.UnPackTo(t)
	return t
}

type RtpEncodingParameters struct {
	_tab flatbuffers.Table
}

func GetRootAsRtpEncodingParameters(buf []byte, offset flatbuffers.UOffsetT) *RtpEncodingParameters {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RtpEncodingParameters{}
	x.Init(buf, n+offset)
	return x
}

func FinishRtpEncodingParametersBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRtpEncodingParameters(buf []byte, offset flatbuffers.UOffsetT) *RtpEncodingParameters {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RtpEncodingParameters{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRtpEncodingParametersBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *RtpEncodingParameters) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RtpEncodingParameters) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RtpEncodingParameters) Ssrc() *uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		v := rcv._tab.GetUint32(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RtpEncodingParameters) MutateSsrc(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *RtpEncodingParameters) Rid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RtpEncodingParameters) CodecPayloadType() *byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		v := rcv._tab.GetByte(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RtpEncodingParameters) MutateCodecPayloadType(n byte) bool {
	return rcv._tab.MutateByteSlot(8, n)
}

func (rcv *RtpEncodingParameters) Rtx(obj *Rtx) *Rtx {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Rtx)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *RtpEncodingParameters) Dtx() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *RtpEncodingParameters) MutateDtx(n bool) bool {
	return rcv._tab.MutateBoolSlot(12, n)
}

func (rcv *RtpEncodingParameters) ScalabilityMode() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RtpEncodingParameters) MaxBitrate() *uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		v := rcv._tab.GetUint32(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RtpEncodingParameters) MutateMaxBitrate(n uint32) bool {
	return rcv._tab.MutateUint32Slot(16, n)
}

func RtpEncodingParametersStart(builder *flatbuffers.Builder) {
	builder.StartObject(7)
}
func RtpEncodingParametersAddSsrc(builder *flatbuffers.Builder, ssrc uint32) {
	builder.PrependUint32(ssrc)
	builder.Slot(0)
}
func RtpEncodingParametersAddRid(builder *flatbuffers.Builder, rid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(rid), 0)
}
func RtpEncodingParametersAddCodecPayloadType(builder *flatbuffers.Builder, codecPayloadType byte) {
	builder.PrependByte(codecPayloadType)
	builder.Slot(2)
}
func RtpEncodingParametersAddRtx(builder *flatbuffers.Builder, rtx flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(rtx), 0)
}
func RtpEncodingParametersAddDtx(builder *flatbuffers.Builder, dtx bool) {
	builder.PrependBoolSlot(4, dtx, false)
}
func RtpEncodingParametersAddScalabilityMode(builder *flatbuffers.Builder, scalabilityMode flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(scalabilityMode), 0)
}
func RtpEncodingParametersAddMaxBitrate(builder *flatbuffers.Builder, maxBitrate uint32) {
	builder.PrependUint32(maxBitrate)
	builder.Slot(6)
}
func RtpEncodingParametersEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
