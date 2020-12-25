package compress

import "github.com/valyala/gozstd"

func NewZstdCompressor() *zstdCompressor {
	return &zstdCompressor{}
}

type zstdCompressor struct{}

func (z zstdCompressor) Compress(b []byte) ([]byte, error) {
	return gozstd.Compress(nil, b), nil
}

func (z zstdCompressor) Decompress(b []byte) ([]byte, error) {
	return gozstd.Decompress(nil, b)
}
