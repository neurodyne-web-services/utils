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

	streams []*v1.Stream
	entries map[string][]*v1.Entry
}

func MakeLokiSyncer(conf LokiConfig) *LokiSyncer {
	return &LokiSyncer{
		conf:    conf,
		entries: make(map[string][]*v1.Entry, MIN_ENTRIES),
	}
}

var labels string

func (l LokiSyncer) Write(p []byte) (n int, err error) {

	var msg zapMsg

	err = json.Unmarshal(p, &msg)
	if err != nil {
		return 0, err
	}

	labels = buildLabels(l.conf.Service, "myJOB")
	l.entries[labels] = append(l.entries[labels], makeEntry(msg.Level, msg.Caller, msg.Message))

	// Process a batch
	if len(l.entries[labels]) >= l.conf.BatchSize {
		if err = l.process(); err != nil {
			return 0, err
		}
	}

	if len(l.streams) > 0 {
		if err = l.procStreams(); err != nil {
			return 0, err
		}
	}

	return 0, nil
}

func (l *LokiSyncer) Sync() error {
	return l.process()
}

func (l *LokiSyncer) process() error {

	for labels, arr := range l.entries {
		l.streams = append(l.streams, &v1.Stream{
			Labels:  labels,
			Entries: arr,
		},
		)
	}

	l.entries = make(map[string][]*v1.Entry, MIN_ENTRIES)
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

	l.streams = make([]*v1.Stream, MIN_ENTRIES)

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

func makeEntry(level, caller, msg string) *v1.Entry {
	now := time.Now().UnixNano()
	return &v1.Entry{
		Timestamp: &timestamppb.Timestamp{
			Seconds: now / int64(time.Second),
			Nanos:   int32(now % int64(time.Second)),
		},
		// Line: fmt.Sprintf(level, msg),
		Line: level + ": " + caller + " " + msg,
	}
}
