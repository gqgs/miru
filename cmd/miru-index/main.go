package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

//go:generate go run github.com/gqgs/argsgen

type options struct {
	storage    string `arg:"storage"`
	db         string `arg:"database name (sqlite)"`
	folder     string `arg:"target folder,+,!"`
	parallel   uint   `arg:"number of files to process in parallel"`
	profile    bool   `arg:"create CPU profile"`
	compressor string `arg:"compression algorithm"`
	cachesize  int    `arg:"max number of elements in cache"`
}

func main() {
	o := options{
		storage:    "sqlite",
		db:         os.Getenv("MIRU_DB"),
		parallel:   uint(runtime.NumCPU()),
		compressor: "zstd",
		cachesize:  256,
	}
	o.MustParse()

	if o.profile {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if err := index(o); err != nil {
		log.Print(err)
	}
}
