package rtc

import (
	"github.com/rs/zerolog"

	FBS__SctpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/SctpParameters"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

type SctpAssociation struct {
	options mediasoupdata.SctpOptions
	logger  zerolog.Logger
}

func newSctpAssociation(options mediasoupdata.SctpOptions) (*SctpAssociation, error) {
	t := &SctpAssociation{
		options: options,
		logger:  zerowrapper.NewScope("sctp-association"),
	}
	return t, nil
}

func (t *SctpAssociation) GetSctpAssociationParam() *FBS__SctpParameters.SctpParametersT {
	return &FBS__SctpParameters.SctpParametersT{
		Port:               5000,
		Os:                 t.options.NumSctpStreams.OS,
		Mis:                t.options.NumSctpStreams.MIS,
		MaxMessageSize:     t.options.MaxSctpMessageSize,
		IsDataChannel:      t.options.IsDataChannel,
		SctpBufferedAmount: uint32(t.options.SctpSendBufferSize),
		SendBufferSize:     uint32(t.options.SctpSendBufferSize),
	}
}
