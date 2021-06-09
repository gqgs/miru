package main

import (
	"log"
	"os"
	"runtime/pprof"
)

//go:generate go run github.com/gqgs/argsgen

type options struct {
	db         string `arg:"database name"`
	file, url  string `arg:"target file|url,+,!"`
	accuracy   uint   `arg:"higher = more accurate = slower"`
	limit      uint   `arg:"number of closest matches to display"`
	open       bool   `arg:"open closest match"`
	profile    bool   `arg:"create CPU profile"`
	compressor string `arg:"compression algorithm"`
}

func main() {
	o := options{
		db:         os.Getenv("MIRU_DB"),
		accuracy:   2,
		limit:      10,
		compressor: "zstd",
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

	if err := search(o); err != nil {
		log.Print(err)
	}
}
