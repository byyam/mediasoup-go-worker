package rtc

type RtpStreamRecv struct {
	ssrc    uint32
	rtxSsrc uint32
}

func (r *RtpStreamRecv) GetSsrc() uint32 {
	return r.ssrc
}

func (r *RtpStreamRecv) GetRtxSsrc() uint32 {
	return r.rtxSsrc
}
