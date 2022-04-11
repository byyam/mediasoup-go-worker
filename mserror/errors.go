package mserror

import "errors"

var (
	ErrInvalidParam            = errors.New("invalid param")
	ErrInvalidMethod           = errors.New("invalid method")
	ErrRouterNotFound          = errors.New("router not found")
	ErrTransportNotFound       = errors.New("transport not found")
	ErrCreateWebrtcTransport   = errors.New("create webrtc-transport failed")
	ErrDuplicatedId            = errors.New("duplicated id")
	ErrProducerExist           = errors.New("producer already exist")
	ErrProducerNotFound        = errors.New("producer not found")
	ErrConsumerNotFound        = errors.New("consumer not found")
	ErrNoSRTPProtectionProfile = errors.New("DTLS Handshake completed and no SRTP Protection Profile was chosen")
)
