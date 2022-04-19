package rtpparser

type FrameMarking struct {
	tid         uint8 // 3 bit
	base        uint8
	discardable uint8
	independent uint8
	end         uint8
	start       uint8
	lid         uint8
	tl0picidx   uint8
}

func Unmarshal(buf []byte) *FrameMarking {
	return &FrameMarking{}
}
