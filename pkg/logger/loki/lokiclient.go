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
)

type LokiLogger struct {
	Client

	conf      LogConfig
	http      http.Client
	zl        *zap.SugaredLogger
	entries   chan streamItem
	done      chan struct{}
	waitGroup sync.WaitGroup
}

// MakeLokiLogger - factory for Loki client
func MakeLokiLogger(conf LogConfig, zl *zap.Logger) *LokiLogger {

	client := LokiLogger{
		conf:    conf,
		zl:      zl.Sugar(),
		done:    make(chan struct{}),
		entries: make(chan streamItem, conf.LokiConf.Batch.BatchSize),
	}

	if conf.Loki.Enable {
		client.waitGroup.Add(1)
		go client.run()
	}

	return &client
}

func (c *LokiLogger) Debug(job string, args ...interface{}) {

	if c.conf.Loki.Enable {
		c.zl.Debug(args...)
	}

	if c.conf.Loki.Enable && zapcore.DebugLevel >= logger.GetZapLevel(c.conf.Loki.Level) {
		c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry("", "Debug: ", args...))
	}
}

func (c *LokiLogger) Fatal(job string, args ...interface{}) {

	c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry("", "Fatal: ", args...))
	c.zl.Fatal(args...)
}

func (c *LokiLogger) Error(job string, args ...interface{}) {

	if c.conf.Loki.Enable {
		c.zl.Error(args...)
	}

	if c.conf.Loki.Enable && zapcore.ErrorLevel >= logger.GetZapLevel(c.conf.Loki.Level) {
		c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry("", "Error: ", args...))
	}
}

func (c *LokiLogger) Warn(job string, args ...interface{}) {

	if c.conf.Loki.Enable {
		c.zl.Warn(args...)
	}

	if c.conf.Loki.Enable && zapcore.WarnLevel >= logger.GetZapLevel(c.conf.Loki.Level) {
		c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry("", "Warn: ", args...))
	}
}

func (c *LokiLogger) Info(job string, args ...interface{}) {

	if c.conf.Loki.Enable {
		c.zl.Info(args...)
	}

	if c.conf.Loki.Enable && zapcore.InfoLevel >= logger.GetZapLevel(c.conf.Loki.Level) {
		c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry("", "Info: ", args...))
	}
}

func (c *LokiLogger) Debugf(job, template string, args ...interface{}) {

	if c.conf.Loki.Enable {
		c.zl.Debugf(template, args...)
	}

	if c.conf.Loki.Enable && zapcore.DebugLevel >= logger.GetZapLevel(c.conf.Loki.Level) {
		c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry(template, "Debug: ", args...))
	}
}

func (c *LokiLogger) Fatalf(job, template string, args ...interface{}) {

	c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry(template, "Fatal: ", args...))
	c.zl.Fatalf(template, args...)
}

func (c *LokiLogger) Errorf(job, template string, args ...interface{}) {

	if c.conf.Loki.Enable {
		c.zl.Errorf(template, args...)
	}

	if c.conf.Loki.Enable && zapcore.ErrorLevel >= logger.GetZapLevel(c.conf.Loki.Level) {
		c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry(template, "Error: ", args...))
	}
}

func (c *LokiLogger) Warnf(job, template string, args ...interface{}) {

	if c.conf.Loki.Enable {
		c.zl.Warnf(template, args...)
	}

	if c.conf.Loki.Enable && zapcore.WarnLevel >= logger.GetZapLevel(c.conf.Loki.Level) {
		c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry(template, "Warn: ", args...))
	}
}

func (c *LokiLogger) Infof(job, template string, args ...interface{}) {

	if c.conf.Loki.Enable {
		c.zl.Infof(template, args...)
	}

	if c.conf.Loki.Enable && zapcore.InfoLevel >= logger.GetZapLevel(c.conf.Loki.Level) {
		c.push(buildLabels(c.conf.LokiConf.Service, job), makeEntry(template, "Info: ", args...))
	}
}

func (c *LokiLogger) push(labels string, entry *v1.Entry) {
	c.entries <- streamItem{
		labels: labels,
		entry:  entry,
	}
}

func (c *LokiLogger) run() {

	maxWait := time.NewTimer(time.Duration(c.conf.LokiConf.Batch.BatchTimeoutSec) * time.Second)
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

			if len(batch) >= c.conf.LokiConf.Batch.BatchSize {
				c.process(batch)
				batch = make(map[string][]*v1.Entry, MIN_ENTRIES)
				maxWait.Reset(time.Duration(c.conf.LokiConf.Batch.BatchTimeoutSec) * time.Second)
			}

		case <-maxWait.C:

			if len(batch) > 0 {
				c.process(batch)
				batch = make(map[string][]*v1.Entry, MIN_ENTRIES)
			}
			maxWait.Reset(time.Duration(c.conf.LokiConf.Batch.BatchTimeoutSec) * time.Second)
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

	req, err := http.NewRequest("POST", c.conf.LokiConf.Url, buff)
	if err != nil {
		return out, err
	}

	req.Header.Set("Content-Type", c.conf.LokiConf.Ctype)

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
