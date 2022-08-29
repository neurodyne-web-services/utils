package loki

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type serverResp struct {
	code int
	body []byte
}

type lokiConfig struct {
	url   string
	ctype string
	level zapcore.Level
}

type lokiClient struct {
	Client
	conf lokiConfig
	http http.Client
	zl   *zap.Logger
}

// MakeLokiClient - factory for Loki client
func MakeLokiClient(conf lokiConfig, zl *zap.Logger) lokiClient {
	return lokiClient{conf: conf, zl: zl}
}

func (c lokiClient) Debugf(format string, args ...interface{}) {
	c.zl.Sugar().Debugf(format, args)

	if c.conf.level == zapcore.DebugLevel {
		fmt.Println("Loki Debug called")
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
