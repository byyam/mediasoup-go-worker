package rtc

import (
	"github.com/rs/zerolog"

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

func (t *SctpAssociation) GetSctpAssociationParam() mediasoupdata.SctpParameters {
	return mediasoupdata.SctpParameters{
		Port:               5000,
		OS:                 t.options.NumSctpStreams.OS,
		MIS:                t.options.NumSctpStreams.MIS,
		MaxMessageSize:     t.options.MaxSctpMessageSize,
		IsDataChannel:      t.options.IsDataChannel,
		SctpBufferedAmount: t.options.SctpSendBufferSize,
		SendBufferSize:     t.options.SctpSendBufferSize,
	}
}
