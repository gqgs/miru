package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

type options struct {
	db         string
	folder     string
	parallel   uint
	profile    bool
	compressor string
}

func main() {
	var o options
	flag.StringVar(&o.db, "db", os.Getenv("MIRU_DB"), "database name")
	flag.StringVar(&o.folder, "folder", ".", "target folder")
	flag.UintVar(&o.parallel, "parallel", uint(runtime.NumCPU()), "number of files to process in parallel")
	flag.BoolVar(&o.profile, "cpuprofile", false, "create CPU profile")
	flag.StringVar(&o.compressor, "compressor", "zstd", "compression algorithm")
	flag.Parse()

	if o.profile {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if err := index(o); err != nil {
		log.Fatal(err)
	}
}
