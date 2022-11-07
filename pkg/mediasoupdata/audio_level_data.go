package mediasoupdata

// AudioLevelObserverOptions define options to create an AudioLevelObserver.
type AudioLevelObserverOptions struct {
	// MaxEntries is maximum int of entries in the 'volumes”' event. Default 1.
	MaxEntries int `json:"maxEntries"`

	// Threshold is minimum average volume (in dBvo from -127 to 0) for entries in the
	// "volumes" event.	Default -80.
	Threshold int `json:"threshold"`

	// Interval in ms for checking audio volumes. Default 1000.
	Interval int `json:"interval"`

	// AppData is custom application data.
	AppData interface{} `json:"appData,omitempty"`
}
