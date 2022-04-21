package ratecalculator

import (
	"log"
	"os"
)

type Logger interface {
	Error(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
}

func newLogger() Logger {
	return &logger{
		log: log.New(os.Stdout, "", log.LstdFlags),
	}
}

type logger struct {
	log *log.Logger
}

func (l *logger) Error(format string, args ...interface{}) {
	l.log.Printf("[Err]"+format, args...)
}
func (l *logger) Warn(format string, args ...interface{}) {
	l.log.Printf("[WRN]"+format, args...)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.log.Printf("[INF]"+format, args...)
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.log.Printf("[DBG]"+format, args...)
}

func (l *logger) Trace(format string, args ...interface{}) {
	l.log.Printf("[TRC]"+format, args...)
}
