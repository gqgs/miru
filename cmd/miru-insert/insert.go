package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/image"
	"github.com/gqgs/miru/pkg/storage"
	"github.com/gqgs/miru/pkg/tree"
)

func insert(o options) error {
	compressor, err := compress.NewCompressor(o.compressor)
	if err != nil {
		return err
	}
	sqliteStorage, err := storage.NewSqliteStorage(o.db, compressor)
	if err != nil {
		return err
	}
	defer sqliteStorage.Close()

	tree := tree.New(sqliteStorage)

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
