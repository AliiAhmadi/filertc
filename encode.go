package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
)

func encode(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	var zb bytes.Buffer
	z, err := gzip.NewWriterLevel(&zb, gzip.BestCompression)
	if err != nil {
		return "", err
	}

	if _, err := z.Write(b); err != nil {
		return "", err
	}

	if err := z.Flush(); err != nil {
		return "", err
	}

	if err := z.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(zb.Bytes()), nil
}
