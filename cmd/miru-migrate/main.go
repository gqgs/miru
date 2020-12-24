package main

import (
	"flag"
	"log"
	"os"

	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/storage"
	"github.com/gqgs/miru/pkg/tree"
)

type options struct {
	db string
}

func main() {
	var o options
	flag.StringVar(&o.db, "db", os.Getenv("MIRU_DB"), "database name")
	flag.Parse()

	if err := migrate(o); err != nil {
		log.Fatal(err)
	}
}

func migrate(o options) error {
	compressor := compress.NewMigrateCompressor()
	sqliteStorage, err := storage.NewSqliteStorage(o.db, compressor)
	if err != nil {
		return err
	}
	defer sqliteStorage.Close()

	return tree.New(sqliteStorage).Migrate()
}
