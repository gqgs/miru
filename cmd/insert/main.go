package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var db = flag.String("db", os.Getenv("MIRU_DB"), "database name")
	var folder = flag.String("folder", ".", "target folder")
	var parallel = flag.Uint("parallel", 10, "number of files to process in parallel")
	flag.Parse()

	if err := insert(*db, *folder, *parallel); err != nil {
		log.Fatal(err)
	}
}
