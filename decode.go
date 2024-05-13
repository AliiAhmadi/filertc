package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
)

func decode(i string, v interface{}) error {
	b, err := base64.StdEncoding.DecodeString(i)
	if err != nil {
		return err
	}

	z, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer z.Close()

	con, err := io.ReadAll(z)
	if err != nil {
		return err
	}

	return json.Unmarshal(con, v)
}
