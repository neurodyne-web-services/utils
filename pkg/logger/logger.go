package logger

import (
	"io"
	"strings"

	"github.com/neurodyne-web-services/utils/pkg/functional"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Type uint8

const (
	Console Type = iota
	JSON
)

const (
	timeFormat = "02 Jan 2006 15:04:05 MST"
)

type Config struct {
	Type  string
	Level string
}

// MakeLogger - simple customized console logger for dev.
func MakeLogger(encoding string, level zapcore.Level) (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         encoding,
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	return cfg.Build()
}

var DevConfig = zap.Config{
	EncoderConfig: zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.TimeEncoderOfLayout(timeFormat),

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	},
}

// NewPipedLogger - zap logger with multiple sinks.
// we use strings for log type and verbosity since those come from confings which are strings in most clients.
func NewPipedLogger(cfg zap.Config, config Config, sinks ...io.Writer) zapcore.Core {
	var enc zapcore.Encoder
	var level zapcore.Level

	// build proper encoder
	switch strings.ToLower(config.Type) {
	case "json":
		enc = zapcore.NewJSONEncoder(cfg.EncoderConfig)

	default:
		enc = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	}

	// select zap verbosity
	switch strings.ToLower(config.Level) {
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "debug":
		level = zapcore.DebugLevel
	case "pinic":
		level = zapcore.PanicLevel

	default:
		level = zapcore.InfoLevel
	}

	syncers := functional.Map(sinks, func(w io.Writer) zapcore.WriteSyncer {
		return zapcore.AddSync(w)
	})

	return zapcore.NewCore(
		enc,
		zap.CombineWriteSyncers(syncers...),
		level,
	)
}

// MakeExtLogger - a multiroute logger, which uses console
// and an external logger thru the Writer interface.
func MakeExtLogger(core zapcore.Core) *zap.Logger {
	return zap.New(core, zap.AddCaller())
}
