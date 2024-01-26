package message_fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/byyam/mediasoup-go-worker/fbs/FBS/Message"
)

type MessageT struct {
	Data *BodyT `json:"data"`
}

func (t *MessageT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	dataOffset := t.Data.Pack(builder)

	Message.MessageStart(builder)
	if t.Data != nil {
		Message.MessageAddDataType(builder, t.Data.Type)
	}
	Message.MessageAddData(builder, dataOffset)
	return Message.MessageEnd(builder)
}

func MessageUnPackTo(rcv *Message.Message, t *MessageT) {
	dataTable := flatbuffers.Table{}
	if rcv.Data(&dataTable) {
		t.Data = BodyUnPack(rcv.DataType(), dataTable)
	}
}

func MessageUnPack(rcv *Message.Message) *MessageT {
	if rcv == nil {
		return nil
	}
	t := &MessageT{}
	MessageUnPackTo(rcv, t)
	return t
}
