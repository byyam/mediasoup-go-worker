// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Tuple struct {
	_tab flatbuffers.Table
}

func GetRootAsTuple(buf []byte, offset flatbuffers.UOffsetT) *Tuple {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Tuple{}
	x.Init(buf, n+offset)
	return x
}

func FinishTupleBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsTuple(buf []byte, offset flatbuffers.UOffsetT) *Tuple {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Tuple{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedTupleBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Tuple) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Tuple) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Tuple) LocalIp() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Tuple) LocalPort() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Tuple) MutateLocalPort(n uint16) bool {
	return rcv._tab.MutateUint16Slot(6, n)
}

func (rcv *Tuple) RemoteIp() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Tuple) RemotePort() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Tuple) MutateRemotePort(n uint16) bool {
	return rcv._tab.MutateUint16Slot(10, n)
}

func (rcv *Tuple) Protocol() Protocol {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return Protocol(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 1
}

func (rcv *Tuple) MutateProtocol(n Protocol) bool {
	return rcv._tab.MutateByteSlot(12, byte(n))
}

func TupleStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func TupleAddLocalIp(builder *flatbuffers.Builder, localIp flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(localIp), 0)
}
func TupleAddLocalPort(builder *flatbuffers.Builder, localPort uint16) {
	builder.PrependUint16Slot(1, localPort, 0)
}
func TupleAddRemoteIp(builder *flatbuffers.Builder, remoteIp flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(remoteIp), 0)
}
func TupleAddRemotePort(builder *flatbuffers.Builder, remotePort uint16) {
	builder.PrependUint16Slot(3, remotePort, 0)
}
func TupleAddProtocol(builder *flatbuffers.Builder, protocol Protocol) {
	builder.PrependByteSlot(4, byte(protocol), 1)
}
func TupleEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
