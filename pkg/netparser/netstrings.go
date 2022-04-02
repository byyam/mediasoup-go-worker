package netparser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

type NetStrings struct {
	w                          io.Writer
	r                          *bufio.Reader
	separatorSymbol, endSymbol byte
}

func NewNetStrings(w io.Writer, r io.Reader) INetParser {
	return &NetStrings{
		w:               w,
		r:               bufio.NewReader(r),
		separatorSymbol: byte(':'),
		endSymbol:       byte(','),
	}
}

func (c NetStrings) WriteBuffer(payload []byte) error {
	var buffer bytes.Buffer
	length := strconv.FormatInt(int64(len(payload)), 10)
	buffer.WriteString(length)
	buffer.WriteByte(c.separatorSymbol)
	buffer.Write(payload)
	buffer.WriteByte(c.endSymbol)

	_, err := c.w.Write(buffer.Bytes())
	return err
}

func (c NetStrings) ReadBuffer() (payload []byte, err error) {
	begin, err := c.r.ReadString(c.separatorSymbol)
	if err != nil {
		return
	}
	if len(begin) < 1 {
		err = errors.New("invalid payload start")
		return
	}
	if separator := begin[len(begin)-1]; separator != c.separatorSymbol {
		err = errors.New("invalid payload separator")
		return
	}
	length, err := strconv.Atoi(begin[:len(begin)-1])
	if err != nil {
		return
	}
	payload = make([]byte, length)
	if _, err = io.ReadFull(c.r, payload); err != nil {
		return
	}
	end, err := c.r.ReadByte()
	if err != nil {
		return
	}
	if end != c.endSymbol {
		err = errors.New("invalid payload end")
		return
	}
	return
}
