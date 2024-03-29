package rtc

import (
	"sync"

	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

type KeyFrameRequestManager struct {
	keyFrameRequestDelay          uint32
	mapSsrcKeyFrameRequestDelayer sync.Map
	logger                        zerolog.Logger
	onKeyFrameNeededHandler       func(ssrc uint32)
}

type KeyFrameRequestManagerParam struct {
	keyFrameRequestDelay uint32
	onKeyFrameNeeded     func(ssrc uint32)
}

func NewKeyFrameRequestManager(param *KeyFrameRequestManagerParam) *KeyFrameRequestManager {
	return &KeyFrameRequestManager{
		keyFrameRequestDelay:    param.keyFrameRequestDelay,
		onKeyFrameNeededHandler: param.onKeyFrameNeeded,
		logger:                  zerowrapper.NewScope("KeyFrameRequestManager"),
	}
}

func (p *KeyFrameRequestManager) ForceKeyFrameNeeded(ssrc uint32) {
	if p.keyFrameRequestDelay > 0 {

	}
}

func (p *KeyFrameRequestManager) KeyFrameNeeded(ssrc uint32) {
	if p.keyFrameRequestDelay > 0 {
		v, ok := p.mapSsrcKeyFrameRequestDelayer.Load(ssrc)

		if ok { // There is a delayer for the given ssrc, so enable it and return.
			keyFrameRequestDelayer := v.(*KeyFrameRequestDelayer)
			keyFrameRequestDelayer.SetKeyFrameRequested(true)
			return
		} else { // Otherwise create a delayer (not yet enabled) and continue.
			p.logger.Debug().Msg("creating a delayer for the given ssrc")
		}
		// todo
	}
	p.onKeyFrameNeededHandler(ssrc)
}

type KeyFrameRequestDelayer struct {
}

func (p *KeyFrameRequestDelayer) SetKeyFrameRequested(v bool) {

}
