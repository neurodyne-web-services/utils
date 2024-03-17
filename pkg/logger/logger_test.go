package logger_test

import (
	"bytes"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	"go.uber.org/zap/zapcore"
)

const (
	url       = "http://localhost:3100/api/prom/push"
	ctype     = "application/x-protobuf"
	batchSize = 4
	loops     = 1
)

func Test_piped_logger(t *testing.T) {
	b := &bytes.Buffer{}

	var lt logger.LoggerType
	t.Run("Console logger", func(_ *testing.T) {
		core := logger.NewPipedLogger(logger.DevConfig, b, lt, zapcore.DebugLevel)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})

	t.Run("JSON logger", func(_ *testing.T) {
		core := logger.NewPipedLogger(logger.DevConfig, b, lt, zapcore.DebugLevel)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})
}
