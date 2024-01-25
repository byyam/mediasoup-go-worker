// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WebRtcTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
)

type IceCandidateT struct {
	Foundation string `json:"foundation"`
	Priority uint32 `json:"priority"`
	Ip string `json:"ip"`
	Protocol FBS__Transport.Protocol `json:"protocol"`
	Port uint16 `json:"port"`
	Type IceCandidateType `json:"type"`
	TcpType *IceCandidateTcpType `json:"tcp_type"`
}

func (t *IceCandidateT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	foundationOffset := flatbuffers.UOffsetT(0)
	if t.Foundation != "" {
		foundationOffset = builder.CreateString(t.Foundation)
	}
	ipOffset := flatbuffers.UOffsetT(0)
	if t.Ip != "" {
		ipOffset = builder.CreateString(t.Ip)
	}
	IceCandidateStart(builder)
	IceCandidateAddFoundation(builder, foundationOffset)
	IceCandidateAddPriority(builder, t.Priority)
	IceCandidateAddIp(builder, ipOffset)
	IceCandidateAddProtocol(builder, t.Protocol)
	IceCandidateAddPort(builder, t.Port)
	IceCandidateAddType(builder, t.Type)
	if t.TcpType != nil {
		IceCandidateAddTcpType(builder, *t.TcpType)
	}
	return IceCandidateEnd(builder)
}

func (rcv *IceCandidate) UnPackTo(t *IceCandidateT) {
	t.Foundation = string(rcv.Foundation())
	t.Priority = rcv.Priority()
	t.Ip = string(rcv.Ip())
	t.Protocol = rcv.Protocol()
	t.Port = rcv.Port()
	t.Type = rcv.Type()
	t.TcpType = rcv.TcpType()
}

func (rcv *IceCandidate) UnPack() *IceCandidateT {
	if rcv == nil {
		return nil
	}
	t := &IceCandidateT{}
	rcv.UnPackTo(t)
	return t
}

type IceCandidate struct {
	_tab flatbuffers.Table
}

func GetRootAsIceCandidate(buf []byte, offset flatbuffers.UOffsetT) *IceCandidate {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &IceCandidate{}
	x.Init(buf, n+offset)
	return x
}

func FinishIceCandidateBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsIceCandidate(buf []byte, offset flatbuffers.UOffsetT) *IceCandidate {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &IceCandidate{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedIceCandidateBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *IceCandidate) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *IceCandidate) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *IceCandidate) Foundation() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *IceCandidate) Priority() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IceCandidate) MutatePriority(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *IceCandidate) Ip() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *IceCandidate) Protocol() FBS__Transport.Protocol {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return FBS__Transport.Protocol(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 1
}

func (rcv *IceCandidate) MutateProtocol(n FBS__Transport.Protocol) bool {
	return rcv._tab.MutateByteSlot(10, byte(n))
}

func (rcv *IceCandidate) Port() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IceCandidate) MutatePort(n uint16) bool {
	return rcv._tab.MutateUint16Slot(12, n)
}

func (rcv *IceCandidate) Type() IceCandidateType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return IceCandidateType(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *IceCandidate) MutateType(n IceCandidateType) bool {
	return rcv._tab.MutateByteSlot(14, byte(n))
}

func (rcv *IceCandidate) TcpType() *IceCandidateTcpType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		v := IceCandidateTcpType(rcv._tab.GetByte(o + rcv._tab.Pos))
		return &v
	}
	return nil
}

func (rcv *IceCandidate) MutateTcpType(n IceCandidateTcpType) bool {
	return rcv._tab.MutateByteSlot(16, byte(n))
}

func IceCandidateStart(builder *flatbuffers.Builder) {
	builder.StartObject(7)
}
func IceCandidateAddFoundation(builder *flatbuffers.Builder, foundation flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(foundation), 0)
}
func IceCandidateAddPriority(builder *flatbuffers.Builder, priority uint32) {
	builder.PrependUint32Slot(1, priority, 0)
}
func IceCandidateAddIp(builder *flatbuffers.Builder, ip flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(ip), 0)
}
func IceCandidateAddProtocol(builder *flatbuffers.Builder, protocol FBS__Transport.Protocol) {
	builder.PrependByteSlot(3, byte(protocol), 1)
}
func IceCandidateAddPort(builder *flatbuffers.Builder, port uint16) {
	builder.PrependUint16Slot(4, port, 0)
}
func IceCandidateAddType(builder *flatbuffers.Builder, type_ IceCandidateType) {
	builder.PrependByteSlot(5, byte(type_), 0)
}
func IceCandidateAddTcpType(builder *flatbuffers.Builder, tcpType IceCandidateTcpType) {
	builder.PrependByte(byte(tcpType))
	builder.Slot(6)
}
func IceCandidateEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
