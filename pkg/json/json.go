package json

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

// MarshalGzip - JSON encoder with gzip compression.
func MarshalGzip(data any) ([]byte, error) {
	buf := &bytes.Buffer{}
	g := gzip.NewWriter(buf)
	enc := json.NewEncoder(g)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	if err := g.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalGzip - JSON decoder with gzip compression.
func UnmarshalGzip(b []byte, dst any) error {
	r := bytes.NewReader(b)
	g, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(g)
	return dec.Decode(dst)
}
