package compress

func newNop() *nopCompressor {
	return &nopCompressor{}
}

type nopCompressor struct{}

func (nopCompressor) Compress(b []byte) ([]byte, error) {
	return b, nil
}

func (nopCompressor) Decompress(b []byte) ([]byte, error) {
	return b, nil
}
