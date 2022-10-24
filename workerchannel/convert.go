package workerchannel

import (
	"errors"
	"strings"
)

func ConvertRequestData(request *RequestData) error {
	methods := strings.Split(request.Method, ".")
	if len(methods) != 2 {
		return errors.New("method invalid")
	}
	switch methods[0] {
	case "worker":
		return nil
	case "router":
		request.HandlerId = request.Internal.RouterId
	case "transport":
		request.HandlerId = request.Internal.TransportId
	case "producer":
		request.HandlerId = request.Internal.ProducerId
	case "consumer":
		request.HandlerId = request.Internal.ConsumerId
	case "dataProducer":
		request.HandlerId = request.Internal.DataProducerId
	case "dataConsumer":
		request.HandlerId = request.Internal.DataConsumerId
	case "rtpObserver":
		request.HandlerId = request.Internal.RtpObserverId
	default:
		return errors.New("unknown method prefix")
	}

	return nil
}
