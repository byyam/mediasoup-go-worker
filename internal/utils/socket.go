package utils

import (
	"net"
	"os"
)

func FileToConn(file *os.File) (net.Conn, error) {
	defer func() {
		_ = file.Close()
	}()
	return net.FileConn(file)
}
