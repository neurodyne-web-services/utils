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
	verbosity = "debug"
	batchSize = 4
)

func Test_buff(t *testing.T) {
	t.Skip()

	b := &bytes.Buffer{}

	logger := MakeExtLogger(b, MakeLoggerConfig("debug", "console"))

	fmt.Println(" >>>>> Raw log:")
	logger.Error("foo")
	logger.Error("bar")

	fmt.Printf(">>>>> Bufferred log: \n%s", b.String())
}

func Test_zap(t *testing.T) {

	conf := MakeLokiConfig(true, url, ctype, batchSize)

	loki := MakeLokiSyncer(conf)
	defer loki.Sync()

	zl := MakeExtLogger(loki, MakeLoggerConfig("debug", "json"))
	logger := zl.Sugar()

	logger.Info("baz", "env", "prod")

	logger.Debugw(fmt.Sprintf("Hello, %s", "Boris"),
		"env", "dev",
		"service", "front")

	logger.Debugw(fmt.Sprintf("Hello, %s", "Emma"),
		"env", "dev",
		"service", "back")

	logger.Debugw(fmt.Sprintf("Hello, %s", "Ivan"),
		"env", "prod",
		"service", "front")

	logger.Warnw("My warn",
		"env", "prod",
		"service", "back")

	logger.Error("My error")

	time.Sleep(time.Second)
}
