package logger_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
)

func Test_piped_logger(t *testing.T) {
	t.Run("Console logger", func(_ *testing.T) {
		b := &bytes.Buffer{}

		conf := logger.Config{
			Type:  "console",
			Level: "debug",
		}

		core := logger.NewPipedLogger(logger.DevConfig, conf, b, os.Stdout)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})

	t.Run("JSON logger", func(_ *testing.T) {
		b := &bytes.Buffer{}

		conf := logger.Config{
			Type:  "json",
			Level: "debug",
		}

		core := logger.NewPipedLogger(logger.DevConfig, conf, b, os.Stdout)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})
}
