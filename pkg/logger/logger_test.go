package logger_test

import (
	"bytes"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func Test_piped_logger(t *testing.T) {
	b := &bytes.Buffer{}

	var lt logger.Type
	t.Run("Console logger", func(_ *testing.T) {
		lt = logger.Console

		core := logger.NewPipedLogger(logger.DevConfig, lt, zapcore.DebugLevel, b)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})

	t.Run("JSON logger", func(_ *testing.T) {
		lt = logger.JSON

		core := logger.NewPipedLogger(logger.DevConfig, lt, zapcore.DebugLevel, b)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})
}
