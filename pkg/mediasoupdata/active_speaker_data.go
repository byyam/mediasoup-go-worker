package mediasoupdata

type ActiveSpeakerObserverOptions struct {
	Interval int         `json:"interval"`
	AppData  interface{} `json:"appData,omitempty"`
}
