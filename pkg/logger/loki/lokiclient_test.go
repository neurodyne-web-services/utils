package loki

import (
	"testing"
	"time"

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

	batch := batchConfig{
		BatchEntriesNumber: 8,
		BatchWait:          300,
	}

	loki := MakeLokiLogger(conf, zl, true, true, batch)

	loki.Debugf(job0, "My message is %s", "Hey There")
	loki.Debugf(job1, "My number is %d", 5)
	loki.Debugf(job1, "How are you, user %s ", "Ivan")

	time.Sleep(time.Second)
}
