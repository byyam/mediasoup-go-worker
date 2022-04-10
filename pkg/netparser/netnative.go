package netparser

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"unsafe"
)

type NetNative struct {
	w                   io.Writer
	r                   io.Reader
	writeFile, readFile *os.File
	nativeEndian        binary.ByteOrder
}

func NewNetNative(w io.Writer, r io.Reader, nativeEndian binary.ByteOrder) INetParser {
	return &NetNative{
		w:            w,
		r:            r,
		nativeEndian: nativeEndian,
	}
}

func NewNetNativeFd(writeFd, readFd int, nativeEndian binary.ByteOrder) (INetParser, error) {
	var err error
	c := &NetNative{
		nativeEndian: nativeEndian,
	}
	c.writeFile = os.NewFile(uintptr(writeFd), "")
	c.readFile = os.NewFile(uintptr(readFd), "")
	c.w, err = net.FileConn(c.writeFile)
	if err != nil {
		return nil, err
	}
	c.r, err = net.FileConn(c.readFile)
	if err != nil {
		return nil, err
	}
	return c, nil
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

func (c NetNative) ReadBuffer(payload []byte) (int, error) {
	var payloadLen uint32
	if err := binary.Read(c.r, c.nativeEndian, &payloadLen); err != nil {
		return 0, err
	}
	return io.ReadFull(c.r, payload[:payloadLen])
}

func (c *NetNative) Close() error {
	var wErr, rErr error
	if c.writeFile != nil {
		wErr = c.writeFile.Close()
	}
	if c.readFile != nil {
		rErr = c.readFile.Close()
	}
	if wErr != nil || rErr != nil {
		return errors.New("close write/read file failed")
	}
	return nil
}

func HostByteOrder() binary.ByteOrder {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		return binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		return binary.BigEndian
	default:
		panic("Could not determine native endian.")
	}
}
