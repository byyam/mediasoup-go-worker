package mediasoupdata

type WorkerSettings struct {
	/**
	 * Logging level for logs generated by the media worker subprocesses (check
	 * the Debugging documentation). Valid values are 'debug', 'warn', 'error' and
	 * 'none'. Default 'error'.
	 */
	LogLevel WorkerLogLevel `json:"logLevel,omitempty"`

	/**
	 * Log tags for debugging. Check the meaning of each available tag in the
	 * Debugging documentation.
	 */
	LogTags []WorkerLogTag `json:"logTags,omitempty"`

	/**
	 * Minimun RTC port for ICE, DTLS, RTP, etc. Default 10000.
	 */
	RtcMinPort uint16 `json:"rtcMinPort,omitempty"`

	/**
	 * Maximum RTC port for ICE, DTLS, RTP, etc. Default 59999.
	 */
	RtcMaxPort uint16 `json:"rtcMaxPort,omitempty"`

	// mux port
	RtcStaticPort uint16 `json:"rtcStaticPort,omitempty"`
	RtcListenIp   string `json:"rtcListenIp,omitempty"`

	/**
	 * Path to the DTLS public certificate file in PEM format. If unset, a
	 * certificate is dynamically created.
	 */
	DtlsCertificateFile string `json:"dtlsCertificateFile,omitempty"`

	/**
	 * Path to the DTLS certificate private key file in PEM format. If unset, a
	 * certificate is dynamically created.
	 */
	DtlsPrivateKeyFile string `json:"dtlsPrivateKeyFile,omitempty"`

	PrometheusPath string `json:"prometheusPath,omitempty"`
	PrometheusPort int    `json:"prometheusPort,omitempty"`

	PipePort int `json:"pipePort,omitempty"`
	/**
	 * Custom application data.
	 */
	AppData interface{} `json:"appData,omitempty"`

	/**
	 * Custom options.
	 */
	CustomOptions map[string]interface{}
}
