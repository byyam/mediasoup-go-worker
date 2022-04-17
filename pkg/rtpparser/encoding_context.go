package rtpparser

type EncodingContextParam struct {
	SpatialLayers  uint8
	TemporalLayers uint8
	Ksvc           bool
}
type EncodingContext struct {
	params               *EncodingContextParam
	targetSpatialLayer   int16
	targetTemporalLayer  int16
	currentSpatialLayer  int16
	currentTemporalLayer int16
}

func NewEncodingContext(param *EncodingContextParam) *EncodingContext {
	return &EncodingContext{
		params:               param,
		targetSpatialLayer:   -1,
		targetTemporalLayer:  -1,
		currentSpatialLayer:  -1,
		currentTemporalLayer: -1,
	}
}

func (p EncodingContext) GetSpatialLayers() uint8 {
	return p.params.SpatialLayers
}

func (p EncodingContext) GetTemporalLayers() uint8 {
	return p.params.TemporalLayers
}

func (p EncodingContext) IsKSvc() bool {
	return p.params.Ksvc
}

func (p EncodingContext) GetTargetSpatialLayer() int16 {
	return p.targetSpatialLayer
}

func (p EncodingContext) GetTargetTemporalLayer() int16 {
	return p.targetTemporalLayer
}

func (p EncodingContext) GetCurrentSpatialLayer() int16 {
	return p.currentSpatialLayer
}

func (p EncodingContext) GetCurrentTemporalLayer() int16 {
	return p.currentTemporalLayer
}

func (p *EncodingContext) SetTargetSpatialLayer(spatialLayer int16) {
	p.targetSpatialLayer = spatialLayer
}

func (p *EncodingContext) SetTargetTemporalLayer(temporalLayer int16) {
	p.targetTemporalLayer = temporalLayer
}

func (p *EncodingContext) SetCurrentSpatialLayer(spatialLayer int16) {
	p.currentSpatialLayer = spatialLayer
}

func (p *EncodingContext) SetCurrentTemporalLayer(temporalLayer int16) {
	p.currentTemporalLayer = temporalLayer
}
