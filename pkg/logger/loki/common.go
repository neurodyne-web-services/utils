package loki

import (
	"time"
)

const LOG_ENTRIES_CHAN_SIZE = 5000

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO  LogLevel = iota
	WARN  LogLevel = iota
	ERROR LogLevel = iota
	// Maximum level, disables sending or printing
	DISABLE LogLevel = iota
)

type ClientConfig struct {
	// E.g. http://localhost:3100/api/prom/push
	PushURL string
	// E.g. "{job=\"somejob\"}"
	Labels             string
	BatchWait          time.Duration
	BatchEntriesNumber int
	// Logs are sent to Promtail if the entry level is >= SendLevel
	SendLevel LogLevel
	// Logs are printed to stdout if the entry level is >= PrintLevel
	PrintLevel LogLevel
}

type Client interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
