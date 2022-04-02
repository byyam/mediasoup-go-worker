package netparser

type INetParser interface {
	WriteBuffer(payload []byte) error
	ReadBuffer() ([]byte, error)
}
