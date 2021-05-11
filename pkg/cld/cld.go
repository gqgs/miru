package cld

import (
	"encoding/json"
	"image/jpeg"
	"os"

	"github.com/gqgs/mpeg7cld"
)

type CLD struct {
	Filename string
	CLD      [64]mpeg7cld.YCbCr
}

func (c *CLD) Compare(b []byte) (result float64, comparedElement string, err error) {
	var cld CLD
	if err := json.Unmarshal(b, &cld); err != nil {
		return 0, "", nil
	}
	return mpeg7cld.Compare(c.CLD, cld.CLD), cld.Filename, nil
}

func (c *CLD) String() string {
	return c.Filename
}

func (c *CLD) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func Load(filename string) (*CLD, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return &CLD{
		Filename: filename,
		CLD:      mpeg7cld.CLD(img),
	}, nil
}
