package zlog

import (
	"io"
	"os"
	"path"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger   *zap.Logger
	initOnce sync.Once
)

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

func Init(config Config) {
	initOnce.Do(func() {
		// lumberjack.Logger is already safe for concurrent use, so we don't need to
		// lock it.
		fileWriter := newRollingFile(config)
		w := zapcore.AddSync(fileWriter)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.InfoLevel,
		)
		logger = zap.New(core)
	})
}

// GetLogger for middleware like gin recovery or log print
func GetLogger() *zap.Logger { // no need to protect nil pointer if not init logger, low error
	return logger
}

func newRollingFile(config Config) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		panic("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}
