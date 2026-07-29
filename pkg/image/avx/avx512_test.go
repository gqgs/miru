package avx

import (
	"math"
	"testing"

	"golang.org/x/sys/cpu"
)

func TestCompareHistAVX512(t *testing.T) {
	if !cpu.X86.HasAVX512F {
		t.Skip("AVX-512F is not supported")
	}

	var h1r, h2r, h1g, h2g, h1b, h2b [256]float32
	for i := range h1r {
		// Leave some bins at zero to exercise masked division.
		if i%7 != 0 {
			h1r[i] = float32(i%13) / 17
			h2r[i] = float32(i%11) / 19
		}
		if i%5 != 0 {
			h1g[i] = float32(i%17) / 23
			h2g[i] = float32(i%19) / 29
		}
		if i%3 != 0 {
			h1b[i] = float32(i%23) / 31
			h2b[i] = float32(i%29) / 37
		}
	}

	want := compareHistScalar(h1r, h2r, h1g, h2g, h1b, h2b)
	got := CompareHistAVX512(h1r, h2r, h1g, h2g, h1b, h2b)
	if diff := math.Abs(float64(got - want)); diff > 1e-3 {
		t.Fatalf("CompareHistAVX512() = %v, want %v (difference %v)", got, want, diff)
	}
}

func compareHistScalar(channels ...[256]float32) float32 {
	var result float32
	for pair := 0; pair < len(channels); pair += 2 {
		for i := range channels[pair] {
			a, b := channels[pair][i], channels[pair+1][i]
			if denominator := a + b; denominator > 0 {
				result += (a - b) * (a - b) / denominator
			}
		}
	}
	return result
}
