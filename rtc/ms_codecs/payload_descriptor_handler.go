package ms_codecs

type PayloadDescriptorHandler interface {
	GetSpatialLayer() uint8
	GetTemporalLayer() uint8
	IsKeyFrame() bool
}
