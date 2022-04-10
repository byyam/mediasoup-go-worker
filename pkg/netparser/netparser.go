package netparser

type INetParser interface {
	WriteBuffer(payload []byte) error
	ReadBuffer(payload []byte) (int, error)
	Close() error
}
