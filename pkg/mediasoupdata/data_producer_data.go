package mediasoupdata

// DataProducerOptions define options to create a DataProducer.
type DataProducerOptions struct {
	// Id is DataProducer id (just for Router.pipeToRouter() method).
	Id string `json:"id,omitempty"`

	Type DataProducerType `json:"type,omitempty"`

	// SctpStreamParameters define how the endpoint is sending the data.
	// Just if messages are sent over SCTP.
	SctpStreamParameters *SctpStreamParameters `json:"sctpStreamParameters,omitempty"`

	// Label can be used to distinguish this DataChannel from others.
	Label string `json:"label,omitempty"`

	// Protocol is the name of the sub-protocol used by this DataChannel.
	Protocol string `json:"protocol,omitempty"`

	// AppData is custom application data.
	AppData interface{} `json:"app_data,omitempty"`
}

// DataProducerStat define the statistic info for DataProducer.
type DataProducerStat struct {
	Type             string
	Timestamp        int64
	Label            string
	Protocol         string
	MessagesReceived int64
	BytesReceived    int64
}

// DataProducerType define DataProducer type.
type DataProducerType string

const (
	DataProducerType_Sctp   DataProducerType = "sctp"
	DataProducerType_Direct DataProducerType = "direct"
)
