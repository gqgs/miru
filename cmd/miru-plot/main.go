package main

import (
	"log"
	"os"
)

//go:generate go run github.com/gqgs/argsgen

type options struct {
	db         string `arg:"database name"`
	out        string `arg:"output file"`
	compressor string `arg:"compression algorithm"`
}

func main() {
	o := options{
		db:         os.Getenv("MIRU_DB"),
		out:        "digraph.dot",
		compressor: "zstd",
	}
	o.MustParse()

	if err := plot(o); err != nil {
		log.Fatal(err)
	}
}
