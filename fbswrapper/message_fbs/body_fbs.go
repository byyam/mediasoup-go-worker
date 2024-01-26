package message_fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/byyam/mediasoup-go-worker/fbs/FBS/Message"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	"github.com/byyam/mediasoup-go-worker/fbswrapper/request_fbs"
)

type BodyT struct {
	Type  Message.Body
	Value interface{}
}

func (t *BodyT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	//switch t.Type {
	//case Message.BodyRequest:
	//	return t.Value.(*FBS__Request.RequestT).Pack(builder)
	//case Message.BodyResponse:
	//	return t.Value.(*FBS__Response.ResponseT).Pack(builder)
	//case Message.BodyNotification:
	//	return t.Value.(*FBS__Notification.NotificationT).Pack(builder)
	//case Message.BodyLog:
	//	return t.Value.(*FBS__Log.LogT).Pack(builder)
	//}
	return 0
}

func BodyUnPack(rcv Message.Body, table flatbuffers.Table) *BodyT {
	switch rcv {
	case Message.BodyRequest:
		var x FBS__Request.Request
		x.Init(table.Bytes, table.Pos)
		return &BodyT{Type: Message.BodyRequest, Value: request_fbs.RequestUnPack(&x)}
		//case Message.BodyResponse:
		//	var x FBS__Response.Response
		//	x.Init(table.Bytes, table.Pos)
		//	return &BodyT{Type: Message.BodyResponse, Value: x.UnPack()}
		//case Message.BodyNotification:
		//	var x FBS__Notification.Notification
		//	x.Init(table.Bytes, table.Pos)
		//	return &BodyT{Type: Message.BodyNotification, Value: x.UnPack()}
		//case Message.BodyLog:
		//	var x FBS__Log.Log
		//	x.Init(table.Bytes, table.Pos)
		//	return &BodyT{Type: Message.BodyLog, Value: x.UnPack()}
	}
	return nil
}
