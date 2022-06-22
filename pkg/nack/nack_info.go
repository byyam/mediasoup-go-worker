package nack

type nackInfo struct {
	createdAtMs uint64
	seq         uint16
	sendAtSeq   uint16
	sentAtMs    uint64
	retries     uint8
}

const (
	nackFilterSEQ = iota
	nackFilterTIME
)
