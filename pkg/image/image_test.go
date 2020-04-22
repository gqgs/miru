package image

import (
	"log"
	"testing"
)

func TestCompare(t *testing.T) {
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
			if got := Compare(img, img); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
