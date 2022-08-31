package logger

import (
	"time"

	v1 "github.com/neurodyne-web-services/utils/pkg/logger/loki/genout/v1"
)

const (
	MIN_ENTRIES = 2
)

type serverResp struct {
	code int
	body []byte
}

type zapMsg struct {
	Level   string
	Time    time.Time
	Message string
}

type LokiConfig struct {
	Url     string
	Ctype   string
	Service string
	Batch   BatchConfig
}

func MakeLokiConfig(url, ctype, service string, batchSize int8) LokiConfig {
	return LokiConfig{
		Url:     url,
		Ctype:   ctype,
		Service: service,
		Batch:   BatchConfig{BatchSize: batchSize},
	}
}

type streamItem struct {
	labels string
	entry  *v1.Entry
}

type BatchConfig struct {
	BatchSize int8
}
