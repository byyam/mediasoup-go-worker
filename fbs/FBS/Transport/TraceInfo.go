// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"strconv"

	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
)

type TraceInfoT struct {
	Type TraceInfo
	Value interface{}
}

func (t *TraceInfoT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	switch t.Type {
	case TraceInfoBweTraceInfo:
		return t.Value.(*FBS__Transport.BweTraceInfoT).Pack(builder)
	}
	return 0
}

func (rcv TraceInfo) UnPack(table flatbuffers.Table) *TraceInfoT {
	switch rcv {
	case TraceInfoBweTraceInfo:
		var x FBS__Transport.BweTraceInfo
		x.Init(table.Bytes, table.Pos)
		return &FBS__Transport.TraceInfoT{Type: TraceInfoBweTraceInfo, Value: x.UnPack()}
	}
	return nil
}