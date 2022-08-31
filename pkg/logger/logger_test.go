package logger

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

const (
	url       = "http://localhost:3100/api/prom/push"
	ctype     = "application/x-protobuf"
	service   = "drevo"
	verbosity = "debug"
	batchSize = 4
)

func Test_buff(t *testing.T) {
	t.Skip()

	b := &bytes.Buffer{}

	logger := MakeBufferLogger(b, "debug", "console")

	fmt.Println("Raw log:")
	logger.Error("foo")
	logger.Error("bar")

	fmt.Printf("Bufferred log: \n%s", b.String())
}

func Test_loki(t *testing.T) {

	conf := MakeLokiConfig(url, ctype, service, batchSize)

	loki := MakeLokiSyncer(conf)
	defer loki.Sync()

	zl := MakeExtLogger(loki, "debug", "json")
	logger := zl.Sugar()

	logger.Error("foo")
	// logger.Debugf("My number is ", zap.Int("Number", 4))
	logger.Debugf("My Number is %d", 4)
	logger.Warn("bar")
	logger.Info("baz")

	for i := 0; i < 16; i++ {
		logger.Debugf("My value is %d", i)
	}

	time.Sleep(2 * time.Second)
}
