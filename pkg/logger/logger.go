package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MakeLogger - simple customized console logger for dev
func MakeLogger(verbosity, encoding string) (*zap.Logger, error) {

	level := GetZapLevel(verbosity)

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

	// Add colors in for console
	config := zap.Config{
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

	return zapcore.NewCore(
		getZapEncoder(encoding, config.EncoderConfig),
		zap.CombineWriteSyncers(os.Stderr, zapcore.AddSync(pipeTo)),
		GetZapLevel(verbosity),
	)
}

// MakeExtLogger - a multiroute logger, which uses console
// and an external logger thru the Writer interface
func MakeExtLogger(w io.Writer, verb, enc string) *zap.Logger {
	return zap.New(newCustomLogger(w, verb, enc), zap.AddCaller())
}

// GetZapLevel - returns a Zap logger verbosity level based
// on the input string
func GetZapLevel(verb string) zapcore.Level {
	var level zapcore.Level

	switch verb {

	case "debug":
		level = zapcore.DebugLevel

	case "fatal":
		level = zapcore.FatalLevel

	case "error":
		level = zapcore.ErrorLevel

	case "warn":
		level = zapcore.WarnLevel

	case "info":
		level = zapcore.InfoLevel

	default:
		level = zapcore.InfoLevel
	}

	return level
}

func getZapEncoder(encoder string, cfg zapcore.EncoderConfig) zapcore.Encoder {
	var enc zapcore.Encoder

	// Build a proper logger type
	switch encoder {

	case "console":
		enc = zapcore.NewConsoleEncoder(cfg)

	case "json":
		enc = zapcore.NewJSONEncoder(cfg)

	default:
		enc = zapcore.NewConsoleEncoder(cfg)
	}

	return enc
}
