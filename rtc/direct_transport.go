package rtc

import (
	"encoding/json"

	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

type DirectTransport struct {
	ITransport
	id     string
	logger zerolog.Logger
}

type directTransportParam struct {
	options mediasoupdata.DirectTransportOptions
	transportParam
}

func newDirectTransport(param directTransportParam) (ITransport, error) {
	var err error
	t := &DirectTransport{
		id:     param.Id,
		logger: zerowrapper.NewScope("direct-transport", param.Id),
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
