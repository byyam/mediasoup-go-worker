package zerowrapper

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logLevel = zerolog.InfoLevel
	initOnce sync.Once
	logger   = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: TimeFormatRFC3339}).With().Timestamp().Logger() // set default logger
)

const (
	TimeFormatRFC3339 = "2006-01-02T15:04:05.999Z"
	TimeFormatDefault = "2006-01-02 15:04:05.000000"
)

func InitLog(config Config) {
	initOnce.Do(func() {
		// set from os.env
		getLevel()
		zerolog.SetGlobalLevel(logLevel)
		// set format from config
		zerolog.TimeFieldFormat = config.LogTimeFieldFormat
		if config.ErrorStackMarshaler {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			zerolog.ErrorStackFieldName = "error.stack_trace"
		}
		// pretty format
		consoleWriter := zerolog.ConsoleWriter{}
		switch config.LogTimeFieldFormat {
		case zerolog.TimeFormatUnixMicro:
			consoleWriter.TimeFormat = TimeFormatDefault
		case zerolog.TimeFormatUnixMs:
			consoleWriter.TimeFormat = TimeFormatRFC3339
		}
		consoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		var writers []io.Writer
		if config.ConsoleLoggingEnabled {
			consoleWriter.Out = os.Stdout
			writers = append(writers, consoleWriter)
		}
		if config.FileLoggingEnabled {
			fileWriter := newRollingFile(config)
			writers = append(writers, fileWriter)
		}
		mw := io.MultiWriter(writers...)

		logger = zerolog.New(mw)

		newLogger := logger.With().Timestamp().Logger()
		newLogger.Info().
			Bool("fileLogging", config.FileLoggingEnabled).
			Str("logDirectory", config.Directory).
			Str("fileName", config.Filename).
			Int("maxSizeMB", config.MaxSize).
			Int("maxBackups", config.MaxBackups).
			Int("maxAgeInDays", config.MaxAge).
			Str("logTimeFormat", config.LogTimeFieldFormat).
			Bool("errorStackMarshaler", config.ErrorStackMarshaler).
			Str("logLevel", logLevel.String()).
			Msg("zero logging configured")
	})
}

func getLevel() {
	setLevel := os.Getenv("WORKER_LOG")
	switch setLevel {
	case "i":
		logLevel = zerolog.InfoLevel
	case "d":
		logLevel = zerolog.DebugLevel
	case "e":
		logLevel = zerolog.ErrorLevel
	case "t":
		logLevel = zerolog.TraceLevel
	case "w":
		logLevel = zerolog.WarnLevel
	case "off":
		logLevel = zerolog.Disabled
	}
}

func NewScope(scope string, ids ...interface{}) zerolog.Logger {
	logCtx := logger.With().Timestamp()

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
		logCtx = logCtx.Str(zerolog.CallerFieldName, caller)
	}
	return logCtx.Logger()
}

// Config - Configuration for logging
type Config struct {
	// ConsoleLoggingEnabled enable console logging
	ConsoleLoggingEnabled bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to when file logging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
	// LogTimeFieldFormat UNIX Time is faster and smaller than most timestamps, if not set use the default UNIX time
	LogTimeFieldFormat string
	// ErrorStackMarshaler extract the stack from err if any.
	ErrorStackMarshaler bool
}

func newRollingFile(config Config) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}
