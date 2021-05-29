package tree

import (
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/image"
	"github.com/gqgs/miru/pkg/storage"
)

func Benchmark_Add(b *testing.B) {
	compressor, err := compress.NewCompressor("nop")
	if err != nil {
		b.Fatal(err)
	}
	storage, err := storage.NewSqliteStorage("file::memory:", compressor, cache.New(0))
	if err != nil {
		b.Fatal(err)
	}
	defer storage.Close()

	files, err := ioutil.ReadDir("./testdata/thumbs")
	if err != nil {
		b.Fatal(err)
	}

	imgs := make([]*image.Image, 0, len(files))
	for _, file := range files {
		img, err := image.Load(filepath.Join("./testdata/thumbs", file.Name()))
		if err != nil {
			b.Fatal(err)
		}
		imgs = append(imgs, img)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := New(storage)
		var wg sync.WaitGroup
		imgCh := make(chan *image.Image)
		go func() {
			semaphore := make(chan struct{}, 10)
			for img := range imgCh {
				img := img
				semaphore <- struct{}{}
				go func() {
					defer func() {
						<-semaphore
						wg.Done()
					}()
					if err := tree.Add(img); err != nil {
						b.Log(err)
						b.Fail()
					}
				}()
			}
		}()

		for _, img := range imgs {
			wg.Add(1)
			imgCh <- img
		}
		wg.Wait()
	}
}
