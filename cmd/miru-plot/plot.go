package main

import (
	"os"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/storage"
	"github.com/gqgs/miru/pkg/tree"
)

func plot(o options) error {
	compressor, err := compress.NewCompressor(o.compressor)
	if err != nil {
		return err
	}
	sqliteStorage, err := storage.NewSqliteStorage(o.db, compressor, cache.New(0))
	if err != nil {
		return err
	}
	defer sqliteStorage.Close()

	file, err := os.Create(o.out)
	if err != nil {
		return err
	}
	defer file.Close()

	return tree.New(sqliteStorage).Plot(file)
}
