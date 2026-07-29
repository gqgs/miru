package avx

import "unsafe"

//go:noescape
func compareHistAVX512(h1r, h2r, h1g, h2g, h1b, h2b, result unsafe.Pointer)

// CompareHistAVX512 compares three 256-bin histogram channels using AVX-512F.
// The caller must check that AVX-512F is available before calling it.
func CompareHistAVX512(h1r, h2r, h1g, h2g, h1b, h2b [256]float32) float32 {
	var result [16]float32
	compareHistAVX512(
		unsafe.Pointer(&h1r[0]), unsafe.Pointer(&h2r[0]),
		unsafe.Pointer(&h1g[0]), unsafe.Pointer(&h2g[0]),
		unsafe.Pointer(&h1b[0]), unsafe.Pointer(&h2b[0]),
		unsafe.Pointer(&result[0]),
	)

	var sum float32
	for _, value := range result {
		sum += value
	}
	return sum
}
