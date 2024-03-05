package workerchannel

import (
	"errors"
	"strings"
)

func ConvertRequestData(method string, internal *InternalData, handlerId *string) error {
	methods := strings.Split(method, ".")
	if len(methods) != 2 {
		return errors.New("method invalid")
	}
	switch methods[0] {
	case "worker":
		return nil
	case "router":
		*handlerId = internal.RouterId
	case "transport":
		*handlerId = internal.TransportId
	case "producer":
		*handlerId = internal.ProducerId
	case "consumer":
		*handlerId = internal.ConsumerId
	case "dataProducer":
		*handlerId = internal.DataProducerId
	case "dataConsumer":
		*handlerId = internal.DataConsumerId
	case "rtpObserver":
		*handlerId = internal.RtpObserverId
	default:
		return errors.New("unknown method prefix")
	}

	return nil
}
