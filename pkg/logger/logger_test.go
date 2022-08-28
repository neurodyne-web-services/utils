package logger

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_loki(t *testing.T) {
	b := &bytes.Buffer{}

	logger := MakeBufferLogger(b, "debug", "console")

	fmt.Println("Raw log:")
	logger.Error("foo")
	logger.Error("bar")

	fmt.Printf("Bufferred log: \n%s", b.String())
}
