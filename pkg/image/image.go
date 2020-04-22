package image

import (
	"image/jpeg"
	"math"
	"os"
)

type Histogram struct {
	Red   [256]uint
	Green [256]uint
	Blue  [256]uint
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
	for i := 0; i < 256; i++ {
		if num := float64(h1.Red[i]) + float64(h2.Red[i]); num > 0 {
			result += (float64(h1.Red[i]) - float64(h2.Red[i])) / num
		}
		if num := float64(h1.Green[i]) + float64(h2.Green[i]); num > 0 {
			result += (float64(h1.Green[i]) - float64(h2.Green[i])) / num
		}
		if num := float64(h1.Blue[i]) + float64(h2.Blue[i]); num > 0 {
			result += (float64(h1.Blue[i]) - float64(h2.Blue[i])) / num
		}
	}
	return math.Abs(2*result) / 3
}

func Load(filename string) (*Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	var hist Histogram
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
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
