package request_fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
)

type RequestT struct {
	Id        uint32         `json:"id"`
	Method    Request.Method `json:"method"`
	HandlerId string         `json:"handler_id"`
	Body      *BodyT         `json:"body"`
}

func (t *RequestT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	handlerIdOffset := flatbuffers.UOffsetT(0)
	if t.HandlerId != "" {
		handlerIdOffset = builder.CreateString(t.HandlerId)
	}
	bodyOffset := t.Body.Pack(builder)

	Request.RequestStart(builder)
	Request.RequestAddId(builder, t.Id)
	Request.RequestAddMethod(builder, t.Method)
	Request.RequestAddHandlerId(builder, handlerIdOffset)
	if t.Body != nil {
		Request.RequestAddBodyType(builder, t.Body.Type)
	}
	Request.RequestAddBody(builder, bodyOffset)
	return Request.RequestEnd(builder)
}

func RequestUnPackTo(rcv *Request.Request, t *RequestT) {
	t.Id = rcv.Id()
	t.Method = rcv.Method()
	t.HandlerId = string(rcv.HandlerId())
	bodyTable := flatbuffers.Table{}
	if rcv.Body(&bodyTable) {
		BodyUnPack(t.Body.Type, bodyTable)
	}
}

func RequestUnPack(rcv *Request.Request) *RequestT {
	if rcv == nil {
		return nil
	}
	t := &RequestT{}
	RequestUnPackTo(rcv, t)
	return t
}
