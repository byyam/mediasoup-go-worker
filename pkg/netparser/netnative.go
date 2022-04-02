package netparser

import (
	"encoding/binary"
	"io"
)

type NetNative struct {
	w            io.Writer
	r            io.Reader
	nativeEndian binary.ByteOrder
}

func NewNetNative(w io.Writer, r io.Reader, nativeEndian binary.ByteOrder) INetParser {
	return &NetNative{
		w:            w,
		r:            r,
		nativeEndian: nativeEndian,
	}
}

func (c NetNative) WriteBuffer(payload []byte) error {
	length := uint32(len(payload))
	if length == 0 {
		return nil
	}
	if err := binary.Write(c.w, c.nativeEndian, length); err != nil {
		return err
	}
	_, err := c.w.Write(payload)
	return err
}

func (c NetNative) ReadBuffer() (payload []byte, err error) {
	var payloadLen uint32
	if err = binary.Read(c.r, c.nativeEndian, &payloadLen); err != nil {
		return
	}
	payload = make([]byte, payloadLen)
	_, err = io.ReadFull(c.r, payload)
	return
}
