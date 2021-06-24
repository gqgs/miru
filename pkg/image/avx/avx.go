package avx

import (
	"unsafe"
)

//go:noescape
func __Z16compare_hist_avxPfS_S_S_S_S_S_(h1r, h2r, h1g, h2g, h1b, h2b, result unsafe.Pointer)

func CompareHist(h1r, h2r, h1g, h2g, h1b, h2b [256]float32) float32 {
	var result [8]float32
	__Z16compare_hist_avxPfS_S_S_S_S_S_(
		unsafe.Pointer(&h1r[0]), unsafe.Pointer(&h2r[0]),
		unsafe.Pointer(&h1g[0]), unsafe.Pointer(&h2g[0]),
		unsafe.Pointer(&h1b[0]), unsafe.Pointer(&h2b[0]),
		unsafe.Pointer(&result[0]))
	return result[0] + result[1] + result[2] + result[3] +
		result[4] + result[5] + result[6] + result[7]
}
