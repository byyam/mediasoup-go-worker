package mediasoupdata

const (
	MethodWorkerCreateRouter     = "worker.createRouter"
	MethodWorkerDump             = "worker.dump"
	MethodWorkerUpdateSettings   = "worker.updateSettings"
	MethodWorkerGetResourceUsage = "worker.getResourceUsage"
	MethodWorkerClose            = "worker.close"

	MethodRouterCreateWebRtcTransport       = "router.createWebRtcTransport"
	MethodRouterCreatePlainTransport        = "router.createPlainTransport"
	MethodRouterCreatePipeTransport         = "router.createPipeTransport"
	MethodRouterCreateDirectTransport       = "router.createDirectTransport"
	MethodRouterCreateActiveSpeakerObserver = "router.createActiveSpeakerObserver"
	MethodRouterCreateAudioLevelObserver    = "router.createAudioLevelObserver"
	MethodRouterDump                        = "router.dump"
	MethodRouterClose                       = "router.close"

	MethodTransportDump                  = "transport.dump"
	MethodTransportSetMaxIncomingBitrate = "transport.setMaxIncomingBitrate"
	MethodTransportSetMaxOutgoingBitrate = "transport.setMaxOutgoingBitrate"
	MethodTransportProduce               = "transport.produce"
	MethodTransportConsume               = "transport.consume"
	MethodTransportProduceData           = "transport.produceData"
	MethodTransportConsumeData           = "transport.consumeData"
	MethodTransportEnableTraceEvent      = "transport.enableTraceEvent"
	MethodTransportClose                 = "transport.close"
	MethodTransportGetStats              = "transport.getStats"
	MethodTransportConnect               = "transport.connect"
	MethodTransportRestartIce            = "transport.restartIce"

	MethodProducerDump             = "producer.dump"
	MethodProducerPause            = "producer.pause"
	MethodProducerResume           = "producer.resume"
	MethodProducerEnableTraceEvent = "producer.enableTraceEvent"
	MethodProducerClose            = "producer.close"
	MethodProducerGetStats         = "producer.getStats"

	MethodConsumerDump               = "consumer.dump"
	MethodConsumerPause              = "consumer.pause"
	MethodConsumerResume             = "consumer.resume"
	MethodConsumerEnableTraceEvent   = "consumer.enableTraceEvent"
	MethodConsumerSetPreferredLayers = "consumer.setPreferredLayers"
	MethodConsumerSetPriority        = "consumer.setPriority"
	MethodConsumerRequestKeyFrame    = "consumer.requestKeyFrame"
	MethodConsumerClose              = "consumer.close"
	MethodConsumerGetStats           = "consumer.getStats"

	MethodDataProducerClose    = "dataProducer.close"
	MethodDataProducerDump     = "dataProducer.dump"
	MethodDataProducerGetStats = "dataProducer.getStats"

	MethodDataConsumerClose                         = "dataConsumer.close"
	MethodDataConsumerDump                          = "dataConsumer.dump"
	MethodDataConsumerGetStats                      = "dataConsumer.getStats"
	MethodDataConsumerSetBufferedAmountLowThreshold = "dataConsumer.setBufferedAmountLowThreshold"
	MethodDataConsumerGetBufferedAmount             = "dataConsumer.getBufferedAmount"

	MethodRtpObserverPause          = "rtpObserver.pause"
	MethodRtpObserverResume         = "rtpObserver.resume"
	MethodRtpObserverAddProducer    = "rtpObserver.addProducer"
	MethodRtpObserverRemoveProducer = "rtpObserver.removeProducer"
	MethodRtpObserverClose          = "rtpObserver.close"
)

const (
	MethodPrefixWorker      = "worker"
	MethodPrefixRouter      = "router"
	MethodPrefixTransport   = "transport"
	MethodPrefixProducer    = "producer"
	MethodPrefixConsumer    = "consumer"
	MethodPrefixRtpObserver = "rtpObserver"
)
