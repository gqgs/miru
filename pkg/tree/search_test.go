package tree

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/image"
	"github.com/gqgs/miru/pkg/storage"
)

func Benchmark_Search(b *testing.B) {
	compressor, err := compress.NewCompressor("nop")
	if err != nil {
		b.Fatal(err)
	}
	storage, err := storage.NewStorage("sqlite", "file::memory:", compressor, cache.New(0))
	if err != nil {
		b.Fatal(err)
	}
	defer storage.Close()

	files, err := ioutil.ReadDir("./testdata/thumbs")
	if err != nil {
		b.Fatal(err)
	}

	tree := New(storage)
	for _, file := range files {
		img, err := image.Load(filepath.Join("./testdata/thumbs", file.Name()))
		if err != nil {
			b.Fatal(err)
		}
		if err := tree.Add(img); err != nil {
			b.Fatal(err)
		}
	}

	i := rand.Intn(len(files))
	query, err := image.Load(filepath.Join("./testdata/thumbs", files[i].Name()))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.Search(query, 2)
	}
}
