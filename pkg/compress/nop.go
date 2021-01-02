package compress

func newNop() *nopCompressor {
	return &nopCompressor{}
}

type nopCompressor struct{}

func (c nopCompressor) Compress(b []byte) ([]byte, error) {
	return b, nil
}

func (c nopCompressor) Decompress(b []byte) ([]byte, error) {
	return b, nil
}
