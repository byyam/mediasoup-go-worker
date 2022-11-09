package rtpparser

import "github.com/pion/rtp"

type Packet struct {
	*rtp.Packet
	packetLen                int
	payloadDescriptorHandler PayloadDescriptorHandler

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

func (p *Packet) GetLen() int {
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

func (p *Packet) SetPayloadDescriptorHandler(payloadDescriptorHandler PayloadDescriptorHandler) {
	p.payloadDescriptorHandler = payloadDescriptorHandler
}

func (p *Packet) GetSpatialLayer() uint8 {
	if p.payloadDescriptorHandler == nil {
		return 0
	}
	return p.payloadDescriptorHandler.GetSpatialLayer()
}

func (p *Packet) GetTemporalLayer() uint8 {
	if p.payloadDescriptorHandler == nil {
		return 0
	}
	return p.payloadDescriptorHandler.GetTemporalLayer()
}

func (p *Packet) IsKeyFrame() bool {
	if p.payloadDescriptorHandler == nil {
		return false
	}
	return p.payloadDescriptorHandler.IsKeyFrame()
}

func (p *Packet) ReadFrameMarking(frameMarking *FrameMarking, length *uint8) bool {
	extenValue := p.GetExtension(p.frameMarkingExtensionId)
	// NOTE: Remove this once framemarking draft becomes RFC.
	if extenValue == nil {
		extenValue = p.GetExtension(p.frameMarking07ExtensionId)
	}
	extenLen := len(extenValue)
	if extenValue == nil || extenLen > 3 {
		return false
	}
	frameMarking = Unmarshal(extenValue)
	*length = uint8(extenLen)
	return true
}

func (p *Packet) GetMid() string {
	extenValue := p.GetExtension(p.midExtensionId)
	return string(extenValue)
}

func (p *Packet) GetRid() string {
	extenValue := p.GetExtension(p.ridExtensionId)
	return string(extenValue)
}

func (p *Packet) UpdateMid(mid string) error {
	return p.SetExtension(p.midExtensionId, []byte(mid))
}
