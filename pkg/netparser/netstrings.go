package netparser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
)

const (
	separatorSymbol = byte(':')
	endSymbol       = byte(',')
)

type NetStrings struct {
	w                          io.Writer
	r                          *bufio.Reader
	writeFile, readFile        *os.File
	separatorSymbol, endSymbol byte
}

func NewNetStrings(w io.Writer, r io.Reader) INetParser {
	return &NetStrings{
		w:               w,
		r:               bufio.NewReader(r),
		separatorSymbol: separatorSymbol,
		endSymbol:       endSymbol,
	}
}

func NewNetStringsFd(writeFd, readFd int) (INetParser, error) {
	var err error
	c := &NetStrings{
		separatorSymbol: separatorSymbol,
		endSymbol:       endSymbol,
	}
	c.writeFile = os.NewFile(uintptr(writeFd), "")
	c.readFile = os.NewFile(uintptr(readFd), "")
	c.w, err = net.FileConn(c.writeFile)
	if err != nil {
		return nil, err
	}
	reader, err := net.FileConn(c.readFile)
	if err != nil {
		return nil, err
	}
	c.r = bufio.NewReader(reader)
	return c, nil
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

func (c NetStrings) ReadBuffer(payload []byte) (int, error) {
	begin, err := c.r.ReadString(c.separatorSymbol)
	if err != nil {
		return 0, err
	}
	if len(begin) < 1 {
		return 0, errors.New("invalid payload start")
	}
	if separator := begin[len(begin)-1]; separator != c.separatorSymbol {
		return 0, errors.New("invalid payload separator")
	}
	length, err := strconv.Atoi(begin[:len(begin)-1])
	if err != nil {
		return 0, err
	}
	if _, err = io.ReadFull(c.r, payload[:length]); err != nil {
		return 0, err
	}
	end, err := c.r.ReadByte()
	if err != nil {
		return 0, err
	}
	if end != c.endSymbol {
		return 0, errors.New("invalid payload end")
	}
	return length, nil
}

func (c *NetStrings) Close() error {
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
