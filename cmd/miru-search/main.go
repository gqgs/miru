package main

import (
	"flag"
	"log"
	"os"
)

type options struct {
	db       string
	file     string
	accuracy uint
	limit    uint
	open     bool
}

func main() {
	var o options
	flag.StringVar(&o.db, "db", os.Getenv("MIRU_DB"), "database name")
	flag.StringVar(&o.file, "file", "", "Target file|url")
	flag.UintVar(&o.accuracy, "accuracy", 2, "Accuracy. Higher = more accurate = slower")
	flag.UintVar(&o.limit, "limit", 10, "Number of closest matches to display")
	flag.BoolVar(&o.open, "open", false, "Open closest match")
	flag.Parse()

	if err := search(o); err != nil {
		log.Fatal(err)
	}
}
