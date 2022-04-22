package rtc

const (
	// 17: 16 bit mask + the initial sequence number.
	MaxRequestedPackets = 17
	// Don't retransmit packets older than this (ms).
	MaxRetransmissionDelay = 2000
	DefaultRtt             = 100
)

type Retransmission struct {
	bufferSize     int
	bufferStartIdx uint16
	buffer         []*StorageItem
	storage        []StorageItem
	container      []*StorageItem
}

func NewRetransmission(bufferSize int) *Retransmission {
	r := &Retransmission{
		bufferSize:     bufferSize,
		bufferStartIdx: 0,
		buffer:         nil,
		storage:        nil,
		container:      make([]*StorageItem, MaxRequestedPackets+1),
	}
	if bufferSize > 0 { // seq 16 bits
		r.buffer = make([]*StorageItem, 65536)
	} else {
		r.buffer = nil
	}
	r.storage = make([]StorageItem, bufferSize)
	return r
}
