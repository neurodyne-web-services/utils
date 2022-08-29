package loki

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

type lokiClient struct {
	Client
	conf          lokiConfig
	http          http.Client
	zl            *zap.Logger
	enableLoki    bool
	enableConsole bool
	streams       chan []*v1.Stream
}

// MakeLokiLogger - factory for Loki client
func MakeLokiLogger(conf lokiConfig, zl *zap.Logger, enaConsole, enaLoki bool) lokiClient {
	return lokiClient{conf: conf, zl: zl, enableConsole: enaConsole, enableLoki: enaLoki}
}

func (c lokiClient) Debugf(job, template string, args ...interface{}) {

	if c.enableConsole {
		c.zl.Sugar().Debugf(template, args...)
	}

	if c.enableLoki && c.conf.level == zapcore.DebugLevel {
		tmp := makeEntry(template, "Debug: ", args...)

		labels := "{service=\"" + c.conf.service + "\",job=\"" + job + "\"}"
		c.Process(labels, tmp)
	}
}

func (c lokiClient) Push(buff *bytes.Buffer) (serverResp, error) {
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

func (c lokiClient) Process(labels string, entry *v1.Entry) error {
	var streams []*v1.Stream

	streams = append(streams, &v1.Stream{
		Labels:  labels,
		Entries: []*v1.Entry{entry},
	})

	req := v1.PushRequest{
		Streams: streams,
	}

	pbuf, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	buf := snappy.Encode(nil, pbuf)

	resp, err := c.Push(bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	if resp.code != 204 {
		return fmt.Errorf("invalid response code: %d", resp.code)
	}

	return nil
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
