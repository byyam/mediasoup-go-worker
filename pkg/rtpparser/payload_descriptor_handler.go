package rtpparser

type PayloadDescriptorHandler interface {
	GetSpatialLayer() uint8
	GetTemporalLayer() uint8
	IsKeyFrame() bool
}
