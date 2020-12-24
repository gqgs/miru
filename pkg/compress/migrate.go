package compress

type migrateCompressor struct {
	gzip Compressor
	zstd Compressor
}

func NewMigrateCompressor() migrateCompressor {
	return migrateCompressor{
		gzip: NewGzip(),
		zstd: NewZstdCompressor(),
	}
}

func (c migrateCompressor) Compress(b []byte) ([]byte, error) {
	return c.zstd.Compress(b)
}

func (c migrateCompressor) Decompress(b []byte) ([]byte, error) {
	return c.gzip.Decompress(b)
}
