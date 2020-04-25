package serialize

import (
	"bytes"
	libZgip "compress/gzip"
	"encoding/json"
	"io/ioutil"
)

type gzip struct{}

func NewGzip() *gzip {
	return &gzip{}
}

func (g gzip) Marshal(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	zipWriter := libZgip.NewWriter(buf)
	if _, err := zipWriter.Write(b); err != nil {
		return nil, err
	}
	if err = zipWriter.Close(); err != nil {
		return nil, err
	}
	compressed, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	return compressed, nil
}

func (g gzip) Unmarshal(b []byte, v interface{}) error {
	reader := bytes.NewReader(b)
	zipReader, err := libZgip.NewReader(reader)
	if err != nil {
		return err
	}
	decompressed, err := ioutil.ReadAll(zipReader)
	if err != nil {
		return err
	}
	if err := zipReader.Close(); err != nil {
		return err
	}
	return json.Unmarshal(decompressed, v)
}
