package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

type options struct {
	db       string
	file     string
	accuracy uint
	limit    uint
	open     bool
	profile  bool
}

func main() {
	var o options
	flag.StringVar(&o.db, "db", os.Getenv("MIRU_DB"), "database name")
	flag.StringVar(&o.file, "file", "", "target file|url")
	flag.UintVar(&o.accuracy, "accuracy", 2, "accuracy. higher = more accurate = slower")
	flag.UintVar(&o.limit, "limit", 10, "number of closest matches to display")
	flag.BoolVar(&o.open, "open", false, "open closest match")
	flag.BoolVar(&o.profile, "cpuprofile", false, "create CPU profile")
	flag.Parse()

	if o.profile {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if err := search(o); err != nil {
		log.Fatal(err)
	}
}
