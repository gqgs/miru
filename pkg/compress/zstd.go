package compress

import "github.com/DataDog/zstd"

func NewZstdCompressor() *zstdCompressor {
	return &zstdCompressor{}
}

type zstdCompressor struct{}

func (z zstdCompressor) Compress(b []byte) ([]byte, error) {
	return zstd.Compress(nil, b)
}

func (z zstdCompressor) Decompress(b []byte) ([]byte, error) {
	return zstd.Decompress(nil, b)
}
