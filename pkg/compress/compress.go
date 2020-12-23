package compress

import "errors"

type Compressor interface {
	Compress(b []byte) ([]byte, error)
	Decompress(b []byte) ([]byte, error)
}

func NewCompressor(name string) (Compressor, error) {
	switch name {
	case "zstd":
		return NewZstdCompressor(), nil
	case "gzip":
		return NewGzip(), nil
	}
	return nil, errors.New("invalid compressor")
}
