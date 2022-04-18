package rtpparser

import (
	"encoding/binary"
	"unsafe"
)

var (
	hostByteOrder binary.ByteOrder
)

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		hostByteOrder = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		hostByteOrder = binary.BigEndian
	default:
		panic("Could not determine native endian.")
	}
}
