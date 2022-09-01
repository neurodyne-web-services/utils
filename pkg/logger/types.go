package logger

import (
	"time"
)

const (
	MIN_ENTRIES = 4
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
