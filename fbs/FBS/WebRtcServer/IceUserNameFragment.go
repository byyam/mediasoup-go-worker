// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WebRtcServer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type IceUserNameFragmentT struct {
	LocalIceUsernameFragment string `json:"local_ice_username_fragment"`
	WebRtcTransportId string `json:"web_rtc_transport_id"`
}

func (t *IceUserNameFragmentT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	localIceUsernameFragmentOffset := flatbuffers.UOffsetT(0)
	if t.LocalIceUsernameFragment != "" {
		localIceUsernameFragmentOffset = builder.CreateString(t.LocalIceUsernameFragment)
	}
	webRtcTransportIdOffset := flatbuffers.UOffsetT(0)
	if t.WebRtcTransportId != "" {
		webRtcTransportIdOffset = builder.CreateString(t.WebRtcTransportId)
	}
	IceUserNameFragmentStart(builder)
	IceUserNameFragmentAddLocalIceUsernameFragment(builder, localIceUsernameFragmentOffset)
	IceUserNameFragmentAddWebRtcTransportId(builder, webRtcTransportIdOffset)
	return IceUserNameFragmentEnd(builder)
}

func (rcv *IceUserNameFragment) UnPackTo(t *IceUserNameFragmentT) {
	t.LocalIceUsernameFragment = string(rcv.LocalIceUsernameFragment())
	t.WebRtcTransportId = string(rcv.WebRtcTransportId())
}

func (rcv *IceUserNameFragment) UnPack() *IceUserNameFragmentT {
	if rcv == nil {
		return nil
	}
	t := &IceUserNameFragmentT{}
	rcv.UnPackTo(t)
	return t
}

type IceUserNameFragment struct {
	_tab flatbuffers.Table
}

func GetRootAsIceUserNameFragment(buf []byte, offset flatbuffers.UOffsetT) *IceUserNameFragment {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &IceUserNameFragment{}
	x.Init(buf, n+offset)
	return x
}

func FinishIceUserNameFragmentBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsIceUserNameFragment(buf []byte, offset flatbuffers.UOffsetT) *IceUserNameFragment {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &IceUserNameFragment{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedIceUserNameFragmentBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *IceUserNameFragment) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *IceUserNameFragment) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *IceUserNameFragment) LocalIceUsernameFragment() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *IceUserNameFragment) WebRtcTransportId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func IceUserNameFragmentStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func IceUserNameFragmentAddLocalIceUsernameFragment(builder *flatbuffers.Builder, localIceUsernameFragment flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(localIceUsernameFragment), 0)
}
func IceUserNameFragmentAddWebRtcTransportId(builder *flatbuffers.Builder, webRtcTransportId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(webRtcTransportId), 0)
}
func IceUserNameFragmentEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
