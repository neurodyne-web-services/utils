package loki

import (
	"time"

	v1 "github.com/neurodyne-web-services/utils/pkg/logger/loki/genout/v1"
	"go.uber.org/zap/zapcore"
)

type Client interface {
	Debugf(job, template string, args ...interface{})
	Infof(job, template string, args ...interface{})
	Warnf(job, template string, args ...interface{})
	Errorf(job, template string, args ...interface{})
	Shutdown()
}

type serverResp struct {
	code int
	body []byte
}

type streamItem struct {
	labels string
	entry  *v1.Entry
}

type batchConfig struct {
	BatchSize int
	BatchWait time.Duration
}

func MakeBatchConfig(size int, wait time.Duration) batchConfig {
	return batchConfig{size, wait}
}

type lokiConfig struct {
	enableLoki    bool
	enableConsole bool
	url           string
	ctype         string
	service       string
	level         zapcore.Level
}

func MakeLokiConfig(enaLoki, enaConsole bool, url, ctype, service string, level zapcore.Level) lokiConfig {
	return lokiConfig{
		enableLoki:    enaLoki,
		enableConsole: enaConsole,
		url:           url,
		ctype:         ctype,
		service:       service,
		level:         level,
	}
}
