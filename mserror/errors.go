package mserror

import "errors"

var (
	ErrInvalidParam                = errors.New("invalid param")
	ErrInvalidMethod               = errors.New("invalid method")
	ErrRouterNotFound              = errors.New("router not found")
	ErrTransportNotFound           = errors.New("transport not found")
	ErrCreateWebrtcTransport       = errors.New("create webrtc-transport failed")
	ErrCreatePipeTransport         = errors.New("create pipe-transport failed")
	ErrCreateDirectTransport       = errors.New("create direct-transport failed")
	ErrCreateAudioLevelObserver    = errors.New("create audio-level observer failed")
	ErrCreateActiveSpeakerObserver = errors.New("create active-speaker observer failed")
	ErrDuplicatedId                = errors.New("duplicated id")
	ErrProducerExist               = errors.New("producer already exist")
	ErrProducerNotFound            = errors.New("producer not found")
	ErrConsumerNotFound            = errors.New("consumer not found")
	ErrDataProducerExist           = errors.New("data producer already exist")
	ErrDataProducerNotFound        = errors.New("data producer not found")
	ErrNoSRTPProtectionProfile     = errors.New("DTLS Handshake completed and no SRTP Protection Profile was chosen")
	ErrSubTypeNotRtx               = errors.New("mimeType.subtype is not RTX")
)
