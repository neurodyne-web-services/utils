package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	timeFormat = "02 Jan 2006 15:04:05 MST"
)

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

// NewPipedLogger - console/json logger with an extra pipe.
func NewPipedLogger(cfg zap.Config, pipeTo io.Writer, loggerType LoggerType, level zapcore.Level) zapcore.Core {
	var enc zapcore.Encoder

	switch loggerType {
	case Console:
		enc = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	case JSON:
		enc = zapcore.NewJSONEncoder(cfg.EncoderConfig)

	default:
		enc = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	}

	return zapcore.NewCore(
		enc,
		zap.CombineWriteSyncers(os.Stderr, zapcore.AddSync(pipeTo)),
		level,
	)
}

// MakeExtLogger - a multiroute logger, which uses console
// and an external logger thru the Writer interface.
func MakeExtLogger(core zapcore.Core) *zap.Logger {
	return zap.New(core, zap.AddCaller())
}
