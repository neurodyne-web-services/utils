package logger_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
)

const (
	url       = "http://localhost:3100/api/prom/push"
	ctype     = "application/x-protobuf"
	lokiMode  = logger.DEV
	mode      = "dev"
	verbosity = "debug"
	batchSize = 4
	loops     = 1
)

func Test_buff(t *testing.T) {
	t.Skip()

	b := &bytes.Buffer{}

	core := logger.NewCustomLogger(logger.DevConfig, b, mode, verbosity)
	logger := logger.MakeExtLogger(core)

	fmt.Println(" >>>>> Raw log:")
	logger.Error("foo")
	logger.Error("bar")

	fmt.Printf(">>>>> Bufferred log: \n%s", b.String())
}

func Test_zap(t *testing.T) {
	conf := logger.MakeLokiConfig(lokiMode, url, ctype, batchSize)

	loki := logger.MakeLokiSyncer(conf)
	defer loki.Sync()

	core := logger.NewCustomLogger(logger.DevConfig, loki, mode, verbosity)
	zl := logger.MakeExtLogger(core)
	logger := zl.Sugar()

	logger.Warn("Starting test...")

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
