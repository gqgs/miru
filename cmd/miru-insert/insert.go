package main

import (
	"log"
	"miru/pkg/image"
	"miru/pkg/tree"
	"os"
	"path/filepath"
	"sync"
)

func insert(dbName, folder string, parallel uint) error {
	tree, err := tree.New(dbName)
	if err != nil {
		return err
	}
	defer tree.Close()

	var wg sync.WaitGroup
	pathCh := make(chan string)
	go func() {
		semaphore := make(chan struct{}, parallel)
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

	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			wg.Add(1)
			pathCh <- path
		}
		return nil
	})
	wg.Wait()

	return err
}
