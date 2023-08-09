package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/image"
	"github.com/gqgs/miru/pkg/storage"
	"github.com/gqgs/miru/pkg/tree"
	"github.com/schollz/progressbar/v3"
)

func index(o options) error {
	compressor, err := compress.NewCompressor(o.compressor)
	if err != nil {
		return err
	}
	storage, err := storage.NewStorage(o.storage, o.db, compressor, cache.New(o.cachesize))
	if err != nil {
		return err
	}
	defer storage.Close()

	tree := tree.New(storage)

	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("imgs"),
	)
	// nolint: errcheck
	defer bar.Finish()

	var wg sync.WaitGroup
	pathCh := make(chan string)
	go func() {
		semaphore := make(chan struct{}, o.parallel)
		for path := range pathCh {
			path := path
			semaphore <- struct{}{}
			go func() {
				defer func() {
					<-semaphore
					wg.Done()
					_ = bar.Add(1)
				}()
				img, err := image.Load(path)
				if err != nil {
					log.Printf("%s: %s", err, path)
					return
				}
				if err = tree.Add(img); err != nil {
					log.Println("tree", err)
				}
			}()
		}
	}()

	err = filepath.Walk(o.folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() && image.IsImage(info.Name()) {
			wg.Add(1)
			pathCh <- path
		}
		return nil
	})
	wg.Wait()

	return err
}
