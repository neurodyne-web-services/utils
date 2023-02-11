package loki_test

import (
	"net"
	"os"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	"github.com/neurodyne-web-services/utils/pkg/loki"
	v1 "github.com/neurodyne-web-services/utils/pkg/loki/loki/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	url       = "http://0.0.0.0:3100/api/prom/push"
	level     = "debug"
	batchSize = 4
)

func Test_loki(t *testing.T) {
	zl := loki.BuildLokiLogger(logger.DEV, level, url, batchSize)
	defer zl.Sync()

	RunServer(t, zl)

	// Create client
	cc, err := grpc.Dial(os.Getenv("GRPC_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer cc.Close()

	client := loki.Client{C: v1.NewLogServiceClient(cc)}

	data := []struct {
		name   string
		action v1.Item_Level
	}{
		{"debug", v1.Item_DEBUG},
		{"error", v1.Item_ERROR},
		{"warn", v1.Item_WARN},
		{"info", v1.Item_INFO},
		// {"fatal", v1.Item_FATAL},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			switch d.action {
			case v1.Item_FATAL:
				err = client.Push(v1.Item_FATAL)

			case v1.Item_DEBUG:
				err = client.Push(v1.Item_DEBUG)

			case v1.Item_ERROR:
				err = client.Push(v1.Item_ERROR)

			case v1.Item_WARN:
				err = client.Push(v1.Item_WARN)

			case v1.Item_INFO:
				err = client.Push(v1.Item_INFO)
			}

			assert.NoError(t, err)
		})
	}
}

func RunServer(t *testing.T, zl *zap.SugaredLogger) {
	t.Helper()
	logServer := loki.BuildServer(zl)

	opts := &[]grpc.ServerOption{}

	lis, err := net.Listen("tcp", os.Getenv("GRPC_URL"))
	assert.NoError(t, err)

	s := grpc.NewServer(*opts...)
	v1.RegisterLogServiceServer(s, logServer)

	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()
}
