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
	Caller  string
	Message string
	Time    time.Time
}

type streamItem struct {
	labels string
	entry  *v1.Entry
}

type LokiConfig struct {
	Url       string
	Ctype     string
	Service   string
	BatchSize int
}

func MakeLokiConfig(url, ctype, service string, batchSize int) LokiConfig {
	return LokiConfig{
		Url:       url,
		Ctype:     ctype,
		Service:   service,
		BatchSize: batchSize,
	}
}
