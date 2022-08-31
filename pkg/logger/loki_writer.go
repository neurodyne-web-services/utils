package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang/snappy"
	v1 "github.com/neurodyne-web-services/utils/pkg/logger/loki/genout/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	MIN_ENTRIES = 2
)

type serverResp struct {
	code int
	body []byte
}

type zapMsg struct {
	Level   string
	Time    time.Time
	Message string
}

type LokiConfig struct {
	Url     string
	Ctype   string
	Service string
	Batch   BatchConfig
}

func MakeLokiConfig(url, ctype, service string, batch BatchConfig) LokiConfig {
	return LokiConfig{
		Url:     url,
		Ctype:   ctype,
		Service: service,
		Batch:   batch,
	}
}

type streamItem struct {
	labels string
	entry  *v1.Entry
}

type BatchConfig struct {
	BatchSize       int
	BatchTimeoutSec int
}

func MakeBatchConfig(size, timeout int) BatchConfig {
	return BatchConfig{size, timeout}
}

type LokiSyncer struct {
	conf    LokiConfig
	http    http.Client
	entries chan streamItem
	done    chan struct{}
}

func MakeLokiSyncer(conf LokiConfig) *LokiSyncer {

	client := LokiSyncer{
		conf:    conf,
		done:    make(chan struct{}),
		entries: make(chan streamItem, conf.Batch.BatchSize),
	}

	go client.run()

	return &client
}

func (l LokiSyncer) Write(p []byte) (n int, err error) {

	var msg zapMsg

	err = json.Unmarshal(p, &msg)
	if err != nil {
		return 0, err
	}

	// fmt.Printf("Decoded level: %s, ts: %v, message: %s \n", msg.Level, msg.Time, msg.Message)
	// entry := makeEntry(msg.Level, msg.Message)
	// fmt.Printf("Entry: %v \n", entry)

	l.push(buildLabels(l.conf.Service, "test-job"), makeEntry(msg.Level, msg.Message))

	return 0, nil
}

func (l *LokiSyncer) Sync() error {
	fmt.Println(">>>> Loki Sync")
	return nil
}

func makeEntry(level, msg string) *v1.Entry {
	now := time.Now().UnixNano()
	return &v1.Entry{
		Timestamp: &timestamppb.Timestamp{
			Seconds: now / int64(time.Second),
			Nanos:   int32(now % int64(time.Second)),
		},
		Line: fmt.Sprintf(level, msg),
	}
}

func buildLabels(service, job string) string {
	return "{service=\"" + service + "\",job=\"" + job + "\"}"
}

func (c *LokiSyncer) push(labels string, entry *v1.Entry) {
	c.entries <- streamItem{
		labels: labels,
		entry:  entry,
	}
}

func (c *LokiSyncer) run() {

	maxWait := time.NewTimer(time.Duration(c.conf.Batch.BatchTimeoutSec) * time.Second)
	batch := make(map[string][]*v1.Entry, MIN_ENTRIES)

	defer func() {
		if len(batch) > 0 {
			c.process(batch)
		}
	}()

	for {
		select {

		case <-c.done:
			return

		case entry := <-c.entries:

			batch[entry.labels] = append(batch[entry.labels], entry.entry)

			if len(batch) >= c.conf.Batch.BatchSize {
				c.process(batch)
				batch = make(map[string][]*v1.Entry, MIN_ENTRIES)
				maxWait.Reset(time.Duration(c.conf.Batch.BatchTimeoutSec) * time.Second)
			}

		case <-maxWait.C:

			if len(batch) > 0 {
				c.process(batch)
				batch = make(map[string][]*v1.Entry, MIN_ENTRIES)
			}
			maxWait.Reset(time.Duration(c.conf.Batch.BatchTimeoutSec) * time.Second)
		}
	}
}

func (c *LokiSyncer) process(entries map[string][]*v1.Entry) error {
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

func (c *LokiSyncer) send(buff *bytes.Buffer) (serverResp, error) {
	var out = serverResp{}

	req, err := http.NewRequest("POST", c.conf.Url, buff)
	if err != nil {
		return out, err
	}

	req.Header.Set("Content-Type", c.conf.Ctype)

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

func (c *LokiSyncer) Shutdown() {
	close(c.done)
}
