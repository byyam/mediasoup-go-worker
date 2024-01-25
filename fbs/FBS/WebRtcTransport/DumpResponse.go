// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WebRtcTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
)

type DumpResponseT struct {
	Base *FBS__Transport.DumpT `json:"base"`
	IceRole IceRole `json:"ice_role"`
	IceParameters *IceParametersT `json:"ice_parameters"`
	IceCandidates []*IceCandidateT `json:"ice_candidates"`
	IceState IceState `json:"ice_state"`
	IceSelectedTuple *FBS__Transport.TupleT `json:"ice_selected_tuple"`
	DtlsParameters *DtlsParametersT `json:"dtls_parameters"`
	DtlsState DtlsState `json:"dtls_state"`
}

func (t *DumpResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	baseOffset := t.Base.Pack(builder)
	iceParametersOffset := t.IceParameters.Pack(builder)
	iceCandidatesOffset := flatbuffers.UOffsetT(0)
	if t.IceCandidates != nil {
		iceCandidatesLength := len(t.IceCandidates)
		iceCandidatesOffsets := make([]flatbuffers.UOffsetT, iceCandidatesLength)
		for j := 0; j < iceCandidatesLength; j++ {
			iceCandidatesOffsets[j] = t.IceCandidates[j].Pack(builder)
		}
		DumpResponseStartIceCandidatesVector(builder, iceCandidatesLength)
		for j := iceCandidatesLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(iceCandidatesOffsets[j])
		}
		iceCandidatesOffset = builder.EndVector(iceCandidatesLength)
	}
	iceSelectedTupleOffset := t.IceSelectedTuple.Pack(builder)
	dtlsParametersOffset := t.DtlsParameters.Pack(builder)
	DumpResponseStart(builder)
	DumpResponseAddBase(builder, baseOffset)
	DumpResponseAddIceRole(builder, t.IceRole)
	DumpResponseAddIceParameters(builder, iceParametersOffset)
	DumpResponseAddIceCandidates(builder, iceCandidatesOffset)
	DumpResponseAddIceState(builder, t.IceState)
	DumpResponseAddIceSelectedTuple(builder, iceSelectedTupleOffset)
	DumpResponseAddDtlsParameters(builder, dtlsParametersOffset)
	DumpResponseAddDtlsState(builder, t.DtlsState)
	return DumpResponseEnd(builder)
}

func (rcv *DumpResponse) UnPackTo(t *DumpResponseT) {
	t.Base = rcv.Base(nil).UnPack()
	t.IceRole = rcv.IceRole()
	t.IceParameters = rcv.IceParameters(nil).UnPack()
	iceCandidatesLength := rcv.IceCandidatesLength()
	t.IceCandidates = make([]*IceCandidateT, iceCandidatesLength)
	for j := 0; j < iceCandidatesLength; j++ {
		x := IceCandidate{}
		rcv.IceCandidates(&x, j)
		t.IceCandidates[j] = x.UnPack()
	}
	t.IceState = rcv.IceState()
	t.IceSelectedTuple = rcv.IceSelectedTuple(nil).UnPack()
	t.DtlsParameters = rcv.DtlsParameters(nil).UnPack()
	t.DtlsState = rcv.DtlsState()
}

func (rcv *DumpResponse) UnPack() *DumpResponseT {
	if rcv == nil {
		return nil
	}
	t := &DumpResponseT{}
	rcv.UnPackTo(t)
	return t
}

type DumpResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsDumpResponse(buf []byte, offset flatbuffers.UOffsetT) *DumpResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DumpResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishDumpResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsDumpResponse(buf []byte, offset flatbuffers.UOffsetT) *DumpResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &DumpResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedDumpResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *DumpResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DumpResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *DumpResponse) Base(obj *FBS__Transport.Dump) *FBS__Transport.Dump {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__Transport.Dump)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *DumpResponse) IceRole() IceRole {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return IceRole(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *DumpResponse) MutateIceRole(n IceRole) bool {
	return rcv._tab.MutateByteSlot(6, byte(n))
}

func (rcv *DumpResponse) IceParameters(obj *IceParameters) *IceParameters {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(IceParameters)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *DumpResponse) IceCandidates(obj *IceCandidate, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *DumpResponse) IceCandidatesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) IceState() IceState {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return IceState(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *DumpResponse) MutateIceState(n IceState) bool {
	return rcv._tab.MutateByteSlot(12, byte(n))
}

func (rcv *DumpResponse) IceSelectedTuple(obj *FBS__Transport.Tuple) *FBS__Transport.Tuple {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__Transport.Tuple)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *DumpResponse) DtlsParameters(obj *DtlsParameters) *DtlsParameters {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(DtlsParameters)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *DumpResponse) DtlsState() DtlsState {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return DtlsState(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *DumpResponse) MutateDtlsState(n DtlsState) bool {
	return rcv._tab.MutateByteSlot(18, byte(n))
}

func DumpResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(8)
}
func DumpResponseAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func DumpResponseAddIceRole(builder *flatbuffers.Builder, iceRole IceRole) {
	builder.PrependByteSlot(1, byte(iceRole), 0)
}
func DumpResponseAddIceParameters(builder *flatbuffers.Builder, iceParameters flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(iceParameters), 0)
}
func DumpResponseAddIceCandidates(builder *flatbuffers.Builder, iceCandidates flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(iceCandidates), 0)
}
func DumpResponseStartIceCandidatesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddIceState(builder *flatbuffers.Builder, iceState IceState) {
	builder.PrependByteSlot(4, byte(iceState), 0)
}
func DumpResponseAddIceSelectedTuple(builder *flatbuffers.Builder, iceSelectedTuple flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(iceSelectedTuple), 0)
}
func DumpResponseAddDtlsParameters(builder *flatbuffers.Builder, dtlsParameters flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(dtlsParameters), 0)
}
func DumpResponseAddDtlsState(builder *flatbuffers.Builder, dtlsState DtlsState) {
	builder.PrependByteSlot(7, byte(dtlsState), 0)
}
func DumpResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
