package nack

import (
	"github.com/pion/rtp"
)

type ParamGenerator struct {
	SendNackDelayMs                 uint32
	OnNackGeneratorNackRequired     func(seqNum []uint16)
	OnNackGeneratorKeyFrameRequired func()
}

type Generator struct {
	SendNackDelayMs uint32
	rtt             uint32

	// members
	nackList      map[uint16]nackInfo
	keyFrameList  []uint16
	recoveredList []uint16
}

func NewNackGenerator(param ParamGenerator) *Generator {
	return &Generator{
		SendNackDelayMs: param.SendNackDelayMs,
	}
}

func (p *Generator) ReceivePacket(packet *rtp.Packet, isKeyFrame, isRecovered bool) bool {

	return false
}

func (p *Generator) UpdateRtt(rtt uint32) {
	p.rtt = rtt
}

func (p *Generator) Reset() {
	// todo
}
