package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var db = flag.String("db", os.Getenv("MIRU_DB"), "database name")
	var folder = flag.String("folder", ".", "target folder")
	flag.Parse()

	if err := insert(*db, *folder); err != nil {
		log.Fatal(err)
	}
}
