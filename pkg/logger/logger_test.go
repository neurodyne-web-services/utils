package logger

import (
	"bytes"
	"fmt"
	"testing"

	"go.uber.org/zap"
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

	loki := dummyLogger{}
	defer loki.Sync()

	logger := MakeExtLogger(loki, "debug", "console")

	fmt.Println("Raw log:")
	logger.Error("foo")
	logger.Debug("My number is ", zap.Int("Number", 4))
	logger.Warn("bar")
	logger.Info("baz")
}

var items []string

// A dummy logger, which implements a zap.WriteSyncer interface
type dummyLogger struct{}

func (l dummyLogger) Write(p []byte) (n int, err error) {
	// fmt.Printf(">>>> Loki logger: %s \n", string(p))
	if p != nil {
		items = append(items, string(p))
	}

	return 0, nil
}

func (l dummyLogger) Sync() error {
	fmt.Printf(">>>> SYNC: items size: %d", len(items))
	return nil
}
