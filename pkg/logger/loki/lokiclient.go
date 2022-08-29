package loki

import (
	"bytes"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type serverResp struct {
	code int
	body []byte
}

type lokiConfig struct {
	url   string
	ctype string
}

type lokiClient struct {
	Client
	cconf ClientConfig
	lconf lokiConfig
	http  http.Client
	zl    *zap.Logger
}

// MakeLokiClient - factory for Loki client
func MakeLokiClient(cconf ClientConfig, lconf lokiConfig, zl *zap.Logger) lokiClient {
	return lokiClient{cconf: cconf, lconf: lconf, zl: zl}
}

func (c *lokiClient) Push(buff *bytes.Buffer) (serverResp, error) {
	var out = serverResp{}

	req, err := http.NewRequest("POST", c.lconf.url, buff)
	if err != nil {
		return out, err
	}

	req.Header.Set("Content-Type", c.lconf.ctype)

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
