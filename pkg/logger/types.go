package logger

import (
	"time"
)

const (
	MinEntries = 4
)

type LokiMode uint8

const (
	DEV LokiMode = iota
	PROD
)

type serverResp struct {
	code int
	body []byte
}

type labelEntry struct {
	Env     string
	Service string
}

type zapMsg struct {
	Level   string
	Caller  string
	Message string
	Env     string
	Service string
	Time    time.Time
}

type LokiConfig struct {
	Mode      LokiMode
	URL       string
	Ctype     string
	BatchSize int
}

type LoggerConfig struct {
	Output string
	Level  string
}

func MakeLoggerConfig(mode LokiMode, lvl string) LoggerConfig {
	if mode == PROD {
		return LoggerConfig{Output: "json", Level: lvl}
	}
	return LoggerConfig{Output: "console", Level: lvl}
}

func MakeLokiConfig(mode LokiMode, url, ctype string, batchSize int) LokiConfig {
	return LokiConfig{
		Mode:      mode,
		URL:       url,
		Ctype:     ctype,
		BatchSize: batchSize,
	}
}
