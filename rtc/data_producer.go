package rtc

import (
	"encoding/json"

	"github.com/kr/pretty"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

type DataProducer struct {
	id             string
	logger         zerolog.Logger
	maxMessageSize uint32
	options        mediasoupdata.DataProducerOptions
}

func newDataProducer(id string, maxMessageSize uint32, opt mediasoupdata.DataProducerOptions) (*DataProducer, error) {
	p := &DataProducer{
		id:             id,
		logger:         zerowrapper.NewScope("data-producer", id),
		maxMessageSize: maxMessageSize,
		options:        opt,
	}
	p.logger.Info().Msgf("newDataProducer options:%# v", pretty.Formatter(opt))
	return p, nil
}

func (p *DataProducer) FillJson() json.RawMessage {
	dumpData := mediasoupdata.DataProducerDump{
		Id:                   p.id,
		Type:                 p.options.Type,
		SctpStreamParameters: p.options.SctpStreamParameters,
		Label:                p.options.Label,
		Protocol:             p.options.Protocol,
	}
	data, _ := json.Marshal(&dumpData)
	p.logger.Debug().Msgf("dumpData:%+v", dumpData)
	return data
}
