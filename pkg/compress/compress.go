package compress

type Compressor interface {
	Compress(b []byte) ([]byte, error)
	Decompress(b []byte) ([]byte, error)
}
