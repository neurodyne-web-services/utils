package serdes

import (
	"bytes"
	"encoding/gob"
)

// Encode - encodes any interface into bytes with Gob.
func Encode[T any](data T) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	err := enc.Encode(data)

	return buff.Bytes(), err
}

// Encode - decodes Gob bytes into a value of a predefined data type.
func Decode[T any](data []byte) (T, error) {
	var out T

	buff := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buff)

	if err := dec.Decode(&out); err != nil {
		return out, err
	}

	return out, nil
}
