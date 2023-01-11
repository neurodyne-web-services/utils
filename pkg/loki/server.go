package loki

import (
	"context"

	"github.com/neurodyne-web-services/utils/pkg/logger"
	v1 "github.com/neurodyne-web-services/utils/pkg/loki/loki/v1"

	"go.uber.org/zap"
)

type Server struct {
	v1.UnimplementedLogServiceServer
	zl *zap.SugaredLogger
}

func BuildServer(zl *zap.SugaredLogger) Server {
	return Server{zl: zl}
}

func (s Server) Push(_ context.Context, item *v1.Item) (*v1.Empty, error) {
	switch item.Level {
	case v1.Item_INFO:
		s.zl.Infow(item.Msg, "env", item.Env, "service", item.Service)

	case v1.Item_WARN:
		s.zl.Warnw(item.Msg, "env", item.Env, "service", item.Service)

	case v1.Item_ERROR:
		s.zl.Errorw(item.Msg, "env", item.Env, "service", item.Service)

	case v1.Item_DEBUG:
		s.zl.Debugw(item.Msg, "env", item.Env, "service", item.Service)

	case v1.Item_FATAL:
		s.zl.Fatalw(item.Msg, "env", item.Env, "service", item.Service)
	}

	s.zl.Sync()
	return &v1.Empty{}, nil
}

// BuildLokiLogger - returns an instancf of Zap logger with Loki.
func BuildLokiLogger(mode logger.LokiMode, lvl, url string, batchSize int) *zap.SugaredLogger {
	conf := logger.MakeLokiConfig(mode, url, "application/x-protobuf", batchSize)
	loki := logger.MakeLokiSyncer(conf)
	zl := logger.MakeExtLogger(loki, logger.MakeLoggerConfig(mode, lvl))

	return zl.Sugar()
}
