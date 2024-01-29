// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Producer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ScoreNotificationT struct {
	Scores []*ScoreT `json:"scores"`
}

func (t *ScoreNotificationT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	scoresOffset := flatbuffers.UOffsetT(0)
	if t.Scores != nil {
		scoresLength := len(t.Scores)
		scoresOffsets := make([]flatbuffers.UOffsetT, scoresLength)
		for j := 0; j < scoresLength; j++ {
			scoresOffsets[j] = t.Scores[j].Pack(builder)
		}
		ScoreNotificationStartScoresVector(builder, scoresLength)
		for j := scoresLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(scoresOffsets[j])
		}
		scoresOffset = builder.EndVector(scoresLength)
	}
	ScoreNotificationStart(builder)
	ScoreNotificationAddScores(builder, scoresOffset)
	return ScoreNotificationEnd(builder)
}

func (rcv *ScoreNotification) UnPackTo(t *ScoreNotificationT) {
	scoresLength := rcv.ScoresLength()
	t.Scores = make([]*ScoreT, scoresLength)
	for j := 0; j < scoresLength; j++ {
		x := Score{}
		rcv.Scores(&x, j)
		t.Scores[j] = x.UnPack()
	}
}

func (rcv *ScoreNotification) UnPack() *ScoreNotificationT {
	if rcv == nil {
		return nil
	}
	t := &ScoreNotificationT{}
	rcv.UnPackTo(t)
	return t
}

type ScoreNotification struct {
	_tab flatbuffers.Table
}

func GetRootAsScoreNotification(buf []byte, offset flatbuffers.UOffsetT) *ScoreNotification {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ScoreNotification{}
	x.Init(buf, n+offset)
	return x
}

func FinishScoreNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsScoreNotification(buf []byte, offset flatbuffers.UOffsetT) *ScoreNotification {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ScoreNotification{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedScoreNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ScoreNotification) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ScoreNotification) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ScoreNotification) Scores(obj *Score, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *ScoreNotification) ScoresLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ScoreNotificationStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ScoreNotificationAddScores(builder *flatbuffers.Builder, scores flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(scores), 0)
}
func ScoreNotificationStartScoresVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ScoreNotificationEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
