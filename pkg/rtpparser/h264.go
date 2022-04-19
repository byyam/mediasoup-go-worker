package rtpparser

import (
	"bytes"
	"encoding/binary"
)

type h264PayloadDescriptor struct {
	isKeyFrame bool
}

func (p h264PayloadDescriptor) GetSpatialLayer() uint8 {
	return 0
}

func (p h264PayloadDescriptor) GetTemporalLayer() uint8 {
	return 0
}

func (p h264PayloadDescriptor) IsKeyFrame() bool {
	return p.isKeyFrame
}

func (p h264PayloadDescriptor) keyFrame(data []byte) bool {
	const (
		typeSTAPA       = 24
		typeSPS         = 7
		naluTypeBitmask = 0x1F
	)

	var word uint32

	payload := bytes.NewReader(data)
	if err := binary.Read(payload, binary.BigEndian, &word); err != nil {
		return false
	}

	naluType := (word >> 24) & naluTypeBitmask
	if naluType == typeSTAPA && word&naluTypeBitmask == typeSPS {
		return true
	} else if naluType == typeSPS {
		return true
	}

	return false
}

func parseH264() PayloadDescriptorHandler {
	return &h264PayloadDescriptor{}
}

func ProcessRtpPacketH264(packet *Packet) {

}
