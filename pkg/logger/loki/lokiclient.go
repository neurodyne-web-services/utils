package loki

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/golang/snappy"
	"github.com/neurodyne-web-services/utils/pkg/logger"
	v1 "github.com/neurodyne-web-services/utils/pkg/logger/loki/genout/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	MIN_ENTRIES = 2
	MAX_ENTRIES = 8
	timeout     = 5 // sec
)

type LokiLogger struct {
	Client

	conf      config
	http      http.Client
	zl        *zap.SugaredLogger
	level     zapcore.Level
	entries   chan streamItem
	done      chan struct{}
	waitGroup sync.WaitGroup
}

// MakeLokiLogger - factory for Loki client
func MakeLokiLogger(conf config, zl *zap.Logger) *LokiLogger {

	client := LokiLogger{
		conf:    conf,
		zl:      zl.Sugar(),
		done:    make(chan struct{}),
		entries: make(chan streamItem, MAX_ENTRIES),
		level:   logger.GetZapLevel(conf.loki.Verbosity),
	}

	if conf.loki.Enable {
		client.waitGroup.Add(1)
		go client.run()
	}

	return &client
}

func (c *LokiLogger) Debug(job string, args ...interface{}) {

	if c.conf.console.Enable {
		c.zl.Debug(args...)
	}

	if c.conf.loki.Enable && logger.GetZapLevel(c.conf.loki.Verbosity) <= c.level {
		c.push(buildLabels(c.conf.lcfg.Service, job), makeEntry("", "Debug: ", args...))
	}
}

func (c *LokiLogger) Error(job string, args ...interface{}) {

	if c.conf.console.Enable {
		c.zl.Error(args...)
	}

	if c.conf.loki.Enable && logger.GetZapLevel(c.conf.loki.Verbosity) <= c.level {
		c.push(buildLabels(c.conf.lcfg.Service, job), makeEntry("", "Error: ", args...))
	}
}

func (c *LokiLogger) Warn(job string, args ...interface{}) {

	if c.conf.console.Enable {
		c.zl.Warn(args...)
	}

	if c.conf.loki.Enable && logger.GetZapLevel(c.conf.loki.Verbosity) <= c.level {
		c.push(buildLabels(c.conf.lcfg.Service, job), makeEntry("", "Warn: ", args...))
	}
}

func (c *LokiLogger) Info(job string, args ...interface{}) {

	if c.conf.console.Enable {
		c.zl.Info(args...)
	}

	if c.conf.loki.Enable && logger.GetZapLevel(c.conf.loki.Verbosity) <= c.level {
		c.push(buildLabels(c.conf.lcfg.Service, job), makeEntry("", "Info: ", args...))
	}
}

func (c *LokiLogger) Debugf(job, template string, args ...interface{}) {

	if c.conf.console.Enable {
		c.zl.Debugf(template, args...)
	}

	if c.conf.loki.Enable && logger.GetZapLevel(c.conf.loki.Verbosity) <= c.level {
		c.push(buildLabels(c.conf.lcfg.Service, job), makeEntry(template, "Debug: ", args...))
	}
}

func (c *LokiLogger) Errorf(job, template string, args ...interface{}) {

	if c.conf.console.Enable {
		c.zl.Errorf(template, args...)
	}

	if c.conf.loki.Enable && logger.GetZapLevel(c.conf.loki.Verbosity) <= c.level {
		c.push(buildLabels(c.conf.lcfg.Service, job), makeEntry(template, "Error: ", args...))
	}
}

func (c *LokiLogger) Warnf(job, template string, args ...interface{}) {

	if c.conf.console.Enable {
		c.zl.Warnf(template, args...)
	}

	if c.conf.loki.Enable && logger.GetZapLevel(c.conf.loki.Verbosity) <= c.level {
		c.push(buildLabels(c.conf.lcfg.Service, job), makeEntry(template, "Warn: ", args...))
	}
}

func (c *LokiLogger) Infof(job, template string, args ...interface{}) {

	if c.conf.console.Enable {
		c.zl.Infof(template, args...)
	}

	if c.conf.loki.Enable && logger.GetZapLevel(c.conf.loki.Verbosity) <= c.level {
		c.push(buildLabels(c.conf.lcfg.Service, job), makeEntry(template, "Info: ", args...))
	}
}

func (c *LokiLogger) push(labels string, entry *v1.Entry) {
	c.entries <- streamItem{
		labels: labels,
		entry:  entry,
	}
}

func (c *LokiLogger) run() {

	maxWait := time.NewTimer(timeout * time.Second)
	batch := make(map[string][]*v1.Entry, MIN_ENTRIES)

	defer func() {
		if len(batch) > 0 {
			c.process(batch)
		}
		c.waitGroup.Done()
	}()

	for {
		select {

		case <-c.done:
			return

		case entry := <-c.entries:

			batch[entry.labels] = append(batch[entry.labels], entry.entry)

			if len(batch) >= c.conf.lcfg.Batch.BatchSize {
				c.process(batch)
				batch = make(map[string][]*v1.Entry, MIN_ENTRIES)
				maxWait.Reset(timeout * time.Second)
			}

		case <-maxWait.C:

			if len(batch) > 0 {
				c.process(batch)
				batch = make(map[string][]*v1.Entry, MIN_ENTRIES)
			}
			maxWait.Reset(timeout * time.Second)
		}
	}
}

func (c *LokiLogger) process(entries map[string][]*v1.Entry) error {
	var streams []*v1.Stream

	for labels, arr := range entries {
		for _, v := range arr {
			streams = append(streams, &v1.Stream{
				Labels: labels,
				Entry:  v,
			})
		}
	}

	req := v1.PushRequest{
		Streams: streams,
	}

	pbuf, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	buf := snappy.Encode(nil, pbuf)

	resp, err := c.send(bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	if resp.code != 204 {
		return fmt.Errorf("invalid response code: %d", resp.code)
	}

	return nil
}

func (c *LokiLogger) send(buff *bytes.Buffer) (serverResp, error) {
	var out = serverResp{}

	req, err := http.NewRequest("POST", c.conf.lcfg.Url, buff)
	if err != nil {
		return out, err
	}

	req.Header.Set("Content-Type", c.conf.lcfg.Ctype)

	resp, err := c.http.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	out.code = resp.StatusCode
	out.body, err = io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (c *LokiLogger) Shutdown() {
	close(c.done)
	c.waitGroup.Wait()
}

func makeEntry(template, prefix string, args ...interface{}) *v1.Entry {
	now := time.Now().UnixNano()
	return &v1.Entry{
		Timestamp: &timestamppb.Timestamp{
			Seconds: now / int64(time.Second),
			Nanos:   int32(now % int64(time.Second)),
		},
		Line: fmt.Sprintf(prefix+template, args...),
	}
}

func buildLabels(service, job string) string {
	return "{service=\"" + service + "\",job=\"" + job + "\"}"
}
