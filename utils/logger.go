package utils

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

const (
	// DebugLevel defines debug log level.
	DebugLevel = zerolog.DebugLevel
	// WarnLevel defines warn log level.
	WarnLevel = zerolog.WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel = zerolog.ErrorLevel
	// Disabled disables the logger.
	Disabled = zerolog.Disabled
)

var (
	// DefaultLevel defines default log level.
	DefaultLevel = DebugLevel
	// NewLogger defines function to create logger instance.
	NewLogger = newDefaultLogger
	// NewLoggerWriter defines function to create logger writer.
	NewLoggerWriter = func() io.Writer {
		writer := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if hideDate := os.Getenv("DEBUG_HIDE_DATE"); len(hideDate) > 0 {
			val, err := strconv.ParseBool(hideDate)
			if err == nil && val {
				writer.FormatTimestamp = func(interface{}) string {
					return ""
				}
			}
		}

		if color := os.Getenv("DEBUG_COLORS"); len(color) > 0 {
			val, err := strconv.ParseBool(color)
			if err == nil {
				writer.NoColor = !val
			}
		}

		return writer
	}
)

type Logger interface {
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

type defaultLogger struct {
	logger zerolog.Logger
}

func newDefaultLogger(scope string) Logger {

	context := zerolog.New(NewLoggerWriter()).With().Timestamp()

	if len(scope) > 0 {
		context = context.Str(zerolog.CallerFieldName, scope)
	}

	return &defaultLogger{
		logger: context.Logger().Level(DefaultLevel),
	}
}

func (l defaultLogger) Debug(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

func (l defaultLogger) Warn(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

func (l defaultLogger) Error(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}
