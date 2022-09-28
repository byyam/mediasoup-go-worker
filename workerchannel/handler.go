package workerchannel

import (
	"errors"
	"sync"

	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

var (
	handlerOnce sync.Once

	channelHandlers *cHandlers
)

type Handler func(request RequestData, response *ResponseData)

type cHandlers struct {
	logger   zerolog.Logger
	handlers map[string]Handler
	sync.RWMutex
}

func InitChannelHandlers() {
	handlerOnce.Do(func() {
		channelHandlers = &cHandlers{
			handlers: make(map[string]Handler),
			logger:   zerowrapper.NewScope("cHandlers"),
		}
	})
}

func RegisterHandler(id string, h Handler) {
	if h == nil || id == "" {
		channelHandlers.logger.Warn().Msg("invalid register")
		return
	}
	channelHandlers.registerHandler(id, h)
}

func UnregisterHandler(id string) {
	if id == "" {
		channelHandlers.logger.Warn().Msg("invalid unregister")
		return
	}
	channelHandlers.unregisterHandler(id)
}

func GetChannelRequestHandler(id string) (Handler, error) {
	if id == "" {
		channelHandlers.logger.Warn().Msg("invalid get")
		return nil, errors.New("invalid id")
	}
	return channelHandlers.getHandler(id)
}

func (c *cHandlers) registerHandler(id string, h Handler) {
	c.Lock()
	defer c.Unlock()

	c.handlers[id] = h
	c.logger.Debug().Str("id", id).Msg("registerHandler")
}

func (c *cHandlers) unregisterHandler(id string) {
	c.Lock()
	defer c.Unlock()

	delete(c.handlers, id)
	c.logger.Debug().Str("id", id).Msg("unregisterHandler")
}

func (c *cHandlers) getHandler(id string) (Handler, error) {
	c.RLock()
	defer c.RUnlock()

	h, ok := c.handlers[id]
	if !ok {
		return nil, errors.New("handler not found")
	}
	return h, nil
}
