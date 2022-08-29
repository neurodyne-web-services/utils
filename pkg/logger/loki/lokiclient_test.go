package loki

import (
	"bytes"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func Test_loki(t *testing.T) {

	conf := lokiConfig{
		url:   "http://localhost:3100/api/prom/push",
		ctype: "application/x-protobuf",
		level: zapcore.DebugLevel,
	}

	b := &bytes.Buffer{}
	zl := logger.MakeBufferLogger(b, "debug", "console")

	loki := MakeLokiClient(conf, zl)

	loki.Debugf("Hey there")
}
