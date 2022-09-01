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
	Job     string
	Message string
	Time    time.Time
}

type LokiConfig struct {
	Enable    bool
	Url       string
	Ctype     string
	Service   string
	BatchSize int
}

func MakeLokiConfig(ena bool, url, ctype, service string, batchSize int) LokiConfig {
	return LokiConfig{
		Enable:    ena,
		Url:       url,
		Ctype:     ctype,
		Service:   service,
		BatchSize: batchSize,
	}
}
