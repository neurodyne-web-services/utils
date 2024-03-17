package logger

import (
	"time"
)

type Type uint8

const (
	Console Type = iota
	JSON
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
	URL         string
	ContentType string
	BatchSize   int
}

type Config struct {
	Output string
	Level  string
}
