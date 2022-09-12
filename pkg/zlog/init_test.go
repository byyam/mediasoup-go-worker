package zlog

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init()

}

func TestGetLogger(t *testing.T) {
	Init()

	logger := GetLogger()
	logger.Info("testing get logger")
}
