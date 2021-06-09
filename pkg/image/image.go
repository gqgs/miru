package image

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pierrre/imageutil"
)

type Histogram struct {
	Red   [256]uint64
	Green [256]uint64
	Blue  [256]uint64

	once            sync.Once
	normalizedRed   [256]float64
	normalizedGreen [256]float64
	normalizedBlue  [256]float64
}

func (h *Histogram) normalize() {
	h.once.Do(func() {
		var sum uint64
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
	image, err := deserialize(b)
	if err != nil {
		return 0, "", err
	}

	return compare(i.Hist, image.Hist), image.String(), nil
}

// Implements the Stringer interface
func (i Image) String() string {
	return i.Filename
}

// Implements the BinaryMarshaler interface
func (img *Image) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer
	filenameBytes := []byte(img.Filename)
	filenameLen := len(filenameBytes)

	buffer.Grow(768*8 + filenameLen)
	binBuffer := make([]byte, 8)
	for i := 0; i < 256; i++ {
		binary.PutUvarint(binBuffer, img.Hist.Red[i])
		if _, err := buffer.Write(binBuffer); err != nil {
			return nil, err
		}
	}
	for i := 0; i < 256; i++ {
		binary.PutUvarint(binBuffer, img.Hist.Green[i])
		if _, err := buffer.Write(binBuffer); err != nil {
			return nil, err
		}
	}
	for i := 0; i < 256; i++ {
		binary.PutUvarint(binBuffer, img.Hist.Blue[i])
		if _, err := buffer.Write(binBuffer); err != nil {
			return nil, err
		}
	}
	if _, err := buffer.Write(filenameBytes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func deserialize(b []byte) (*Image, error) {
	img := new(Image)
	img.Hist = new(Histogram)
	binBuffer := make([]byte, 8)
	reader := bytes.NewReader(b)
	for i := 0; i < 256; i++ {
		if _, err := reader.Read(binBuffer); err != nil {
			return nil, err
		}
		value, _ := binary.Uvarint(binBuffer)
		img.Hist.Red[i] = value
	}
	for i := 0; i < 256; i++ {
		if _, err := reader.Read(binBuffer); err != nil {
			return nil, err
		}
		value, _ := binary.Uvarint(binBuffer)
		img.Hist.Green[i] = value
	}
	for i := 0; i < 256; i++ {
		if _, err := reader.Read(binBuffer); err != nil {
			return nil, err
		}
		value, _ := binary.Uvarint(binBuffer)
		img.Hist.Blue[i] = value
	}
	binBuffer = make([]byte, reader.Len())
	if _, err := reader.Read(binBuffer); err != nil {
		return nil, err
	}
	img.Filename = string(binBuffer)
	return img, nil
}

// Alternative Chi-Square
func compare(h1, h2 *Histogram) float64 {
	var result float64
	h1.normalize()
	h2.normalize()
	for i := 0; i < 256; i++ {
		if num := (h1.normalizedRed[i] + h2.normalizedRed[i]); num > 0 {
			result += (h1.normalizedRed[i] - h2.normalizedRed[i]) * (h1.normalizedRed[i] - h2.normalizedRed[i]) / num
		}
		if num := (h1.normalizedGreen[i] + h2.normalizedGreen[i]); num > 0 {
			result += (h1.normalizedGreen[i] - h2.normalizedGreen[i]) * (h1.normalizedGreen[i] - h2.normalizedGreen[i]) / num
		}
		if num := (h1.normalizedBlue[i] + h2.normalizedBlue[i]); num > 0 {
			result += (h1.normalizedBlue[i] - h2.normalizedBlue[i]) * (h1.normalizedBlue[i] - h2.normalizedBlue[i]) / num
		}
	}
	return math.Abs(2 * result)
}

func readFile(filename string) (io.ReadCloser, error) {
	if isURL(filename) {
		resp, err := http.Get(filename)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil

	}
	return os.Open(filename)
}

func Load(filename string) (*Image, error) {
	file, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reader io.Reader
	reader = bufio.NewReader(file)
	magicNumber, err := reader.(*bufio.Reader).Peek(2)
	if err != nil {
		return nil, err
	}

	if isJPEG := string(magicNumber) == string([]byte{0xFF, 0xD8}); isJPEG {
		hist, err := decodeJpeg(reader)
		if err == nil {
			return &Image{
				Filename: filename,
				Hist:     hist,
			}, nil
		}

		// Fallback to default decoder
		// The body has been consumed already
		// Therefore, we have to initilize it again
		// This shouldn't happen too often so it's cheaper than keeping a copy of the body
		reader, err = readFile(filename)
		if err != nil {
			return nil, err
		}
		//nolint:errcheck
		defer reader.(io.Closer).Close()
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	var hist Histogram
	var bounds = img.Bounds()
	atFunc := imageutil.NewAtFunc(img)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			red, green, blue, _ := atFunc(x, y)
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

func IsImage(filename string) bool {
	return strings.HasPrefix(mime.TypeByExtension(filepath.Ext(filename)), "image")
}
