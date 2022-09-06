package loki

import (
	"testing"
	"time"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	url       = "http://localhost:3100/api/prom/push"
	ctype     = "application/x-protobuf"
	verbosity = "debug"
)

func Test_loki(t *testing.T) {

	consConn := MakeLogConnector(verbosity, true)
	lokiConn := MakeLogConnector(verbosity, true)

	lokiConf := LokiConfig{
		Url:   url,
		Ctype: ctype,
		Batch: MakeBatchConfig(8, 30),
	}

	conf := MakeConfig(consConn, lokiConn, lokiConf)

	job0 := "list"
	job1 := "put"

	zl, err := logger.MakeLogger("debug", "console")
	assert.NoError(t, err)

	loki := MakeLokiLogger(conf, zl)

	loki.Warnf(job0, "My message is %s", "Hey There")
	loki.Infof(job0, "Starting the test...")

	for i := 0; i < 3; i++ {
		loki.Debugf(job1, "My number is %d", i)
	}

	loki.Errorf(job0, "Done logging, %s ", "Ivan")

	time.Sleep(time.Second)
	loki.Shutdown()
}
