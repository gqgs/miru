package main

import (
	"flag"
	"log"
	"os"
)

type options struct {
	db       string
	folder   string
	parallel uint
}

func main() {
	var o options
	flag.StringVar(&o.db, "db", os.Getenv("MIRU_DB"), "database name")
	flag.StringVar(&o.folder, "folder", ".", "target folder")
	flag.UintVar(&o.parallel, "parallel", 10, "number of files to process in parallel")
	flag.Parse()

	if err := insert(o); err != nil {
		log.Fatal(err)
	}
}
