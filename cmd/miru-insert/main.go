package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

type options struct {
	db       string
	folder   string
	parallel uint
	profile  bool
}

func main() {
	var o options
	flag.StringVar(&o.db, "db", os.Getenv("MIRU_DB"), "database name")
	flag.StringVar(&o.folder, "folder", ".", "target folder")
	flag.UintVar(&o.parallel, "parallel", 10, "number of files to process in parallel")
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

	if err := insert(o); err != nil {
		log.Fatal(err)
	}
}
