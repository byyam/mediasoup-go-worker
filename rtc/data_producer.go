package rtc

import (
	"github.com/kr/pretty"
	"github.com/rs/zerolog"

	FBS__DataProducer "github.com/byyam/mediasoup-go-worker/fbs/FBS/DataProducer"
	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Response "github.com/byyam/mediasoup-go-worker/fbs/FBS/Response"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/pkg/rtctime"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
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
	p.logger.Debug().Msgf("newDataProducer options:%# v", pretty.Formatter(opt))
	workerchannel.RegisterHandler(p.id, p.HandleRequest)
	return p, nil
}

func (p *DataProducer) FillJson() *FBS__DataProducer.DumpResponseT {
	dumpData := &FBS__DataProducer.DumpResponseT{
		Id:                   p.id,
		Type:                 p.opt.Type,
		SctpStreamParameters: p.opt.SctpStreamParameters,
		Label:                p.opt.Label,
		Protocol:             p.opt.Protocol,
	}
	return dumpData
}

func (p *DataProducer) FillJsonStats() *FBS__DataProducer.GetStatsResponseT {
	dumpData := &FBS__DataProducer.GetStatsResponseT{
		Timestamp:        uint64(rtctime.GetTimeMs()),
		Label:            p.opt.Label,
		Protocol:         p.opt.Protocol,
		MessagesReceived: 0,
		BytesReceived:    0,
		BufferedAmount:   0,
	}
	return dumpData
}

func (p *DataProducer) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		p.logger.Debug().Msgf("method=%s,internal=%+v,response:%s", request.Method, request.Internal, response)
	}()

	switch request.MethodType {
	case FBS__Request.MethodDATAPRODUCER_GET_STATS:
		dataDump := p.FillJsonStats()
		// set rsp
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyDataProducer_GetStatsResponse,
			Value: dataDump,
		}
		response.RspBody = rspBody

	default:
		response.Err = mserror.ErrInvalidMethod
		p.logger.Error().Msgf("unknown method:%s", request.Method)
	}

}

func (p *DataProducer) Close() {
	p.logger.Info().Msgf("DataProducer:%s closed", p.id)
	workerchannel.UnregisterHandler(p.id)
}
