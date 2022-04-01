package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	// TraceLevel defines trace log level.
	TraceLevel = zerolog.TraceLevel
	// DebugLevel defines debug log level.
	DebugLevel = zerolog.DebugLevel
	// InfoLevel defines info log level.
	InfoLevel = zerolog.InfoLevel
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
	// Scopes defines default log scopes.
	NoneScopeLimit = true // if not set, all scope is valid
	NoneScopeLevel = DefaultLevel
	Scopes         = make(map[string]bool)
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

func SetScopes(scopes ...string) {
	NoneScopeLimit = false
	for _, s := range scopes {
		Scopes[s] = true
	}
}

type Logger interface {
	Trace(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

type defaultLogger struct {
	logger zerolog.Logger
}

func newDefaultLogger(scope string, ids ...interface{}) Logger {

	context := zerolog.New(NewLoggerWriter()).With().Timestamp()

	var caller string
	if len(ids) > 0 {
		var scopeIds []string
		for _, id := range ids {
			scopeIds = append(scopeIds, fmt.Sprintf("%v", id))
		}
		idsStr := strings.Join(scopeIds, "|")
		caller = fmt.Sprintf("%s[%s]", scope, idsStr)
	} else {
		caller = scope
	}
	if len(scope) > 0 {
		context = context.Str(zerolog.CallerFieldName, caller)
	}

	logLevel := DefaultLevel
	if !Scopes[scope] && !NoneScopeLimit {
		logLevel = NoneScopeLevel
	}
	return &defaultLogger{
		logger: context.Logger().Level(logLevel),
	}
}

func (l defaultLogger) Trace(format string, v ...interface{}) {
	l.logger.Trace().Msgf(format, v...)
}

func (l defaultLogger) Debug(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

func (l defaultLogger) Info(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

func (l defaultLogger) Warn(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

func (l defaultLogger) Error(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}
