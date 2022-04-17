package rtpparser

import "github.com/pion/rtp"

type Packet struct {
	*rtp.Packet
	packetLen int

	midExtensionId               uint8
	ridExtensionId               uint8
	rridExtensionId              uint8
	absSendTimeExtensionId       uint8
	transportWideCc01ExtensionId uint8
	frameMarking07ExtensionId    uint8
	frameMarkingExtensionId      uint8
	ssrcAudioLevelExtensionId    uint8
	videoOrientationExtensionId  uint8
}

func Parse(buf []byte) (*Packet, error) {
	packet := &rtp.Packet{}
	if err := packet.Unmarshal(buf); err != nil {
		return nil, err
	}
	p := &Packet{
		Packet:    packet,
		packetLen: len(buf),
	}
	return p, nil
}

func (p Packet) GetLen() int {
	return p.packetLen
}

func (p *Packet) SetMidExtensionId(id uint8) {
	p.midExtensionId = id
}

func (p *Packet) SetRidExtensionId(id uint8) {
	p.ridExtensionId = id
}

func (p *Packet) SetRepairedRidExtensionId(id uint8) {
	p.rridExtensionId = id
}

func (p *Packet) SetAbsSendTimeExtensionId(id uint8) {
	p.absSendTimeExtensionId = id
}

func (p *Packet) SetTransportWideCc01ExtensionId(id uint8) {
	p.transportWideCc01ExtensionId = id
}

func (p *Packet) SetFrameMarking07ExtensionId(id uint8) {
	p.frameMarking07ExtensionId = id
}

func (p *Packet) SetFrameMarkingExtensionId(id uint8) {
	p.frameMarkingExtensionId = id
}

func (p *Packet) SetSsrcAudioLevelExtensionId(id uint8) {
	p.ssrcAudioLevelExtensionId = id
}

func (p *Packet) SetVideoOrientationExtensionId(id uint8) {
	p.videoOrientationExtensionId = id
}
