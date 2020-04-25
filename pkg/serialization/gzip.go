package serialization

import (
	"bytes"
	gzip "compress/gzip"
	"encoding/json"
	"io/ioutil"
)

type gzipSerializer struct{}

func NewGzipSerializer() *gzipSerializer {
	return &gzipSerializer{}
}

func (g gzipSerializer) Marshal(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(buf)
	if _, err := gzipWriter.Write(b); err != nil {
		return nil, err
	}
	if err = gzipWriter.Close(); err != nil {
		return nil, err
	}
	compressed, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	return compressed, nil
}

func (g gzipSerializer) Unmarshal(b []byte, v interface{}) error {
	reader := bytes.NewReader(b)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	decompressed, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return err
	}
	if err := gzipReader.Close(); err != nil {
		return err
	}
	return json.Unmarshal(decompressed, v)
}
