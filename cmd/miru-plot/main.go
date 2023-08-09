package main

import (
	"log"
	"os"
)

//go:generate go run github.com/gqgs/argsgen

type options struct {
	storage    string `arg:"storage"`
	db         string `arg:"database name (sqlite)"`
	out        string `arg:"output file"`
	compressor string `arg:"compression algorithm"`
}

func main() {
	o := options{
		storage:    "sqlite",
		db:         os.Getenv("MIRU_DB"),
		out:        "digraph.dot",
		compressor: "zstd",
	}
	o.MustParse()

	if err := plot(o); err != nil {
		log.Fatal(err)
	}
}
