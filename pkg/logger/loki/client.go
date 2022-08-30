package loki

import (
	"time"

	v1 "github.com/neurodyne-web-services/utils/pkg/logger/loki/genout/v1"
)

type Client interface {
	Debug(job string, args ...interface{})
	Info(job string, args ...interface{})
	Warn(job string, args ...interface{})
	Error(job string, args ...interface{})

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

type LogConnector struct {
	enable    bool
	verbosity string
}

func MakeLogConnector(verb string, ena bool) LogConnector {
	return LogConnector{ena, verb}
}

type LokiConfig struct {
	url     string
	ctype   string
	service string
}

type config struct {
	console LogConnector
	loki    LogConnector
	lcfg    LokiConfig
}

func MakeConfig(cons, loki LogConnector, lcfg LokiConfig) config {
	return config{
		console: cons,
		loki:    loki,
		lcfg:    lcfg,
	}
}
