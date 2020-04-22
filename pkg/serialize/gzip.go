package serialize

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/vmihailenco/msgpack/v4"
)

type Gzip struct{}

func NewGzip() *Gzip {
	return &Gzip{}
}

func (g Gzip) Marshal(v interface{}) ([]byte, error) {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	zipWriter := gzip.NewWriter(buf)
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

func (g Gzip) Unmarshal(b []byte, v interface{}) error {
	reader := bytes.NewReader(b)
	zipReader, err := gzip.NewReader(reader)
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
	return msgpack.Unmarshal(decompressed, v)
}
