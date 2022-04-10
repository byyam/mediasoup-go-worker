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

func (c *NetStrings) Close() error {
	if c.writeFile != nil {
		err := c.writeFile.Close()
		if err != nil {
			return err
		}
	}
	if c.readFile != nil {
		err := c.readFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
