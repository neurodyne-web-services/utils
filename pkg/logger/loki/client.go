package loki

import (
	v1 "github.com/neurodyne-web-services/utils/pkg/logger/loki/genout/v1"
)

type Client interface {
	Debug(job string, args ...interface{})
	Info(job string, args ...interface{})
	Warn(job string, args ...interface{})
	Error(job string, args ...interface{})
	Fatal(job string, args ...interface{})

	Debugf(job, template string, args ...interface{})
	Infof(job, template string, args ...interface{})
	Warnf(job, template string, args ...interface{})
	Errorf(job, template string, args ...interface{})
	Fatalf(job, template string, args ...interface{})
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

type BatchConfig struct {
	BatchSize       int
	BatchTimeoutSec int
}

func MakeBatchConfig(size, timeout int) BatchConfig {
	return BatchConfig{size, timeout}
}

type LogConnector struct {
	Enable bool
	Level  string
}

func MakeLogConnector(verb string, ena bool) LogConnector {
	return LogConnector{ena, verb}
}

type LokiConfig struct {
	Url     string
	Ctype   string
	Service string
	Batch   BatchConfig
}

type LogConfig struct {
	Console  LogConnector
	Loki     LogConnector
	LokiConf LokiConfig
}

func MakeConfig(cons, loki LogConnector, lcfg LokiConfig) LogConfig {
	return LogConfig{
		Console:  cons,
		Loki:     loki,
		LokiConf: lcfg,
	}
}
