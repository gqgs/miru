package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var db = flag.String("db", os.Getenv("MIRU_DB"), "database name")
	var file = flag.String("file", "", "Target file")
	var accuracy = flag.Uint("accuracy", 2, "Accuracy. Higher = more accurate = slower")
	var limit = flag.Uint("limit", 10, "Number of closest matches to display")
	flag.Parse()

	if err := search(*db, *file, *accuracy, *limit); err != nil {
		log.Fatal(err)
	}
}
