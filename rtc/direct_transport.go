package rtc

import (
	"encoding/json"
	"github.com/byyam/mediasoup-go-worker/utils"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
)

type DirectTransport struct {
	ITransport
	id     string
	logger utils.Logger
}

type directTransportParam struct {
	options mediasoupdata.DirectTransportOptions
	transportParam
}

func newDirectTransport(param directTransportParam) (ITransport, error) {
	var err error
	t := &DirectTransport{
		id:     param.Id,
		logger: utils.NewLogger("direct-transport", param.Id),
	}
	t.ITransport, err = newTransport(param.transportParam)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *DirectTransport) FillJson() json.RawMessage {
	// todo
	return nil
}
