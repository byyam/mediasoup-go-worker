package common

import "errors"

var (
	ErrRouterNotFound        = errors.New("router not found")
	ErrTransportNotFound     = errors.New("transport not found")
	ErrCreateWebrtcTransport = errors.New("create webrtc-transport failed")
)
