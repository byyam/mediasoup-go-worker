// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RestartIceResponseT struct {
	UsernameFragment string `json:"username_fragment"`
	Password string `json:"password"`
	IceLite bool `json:"ice_lite"`
}

func (t *RestartIceResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	usernameFragmentOffset := flatbuffers.UOffsetT(0)
	if t.UsernameFragment != "" {
		usernameFragmentOffset = builder.CreateString(t.UsernameFragment)
	}
	passwordOffset := flatbuffers.UOffsetT(0)
	if t.Password != "" {
		passwordOffset = builder.CreateString(t.Password)
	}
	RestartIceResponseStart(builder)
	RestartIceResponseAddUsernameFragment(builder, usernameFragmentOffset)
	RestartIceResponseAddPassword(builder, passwordOffset)
	RestartIceResponseAddIceLite(builder, t.IceLite)
	return RestartIceResponseEnd(builder)
}

func (rcv *RestartIceResponse) UnPackTo(t *RestartIceResponseT) {
	t.UsernameFragment = string(rcv.UsernameFragment())
	t.Password = string(rcv.Password())
	t.IceLite = rcv.IceLite()
}

func (rcv *RestartIceResponse) UnPack() *RestartIceResponseT {
	if rcv == nil {
		return nil
	}
	t := &RestartIceResponseT{}
	rcv.UnPackTo(t)
	return t
}

type RestartIceResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsRestartIceResponse(buf []byte, offset flatbuffers.UOffsetT) *RestartIceResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RestartIceResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishRestartIceResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRestartIceResponse(buf []byte, offset flatbuffers.UOffsetT) *RestartIceResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RestartIceResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRestartIceResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *RestartIceResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RestartIceResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RestartIceResponse) UsernameFragment() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RestartIceResponse) Password() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RestartIceResponse) IceLite() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *RestartIceResponse) MutateIceLite(n bool) bool {
	return rcv._tab.MutateBoolSlot(8, n)
}

func RestartIceResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func RestartIceResponseAddUsernameFragment(builder *flatbuffers.Builder, usernameFragment flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(usernameFragment), 0)
}
func RestartIceResponseAddPassword(builder *flatbuffers.Builder, password flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(password), 0)
}
func RestartIceResponseAddIceLite(builder *flatbuffers.Builder, iceLite bool) {
	builder.PrependBoolSlot(2, iceLite, false)
}
func RestartIceResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
