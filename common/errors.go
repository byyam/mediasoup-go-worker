package common

import "errors"

var (
	ErrInvalidParam          = errors.New("invalid param")
	ErrRouterNotFound        = errors.New("router not found")
	ErrTransportNotFound     = errors.New("transport not found")
	ErrCreateWebrtcTransport = errors.New("create webrtc-transport failed")
)
