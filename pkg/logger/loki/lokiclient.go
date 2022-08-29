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

type serverResp struct {
	code int
	body []byte
}

type streamItem struct {
	labels string
	entry  *v1.Entry
}

type lokiConfig struct {
	url     string
	ctype   string
	service string
	level   zapcore.Level
}

func MakeLokiConfig(url, ctype, service string, level zapcore.Level) lokiConfig {
	return lokiConfig{
		url:     url,
		ctype:   ctype,
		service: service,
		level:   level,
	}
}

type batchConfig struct {
	BatchEntriesNumber int
	BatchWait          time.Duration
}

type lokiClient struct {
	Client
	enableLoki    bool
	enableConsole bool
	batch         batchConfig
	conf          lokiConfig
	http          http.Client
	zl            *zap.Logger
	entries       chan streamItem
	quit          chan struct{}
	waitGroup     sync.WaitGroup
}

// MakeLokiLogger - factory for Loki client
func MakeLokiLogger(conf lokiConfig, zl *zap.Logger, enaConsole, enaLoki bool, batch batchConfig) *lokiClient {
	ch := make(chan streamItem, MAX_ENTRIES)

	client := lokiClient{
		enableConsole: enaConsole,
		enableLoki:    enaLoki,
		conf:          conf,
		zl:            zl,
		entries:       ch,
	}

	if enaLoki {
		client.waitGroup.Add(1)
		go client.run()
	}

	return &client
}

func (c *lokiClient) Debugf(job, template string, args ...interface{}) {

	if c.enableConsole {
		c.zl.Sugar().Debugf(template, args...)
	}

	if c.enableLoki && c.conf.level == zapcore.DebugLevel {
		labels := "{service=\"" + c.conf.service + "\",job=\"" + job + "\"}"
		c.push(labels, makeEntry(template, "Debug: ", args...))
	}
}

func (c *lokiClient) push(labels string, entry *v1.Entry) {
	c.entries <- streamItem{
		labels: labels,
		entry:  entry,
	}
}

func (c *lokiClient) run() {

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

		case <-c.quit:
			if batchSize > 0 {
				c.process(batch)
			}
			return

		case entry := <-c.entries:

			batch[entry.labels] = entry.entry
			batchSize++

			if batchSize >= c.batch.BatchEntriesNumber {
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

func (c *lokiClient) process(entries map[string]*v1.Entry) error {
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

func (c *lokiClient) send(buff *bytes.Buffer) (serverResp, error) {
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
