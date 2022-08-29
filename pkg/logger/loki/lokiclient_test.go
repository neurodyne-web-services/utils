package loki

import (
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

const (
	url     = "http://localhost:3100/api/prom/push"
	ctype   = "application/x-protobuf"
	service = "drevo"
)

func Test_loki(t *testing.T) {

	conf := MakeLokiConfig(url, ctype, service, zapcore.DebugLevel)

	job0 := "list"
	job1 := "put"

	zl, err := logger.MakeLogger("debug", "console")
	assert.NoError(t, err)

	loki := MakeLokiLogger(conf, zl, true, true)

	loki.Debugf(job0, "My message is %s", "Hey There")
	loki.Debugf(job1, "My number is %d", 5)
}
