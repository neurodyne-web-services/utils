package logger

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/neurodyne-web-services/utils/pkg/random"
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

	logger := MakeExtLogger(b, MakeLoggerConfig("debug", "console"))

	fmt.Println(" >>>>> Raw log:")
	logger.Error("foo")
	logger.Error("bar")

	fmt.Printf(">>>>> Bufferred log: \n%s", b.String())
}

func Test_zap(t *testing.T) {

	conf := MakeLokiConfig(true, url, ctype, service, batchSize)

	loki := MakeLokiSyncer(conf)
	defer loki.Sync()

	zl := MakeExtLogger(loki, MakeLoggerConfig("debug", "json"))
	logger := zl.Sugar()

	logger.Info("baz")
	logger.Debugw(fmt.Sprintf("Hello, %s", "Boris"),
		"job", random.GenRandomName("job"))

	logger.Warnw("My warn", "job", random.GenRandomName("job"))
	logger.Error("My error")

	time.Sleep(time.Second)
}
