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

type LokiSyncer struct {
	conf LokiConfig
	http http.Client

	streams      []*v1.Stream
	entries      map[string][]*v1.Entry
	totalEntries int8
}

func MakeLokiSyncer(conf LokiConfig) *LokiSyncer {

	client := LokiSyncer{
		conf:    conf,
		entries: make(map[string][]*v1.Entry),
	}

	return &client
}

func (l LokiSyncer) Write(p []byte) (n int, err error) {

	var msg zapMsg

	err = json.Unmarshal(p, &msg)
	if err != nil {
		return 0, err
	}

	l.process(streamItem{
		labels: buildLabels(l.conf.Service, "test-job"),
		entry:  makeEntry(msg.Level, msg.Message),
	})

	return 0, nil
}

func (l *LokiSyncer) Sync() error {
	fmt.Println(">>>> Loki Sync Started")
	return l.procStreams()
}

func (l *LokiSyncer) process(item streamItem) error {

	if item.entry != nil {
		l.entries[item.labels] = append(l.entries[item.labels], item.entry)
		l.totalEntries++
	}

	if l.totalEntries >= l.conf.Batch.BatchSize {

		for labels, arr := range l.entries {
			for _, v := range arr {
				l.streams = append(l.streams, &v1.Stream{
					Labels: labels,
					Entry:  v,
				})
			}
		}
		return l.procStreams()
	}
	return nil
}

func (l *LokiSyncer) procStreams() error {
	req := v1.PushRequest{
		Streams: l.streams,
	}

	pbuf, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	buf := snappy.Encode(nil, pbuf)

	resp, err := l.send(bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	if resp.code != 204 {
		return fmt.Errorf("invalid response code: %d", resp.code)
	}

	return nil
}

func (l *LokiSyncer) send(buff *bytes.Buffer) (serverResp, error) {
	var out = serverResp{}

	req, err := http.NewRequest("POST", l.conf.Url, buff)
	if err != nil {
		return out, err
	}

	req.Header.Set("Content-Type", l.conf.Ctype)

	resp, err := l.http.Do(req)
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

func buildLabels(service, job string) string {
	return "{service=\"" + service + "\",job=\"" + job + "\"}"
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
