package main

import (
	"fmt"
	"miru/pkg/compress"
	"miru/pkg/image"
	"miru/pkg/storage"
	"miru/pkg/tree"
	"os/exec"
)

func search(o options) error {
	compressor := compress.NewGzip()
	sqliteStorage, err := storage.NewSqliteStorage(o.db, compressor)
	if err != nil {
		return err
	}
	defer sqliteStorage.Close()

	tree, err := tree.New(sqliteStorage)
	if err != nil {
		return err
	}

	img, err := image.Load(o.file)
	if err != nil {
		return err
	}

	res, err := tree.Search(img, o.accuracy)
	if err != nil {
		return err
	}

	top := res.Top(o.limit)
	for _, t := range top {
		fmt.Println(t)
	}

	if o.open && len(top) > 0 {
		return exec.Command("xdg-open", top[0].Filename).Start()
	}

	return nil
}
