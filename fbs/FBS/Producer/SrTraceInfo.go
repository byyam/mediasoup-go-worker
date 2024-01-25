// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Producer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SrTraceInfo struct {
	_tab flatbuffers.Table
}

func GetRootAsSrTraceInfo(buf []byte, offset flatbuffers.UOffsetT) *SrTraceInfo {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SrTraceInfo{}
	x.Init(buf, n+offset)
	return x
}

func FinishSrTraceInfoBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSrTraceInfo(buf []byte, offset flatbuffers.UOffsetT) *SrTraceInfo {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SrTraceInfo{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSrTraceInfoBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SrTraceInfo) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SrTraceInfo) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SrTraceInfo) Ssrc() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SrTraceInfo) MutateSsrc(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *SrTraceInfo) NtpSec() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SrTraceInfo) MutateNtpSec(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *SrTraceInfo) NtpFrac() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SrTraceInfo) MutateNtpFrac(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func (rcv *SrTraceInfo) RtpTs() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SrTraceInfo) MutateRtpTs(n uint32) bool {
	return rcv._tab.MutateUint32Slot(10, n)
}

func (rcv *SrTraceInfo) PacketCount() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SrTraceInfo) MutatePacketCount(n uint32) bool {
	return rcv._tab.MutateUint32Slot(12, n)
}

func (rcv *SrTraceInfo) OctetCount() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *SrTraceInfo) MutateOctetCount(n uint32) bool {
	return rcv._tab.MutateUint32Slot(14, n)
}

func SrTraceInfoStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func SrTraceInfoAddSsrc(builder *flatbuffers.Builder, ssrc uint32) {
	builder.PrependUint32Slot(0, ssrc, 0)
}
func SrTraceInfoAddNtpSec(builder *flatbuffers.Builder, ntpSec uint32) {
	builder.PrependUint32Slot(1, ntpSec, 0)
}
func SrTraceInfoAddNtpFrac(builder *flatbuffers.Builder, ntpFrac uint32) {
	builder.PrependUint32Slot(2, ntpFrac, 0)
}
func SrTraceInfoAddRtpTs(builder *flatbuffers.Builder, rtpTs uint32) {
	builder.PrependUint32Slot(3, rtpTs, 0)
}
func SrTraceInfoAddPacketCount(builder *flatbuffers.Builder, packetCount uint32) {
	builder.PrependUint32Slot(4, packetCount, 0)
}
func SrTraceInfoAddOctetCount(builder *flatbuffers.Builder, octetCount uint32) {
	builder.PrependUint32Slot(5, octetCount, 0)
}
func SrTraceInfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
