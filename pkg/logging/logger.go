package logging

import (
	"bytes"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MakeLogger - simple customized console logger for dev
func MakeLogger(verbosity, encoding string) (*zap.Logger, error) {

	var level zapcore.Level

	switch verbosity {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	default:
		level = zapcore.ErrorLevel
	}

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

func newCustomLogger(pipeTo io.Writer, verbosity, encoding string) zapcore.Core {

	var enc zapcore.Encoder
	var level zapcore.Level

	switch verbosity {

	case "debug":
		level = zapcore.DebugLevel

	case "err":
		level = zapcore.WarnLevel

	case "warn":
		level = zapcore.WarnLevel

	case "info":
		level = zapcore.WarnLevel

	default:
		level = zapcore.ErrorLevel
	}

	// Add colors in for console
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Build a proper logger type
	switch encoding {

	case "console":
		enc = zapcore.NewConsoleEncoder(config)

	case "json":
		enc = zapcore.NewJSONEncoder(config)

	default:
		enc = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(
		enc,
		zap.CombineWriteSyncers(os.Stderr, zapcore.AddSync(pipeTo)),
		level,
	)
}

// MakeBufferLogger - a multiroute logger, which supports JSON/Console stdout and writing to buffer
func MakeBufferLogger(b *bytes.Buffer, verb, enc string) *zap.Logger {
	return zap.New(newCustomLogger(b, verb, enc))
}
