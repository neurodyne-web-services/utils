package loki

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/golang/snappy"
	v1 "github.com/neurodyne-web-services/utils/pkg/logger/loki/genout/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	MAX_ENTRIES = 8
)

type LokiLogger struct {
	Client

	batch     batchConfig
	conf      lokiConfig
	http      http.Client
	zl        *zap.SugaredLogger
	entries   chan streamItem
	done      chan struct{}
	waitGroup sync.WaitGroup
}

// MakeLokiLogger - factory for Loki client
func MakeLokiLogger(conf lokiConfig, zl *zap.Logger, batch batchConfig) *LokiLogger {

	client := LokiLogger{
		conf:    conf,
		zl:      zl.Sugar(),
		entries: make(chan streamItem, MAX_ENTRIES),
		done:    make(chan struct{}),
	}

	if conf.enableLoki {
		client.waitGroup.Add(1)
		go client.run()
	}

	return &client
}

func (c *LokiLogger) Debugf(job, template string, args ...interface{}) {

	if c.conf.enableConsole {
		c.zl.Debugf(template, args...)
	}

	if c.conf.enableLoki && c.conf.level == zapcore.DebugLevel {
		c.push(buildLabels(c.conf.service, job), makeEntry(template, "Debug: ", args...))
	}
}

func (c *LokiLogger) Errorf(job, template string, args ...interface{}) {

	if c.conf.enableConsole {
		c.zl.Errorf(template, args...)
	}

	if c.conf.enableLoki && c.conf.level <= zapcore.ErrorLevel {
		c.push(buildLabels(c.conf.service, job), makeEntry(template, "Error: ", args...))
	}
}

func (c *LokiLogger) Warnf(job, template string, args ...interface{}) {

	if c.conf.enableConsole {
		c.zl.Warnf(template, args...)
	}

	if c.conf.enableLoki && c.conf.level <= zapcore.WarnLevel {
		c.push(buildLabels(c.conf.service, job), makeEntry(template, "Warn: ", args...))
	}
}

func (c *LokiLogger) Infof(job, template string, args ...interface{}) {

	if c.conf.enableConsole {
		c.zl.Infof(template, args...)
	}

	if c.conf.enableLoki && c.conf.level <= zapcore.InfoLevel {
		c.push(buildLabels(c.conf.service, job), makeEntry(template, "Info: ", args...))
	}
}

func (c *LokiLogger) push(labels string, entry *v1.Entry) {
	c.entries <- streamItem{
		labels: labels,
		entry:  entry,
	}
}

func (c *LokiLogger) run() {

	batch := make(map[string]*v1.Entry)

	batchSize := 0
	maxWait := time.NewTimer(c.batch.BatchWait)

	defer func() {
		if batchSize > 0 {
			c.process(batch)
		}
		c.waitGroup.Done()
	}()

	for {
		select {

		case <-c.done:
			if batchSize > 0 {
				c.process(batch)
			}
			return

		case entry := <-c.entries:

			batch[entry.labels] = entry.entry
			batchSize++

			if batchSize >= c.batch.BatchSize {
				c.process(batch)
				batch = make(map[string]*v1.Entry)
				batchSize = 0
				maxWait.Reset(c.batch.BatchWait)
			}

		case <-maxWait.C:

			if batchSize > 0 {
				c.process(batch)
				batch = make(map[string]*v1.Entry)
				batchSize = 0
			}
			maxWait.Reset(c.batch.BatchWait)
		}
	}
}

func (c *LokiLogger) process(entries map[string]*v1.Entry) error {
	var streams []*v1.Stream

	for key, v := range entries {
		streams = append(streams, &v1.Stream{
			Labels: key,
			Entry:  v,
		})
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

	req, err := http.NewRequest("POST", c.conf.url, buff)
	if err != nil {
		return out, err
	}

	req.Header.Set("Content-Type", c.conf.ctype)

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

func makeEntry(format, prefix string, args ...interface{}) *v1.Entry {
	now := time.Now().UnixNano()
	return &v1.Entry{
		Timestamp: &timestamppb.Timestamp{
			Seconds: now / int64(time.Second),
			Nanos:   int32(now % int64(time.Second)),
		},
		Line: fmt.Sprintf(prefix+format, args...),
	}
}

func buildLabels(service, job string) string {
	return "{service=\"" + service + "\",job=\"" + job + "\"}"
}
