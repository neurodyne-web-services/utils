package logger

import (
	"strings"
	"time"
)

const (
	MinEntries = 4
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
	Mode      string
	URL       string
	Ctype     string
	BatchSize int
}

type Config struct {
	Output string
	Level  string
}

func MakeLoggerConfig(mode, lvl string) Config {
	if strings.ToLower(mode) == "prod" {
		return Config{Output: "json", Level: lvl}
	}
	return Config{Output: "console", Level: lvl}
}

func MakeLokiConfig(mode, url, ctype string, batchSize int) LokiConfig {
	return LokiConfig{
		Mode:      mode,
		URL:       url,
		Ctype:     ctype,
		BatchSize: batchSize,
	}
}
