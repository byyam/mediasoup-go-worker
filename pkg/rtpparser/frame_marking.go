package rtpparser

type FrameMarking struct {
	tid         uint8
	base        uint8
	discardable uint8
	independent uint8
	end         uint8
	start       uint8
	lid         uint8
	tl0picidx   uint8
}

func NewFrameMarking() *FrameMarking {
	return &FrameMarking{
		tid:         3,
		base:        1,
		discardable: 1,
		independent: 1,
		end:         1,
		start:       1,
		lid:         0,
		tl0picidx:   0,
	}
}

func (p *FrameMarking) Unmarshal(buf []byte) {

}
