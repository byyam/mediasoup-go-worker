package rtpparser

type h264PayloadDescriptor struct {
}

func (p h264PayloadDescriptor) GetSpatialLayer() uint8 {
	return 0
}

func (p h264PayloadDescriptor) GetTemporalLayer() uint8 {
	return 0
}

func (p h264PayloadDescriptor) IsKeyFrame() bool {
	return false
}

func parseH264() PayloadDescriptorHandler {
	return &h264PayloadDescriptor{}
}

func ProcessRtpPacketH264(packet *Packet) {

}
