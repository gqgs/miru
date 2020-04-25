package image

import (
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Histogram struct {
	Red   [256]uint
	Green [256]uint
	Blue  [256]uint

	once            sync.Once
	normalizedRed   [256]float64
	normalizedGreen [256]float64
	normalizedBlue  [256]float64
}

func (h *Histogram) normalize() {
	h.once.Do(func() {
		var sum uint
		for i := 0; i < 256; i++ {
			sum += h.Red[i] * h.Red[i]
			sum += h.Green[i] * h.Green[i]
			sum += h.Blue[i] * h.Blue[i]
		}
		norm := math.Sqrt(float64(sum))
		for i := 0; i < 256; i++ {
			h.normalizedRed[i] = float64(h.Red[i]) / norm
			h.normalizedGreen[i] = float64(h.Green[i]) / norm
			h.normalizedBlue[i] = float64(h.Blue[i]) / norm
		}
	})
}

type Image struct {
	Filename string
	Hist     *Histogram
}

// Implements the Comparer interface
func (i *Image) Compare(b []byte) (result float64, comparedElement string, err error) {
	image := new(Image)
	if err := json.Unmarshal(b, image); err != nil {
		return 0, "", err
	}
	return compare(i.Hist, image.Hist), image.String(), nil
}

func (i Image) String() string {
	return i.Filename
}

// Alternative Chi-Square
func compare(h1, h2 *Histogram) float64 {
	var result float64
	h1.normalize()
	h2.normalize()
	for i := 0; i < 256; i++ {
		if num := (h1.normalizedRed[i] + h2.normalizedRed[i]); num > 0 {
			result += math.Pow(h1.normalizedRed[i]-h2.normalizedRed[i], 2) / num
		}
		if num := (h1.normalizedGreen[i] + h2.normalizedGreen[i]); num > 0 {
			result += math.Pow(h1.normalizedGreen[i]-h2.normalizedGreen[i], 2) / num
		}
		if num := (h1.normalizedBlue[i] + h2.normalizedBlue[i]); num > 0 {
			result += math.Pow(h1.normalizedBlue[i]-h2.normalizedBlue[i], 2) / num
		}
	}
	return math.Abs(2 * result)
}

func Load(filename string) (*Image, error) {
	var file io.ReadCloser
	var err error
	if isURL(filename) {
		resp, err := http.Get(filename)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		file = resp.Body

	} else {
		file, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}

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

func isURL(s string) bool {
	return strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "http://")
}
