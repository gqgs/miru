package main

import (
	"flag"
	"log"
)

func main() {
	var db = flag.String("db", "miru.db", "database name")
	var folder = flag.String("folder", ".", "target folder")
	flag.Parse()

	if err := insert(*db, *folder); err != nil {
		log.Fatal(err)
	}
}
