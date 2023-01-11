package loki

import (
	"context"
	"os"
	"time"

	v1 "github.com/neurodyne-web-services/utils/pkg/loki/loki/v1"
)

const (
	timeout = 5 // sec
)

type Client struct {
	C v1.LogServiceClient
}

func (c Client) Push(level v1.Item_Level) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	item := &v1.Item{
		Level:   level,
		Env:     os.Getenv("ARTEMIDA_ENV"),
		Service: "front",
		Msg:     "Hello world",
	}

	_, err := c.C.Push(ctx, item)
	if err != nil {
		return err
	}

	return nil
}
