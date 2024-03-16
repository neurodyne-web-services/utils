package logger_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	"go.uber.org/zap/zapcore"
)

const (
	url       = "http://localhost:3100/api/prom/push"
	ctype     = "application/x-protobuf"
	mode      = "json"
	verbosity = "debug"
	batchSize = 4
	loops     = 1
)

func Test_buffered_logger(t *testing.T) {
	b := &bytes.Buffer{}

	t.Run("Console logger", func(_ *testing.T) {
		core := logger.NewConsolePipeLogger(logger.DevConfig, b, zapcore.DebugLevel)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})

	t.Run("JSON logger", func(_ *testing.T) {
		core := logger.NewJSONPipeLogger(logger.DevConfig, b, zapcore.DebugLevel)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})
}

func Test_loki_logger(t *testing.T) {
	conf := logger.MakeLokiConfig(mode, url, ctype, batchSize)

	loki := logger.MakeLokiSyncer(conf)
	defer loki.Sync()

	core := logger.NewConsolePipeLogger(logger.DevConfig, loki, zapcore.DebugLevel)
	zl := logger.MakeExtLogger(core)
	logger := zl.Sugar()

	logger.Warn("Starting test...")
	logger.Infof("Log level: %s", logger.Level())

	for i := 0; i < loops; i++ {
		logger.Infow(fmt.Sprintf("PROD Info value, %d", i),
			"env", "prod",
			"service", "front")

		logger.Warnw(fmt.Sprintf("PROD Warn value, %d", i),
			"env", "prod",
			"service", "front")

		logger.Errorw(fmt.Sprintf("PROD Error value, %d", i),
			"env", "prod",
			"service", "front")

		logger.Debugw(fmt.Sprintf("PROD Debug value, %d", i),
			"env", "prod",
			"service", "front")

		logger.Debugw(fmt.Sprintf("Console value: %d", i),
			"service", "front")

		logger.Infow(fmt.Sprintf("DEV Info value, %d", i),
			"env", "dev",
			"service", "back")

		logger.Warnw(fmt.Sprintf("DEV Warn value, %d", i),
			"env", "dev",
			"service", "back")

		logger.Errorw(fmt.Sprintf("DEV Error value, %d", i),
			"env", "dev",
			"service", "back")

		logger.Debugw(fmt.Sprintf("DEV Debug value, %d", i),
			"env", "dev",
			"service", "back")
	}
}
