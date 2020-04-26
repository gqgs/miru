package main

import (
	"fmt"
	"os/exec"

	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/image"
	"github.com/gqgs/miru/pkg/storage"
	"github.com/gqgs/miru/pkg/tree"
)

func search(o options) error {
	compressor := compress.NewGzip()
	sqliteStorage, err := storage.NewSqliteStorage(o.db, compressor)
	if err != nil {
		return err
	}
	defer sqliteStorage.Close()

	img, err := image.Load(o.file)
	if err != nil {
		return err
	}

	res, err := tree.New(sqliteStorage).Search(img, o.accuracy)
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
