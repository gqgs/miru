package main

import (
	"flag"
	"log"
	"os"
)

type options struct {
	db         string
	out        string
	compressor string
}

func main() {
	var o options
	flag.StringVar(&o.db, "db", os.Getenv("MIRU_DB"), "database name")
	flag.StringVar(&o.out, "out", "digraph.dot", "output file")
	flag.StringVar(&o.compressor, "compressor", "zstd", "compression algorithm")
	flag.Parse()

	if err := plot(o); err != nil {
		log.Fatal(err)
	}
}
