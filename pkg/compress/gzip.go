package compress

import (
	"bytes"
	gzip "compress/gzip"
	"io/ioutil"
)

type gzipCompressor struct{}

func NewGzip() *gzipCompressor {
	return &gzipCompressor{}
}

func (g gzipCompressor) Compress(b []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(buf)
	if _, err := gzipWriter.Write(b); err != nil {
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}
	compressed, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	return compressed, nil
}

func (g gzipCompressor) Decompress(b []byte) ([]byte, error) {
	reader := bytes.NewReader(b)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	decompressed, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	if err := gzipReader.Close(); err != nil {
		return nil, err
	}

	return decompressed, nil
}
