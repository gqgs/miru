package image

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"sync"
)

type Histogram struct {
	Red   [256]uint
	Green [256]uint
	Blue  [256]uint

	once            sync.Once
	NormalizedRed   [256]float64
	NormalizedGreen [256]float64
	NormalizedBlue  [256]float64
}

func (h *Histogram) Normalize() {
	h.once.Do(func() {
		var sum uint
		for i := 0; i < 256; i++ {
			sum += h.Red[i] * h.Red[i]
			sum += h.Green[i] * h.Green[i]
			sum += h.Blue[i] * h.Blue[i]
		}
		norm := math.Sqrt(float64(sum))
		for i := 0; i < 256; i++ {
			h.NormalizedRed[i] = float64(h.Red[i]) / norm
			h.NormalizedGreen[i] = float64(h.Green[i]) / norm
			h.NormalizedBlue[i] = float64(h.Blue[i]) / norm
		}
	})
}

type Image struct {
	Filename string
	Hist     *Histogram
}

func (i Image) String() string {
	return i.Filename
}

func Compare(img1, img2 *Image) float64 {
	return compare(img1.Hist, img2.Hist)
}

// Alternative Chi-Square
func compare(h1, h2 *Histogram) float64 {
	var result float64
	h1.Normalize()
	h2.Normalize()
	for i := 0; i < 256; i++ {
		if num := (h1.NormalizedRed[i] + h2.NormalizedRed[i]); num > 0 {
			result += math.Pow(h1.NormalizedRed[i]-h2.NormalizedRed[i], 2) / num
		}
		if num := (h1.NormalizedGreen[i] + h2.NormalizedGreen[i]); num > 0 {
			result += math.Pow(h1.NormalizedGreen[i]-h2.NormalizedGreen[i], 2) / num
		}
		if num := (h1.NormalizedBlue[i] + h2.NormalizedBlue[i]); num > 0 {
			result += math.Pow(h1.NormalizedBlue[i]-h2.NormalizedBlue[i], 2) / num
		}
	}
	return math.Abs(2 * result)
}

func Load(filename string) (*Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	var hist Histogram
	var bounds = img.Bounds()
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			red, green, blue, _ := img.At(x, y).RGBA()
			hist.Red[red>>8]++
			hist.Green[green>>8]++
			hist.Blue[blue>>8]++
		}
	}

	return &Image{
		Filename: filename,
		Hist:     &hist,
	}, nil
}
