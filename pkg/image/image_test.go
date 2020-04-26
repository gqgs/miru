package image

import (
	"log"
	"testing"
)

func Test_compare(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     float64
	}{
		{
			"given a jpeg image",
			"./testdata/jpeg_image.jpg",
			0.0,
		},
		{
			"given a png image",
			"./testdata/png_image.png",
			0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := Load(tt.filename)
			if err != nil {
				log.Fatalf("%s: %s", err, tt.filename)
			}
			if got := compare(img.Hist, img.Hist); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_Load(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Load("./testdata/jpeg_image.jpg")
	}
}

func Benchmark_ImageCompare(b *testing.B) {
	img, err := Load("./testdata/jpeg_image.jpg")
	if err != nil {
		b.Fatal("err", err)
	}
	data, err := img.MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = img.Compare(data)
	}
}

func Benchmark_compare(b *testing.B) {
	img, err := Load("./testdata/jpeg_image.jpg")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = compare(img.Hist, img.Hist)
	}
}
