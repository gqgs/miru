package main

import (
	"fmt"
	"os/exec"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/image"
	"github.com/gqgs/miru/pkg/storage"
	"github.com/gqgs/miru/pkg/tree"
)

func search(o options) error {
	compressor, err := compress.NewCompressor(o.compressor)
	if err != nil {
		return err
	}
	storage, err := storage.NewStorage(o.storage, o.db, compressor, cache.New(0))
	if err != nil {
		return err
	}
	defer storage.Close()

	img, err := image.Load(o.file)
	if err != nil {
		return err
	}

	res, err := tree.New(storage).Search(img, o.accuracy)
	if err != nil {
		return err
	}

	top := res.Top(o.limit)
	fmt.Println(tree.NewFormatter(top).JSON(o.json))

	if o.open && len(top) > 0 {
		return exec.Command("xdg-open", top[0].Filename).Start()
	}

	return nil
}
