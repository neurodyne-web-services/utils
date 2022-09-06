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
	Service string
	Job     string
	Message string
	Time    time.Time
}

type LokiConfig struct {
	Enable    bool
	Url       string
	Ctype     string
	BatchSize int
}

type LoggerConfig struct {
	Output string
	Level  string
}

func MakeLoggerConfig(lvl, out string) LoggerConfig {
	return LoggerConfig{Output: out, Level: lvl}
}

func MakeLokiConfig(ena bool, url, ctype string, batchSize int) LokiConfig {
	return LokiConfig{
		Enable:    ena,
		Url:       url,
		Ctype:     ctype,
		BatchSize: batchSize,
	}
}
