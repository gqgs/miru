package image

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func Test_Load(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			"given a jpeg image",
			"./testdata/jpeg_image.jpg",
			"./testdata/hist_jpeg.golden",
		},
		{
			"given a png image",
			"./testdata/png_image.png",
			"./testdata/hist_png.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist, err := Load(tt.filename)
			if err != nil {
				t.Fatal(err)
			}

			file, err := os.Open(tt.want)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			for i := 0; i < 256; i++ {
				var red, green, blue uint64
				if _, err = fmt.Fscanf(file, "%d %d %d\n", &red, &green, &blue); err != nil {
					t.Fatal(err)
				}
				if hist.Hist.Red[i] != red {
					t.Errorf("wrong color for red[%d]: want %d got %d", i, red, hist.Hist.Red[i])
				}
				if hist.Hist.Green[i] != green {
					t.Errorf("wrong color for green[%d]: want %d got %d", i, green, hist.Hist.Green[i])
				}
				if hist.Hist.Blue[i] != blue {
					t.Errorf("wrong color for blue[%d]: want %d got %d", i, blue, hist.Hist.Blue[i])
				}
			}
		})
	}
}

func Test_compare(t *testing.T) {
	tests := []struct {
		name      string
		filename1 string
		filename2 string
		want      float64
	}{
		{
			"given a jpeg image",
			"./testdata/jpeg_image.jpg",
			"./testdata/jpeg_image.jpg",
			0.0,
		},
		{
			"given a png image",
			"./testdata/png_image.png",
			"./testdata/png_image.png",
			0.0,
		},
		{
			"given distinct images",
			"./testdata/jpeg_image.jpg",
			"./testdata/png_image.png",
			10.887659072875977,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img1, err := Load(tt.filename1)
			if err != nil {
				log.Fatalf("%s: %s", err, tt.filename1)
			}
			img2, err := Load(tt.filename2)
			if err != nil {
				log.Fatalf("%s: %s", err, tt.filename2)
			}
			if got := compare(img1.Hist, img2.Hist); got != tt.want {
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

func Test_IsImage(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			"empty string",
			"",
			false,
		},
		{
			"jpeg image",
			"image.jpg",
			true,
		},
		{
			"another jpeg image",
			"image.jpeg",
			true,
		},
		{
			"png image",
			"image.png",
			true,
		},
		{
			"webm image",
			"image.webp",
			true,
		},
		{
			"webm video",
			"image.webm",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsImage(tt.filename); got != tt.want {
				t.Errorf("isImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
