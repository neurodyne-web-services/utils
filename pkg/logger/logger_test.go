package logger_test

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	zaploki "github.com/neurodyne-web-services/zap-loki"
	"go.uber.org/zap"
)

func Test_piped_logger(t *testing.T) {
	t.Run("Console logger", func(_ *testing.T) {
		conf := logger.Config{
			Type:  "console",
			Level: "debug",
		}

		// loki := zaploki.New(context.Background(), zaploki.Config{
		// 	Url:          "http://localhost:3100",
		// 	BatchMaxSize: 100,
		// 	BatchMaxWait: 2 * time.Second,
		// 	Labels:       map[string]string{"app": "test", "env": "dev"},
		// })

		buf := &bytes.Buffer{}
		bufSink := &zaploki.BufferSink{Buf: buf}

		syncs := make(map[string]zaploki.SyncFactory)
		// syncs["loki"] = loki.Sink
		// syncs["console"] = func(_ *url.URL) (zap.Sink, error) {
		// 	return os.Stdout, nil
		// }
		syncs["buffer"] = func(_ *url.URL) (zap.Sink, error) {
			return bufSink, nil
		}

		core := logger.NewPipedLogger(logger.DevConfig, conf)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})

	t.Run("JSON logger", func(_ *testing.T) {
		t.Skip()

		conf := logger.Config{
			Type:  "json",
			Level: "debug",
		}

		core := logger.NewPipedLogger(logger.DevConfig, conf, nil)
		logger := logger.MakeExtLogger(core)

		logger.Error("foo")
	})
}
