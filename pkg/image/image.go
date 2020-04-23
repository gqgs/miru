package image

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

type Histogram struct {
	Red   [256]uint
	Green [256]uint
	Blue  [256]uint
}

// TODO: replace above
type NormalizedHistogram struct {
	Red   [256]float64
	Green [256]float64
	Blue  [256]float64
}

type Image struct {
	Filename string
	Hist     Histogram
}

func (i Image) String() string {
	return i.Filename
}

func Compare(img1, img2 *Image) float64 {
	return compare(img1.Hist, img2.Hist)
}

// Alternative Chi-Square
func compare(h1, h2 Histogram) float64 {
	var result float64
	nh1 := normalize(h1)
	nh2 := normalize(h2)
	for i := 0; i < 256; i++ {
		if num := (nh1.Red[i] + nh2.Red[i]); num > 0 {
			result += math.Pow(nh1.Red[i]-nh2.Red[i], 2) / num
		}
		if num := (nh1.Green[i] + nh2.Green[i]); num > 0 {
			result += math.Pow(nh1.Green[i]-nh2.Green[i], 2) / num
		}
		if num := (nh1.Blue[i] + nh2.Blue[i]); num > 0 {
			result += math.Pow(nh1.Blue[i]-nh2.Blue[i], 2) / num
		}
	}
	return math.Abs(2 * result)
}

func normalize(h Histogram) NormalizedHistogram {
	var normalized NormalizedHistogram
	var sum uint
	for i := 0; i < 256; i++ {
		sum += h.Red[i] * h.Red[i]
		sum += h.Green[i] * h.Green[i]
		sum += h.Blue[i] * h.Blue[i]
	}
	norm := math.Sqrt(float64(sum))
	for i := 0; i < 256; i++ {
		normalized.Red[i] = float64(h.Red[i]) / norm
		normalized.Green[i] = float64(h.Green[i]) / norm
		normalized.Blue[i] = float64(h.Blue[i]) / norm
	}
	return normalized
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
		Hist:     hist,
	}, nil
}
