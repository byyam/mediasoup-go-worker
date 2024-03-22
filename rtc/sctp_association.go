package rtc

import (
	"github.com/rs/zerolog"

	FBS__SctpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/SctpParameters"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

type SctpAssociation struct {
	options *FBS__Transport.OptionsT
	logger  zerolog.Logger
}

func newSctpAssociation(options *FBS__Transport.OptionsT) (*SctpAssociation, error) {
	t := &SctpAssociation{
		options: options,
		logger:  zerowrapper.NewScope("sctp-association"),
	}
	return t, nil
}

func (t *SctpAssociation) GetSctpAssociationParam() *FBS__SctpParameters.SctpParametersT {
	return &FBS__SctpParameters.SctpParametersT{
		Port:               5000,
		Os:                 t.options.NumSctpStreams.Os,
		Mis:                t.options.NumSctpStreams.Mis,
		MaxMessageSize:     t.options.MaxSctpMessageSize,
		IsDataChannel:      t.options.IsDataChannel,
		SctpBufferedAmount: t.options.SctpSendBufferSize,
		SendBufferSize:     t.options.SctpSendBufferSize,
	}
}
