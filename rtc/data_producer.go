package rtc

import (
	"encoding/json"

	"github.com/kr/pretty"
	"github.com/rs/zerolog"

	FBS__DataProducer "github.com/byyam/mediasoup-go-worker/fbs/FBS/DataProducer"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

type DataProducer struct {
	id             string
	logger         zerolog.Logger
	maxMessageSize uint32
	opt            *FBS__Transport.ProduceDataRequestT
}

func newDataProducer(id string, maxMessageSize uint32, opt *FBS__Transport.ProduceDataRequestT) (*DataProducer, error) {
	p := &DataProducer{
		id:             id,
		logger:         zerowrapper.NewScope("data-producer", id),
		maxMessageSize: maxMessageSize,
		opt:            opt,
	}
	p.logger.Info().Msgf("newDataProducer options:%# v", pretty.Formatter(opt))
	return p, nil
}

func (p *DataProducer) FillJson() json.RawMessage {
	dumpData := &FBS__DataProducer.DumpResponseT{
		Id:                   p.id,
		Type:                 p.opt.Type,
		SctpStreamParameters: p.opt.SctpStreamParameters,
		Label:                p.opt.Label,
		Protocol:             p.opt.Protocol,
	}
	data, _ := json.Marshal(&dumpData)
	p.logger.Debug().Msgf("dumpData:%+v", dumpData)
	return data
}
