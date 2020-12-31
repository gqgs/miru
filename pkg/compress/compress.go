package compress

import "errors"

type Compressor interface {
	Compress(b []byte) ([]byte, error)
	Decompress(b []byte) ([]byte, error)
}

func NewCompressor(name string) (Compressor, error) {
	switch name {
	case "zstd":
		return newZstd(), nil
	case "gzip":
		return newGzip(), nil
	}
	return nil, errors.New("invalid compressor")
}
