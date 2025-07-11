package main

import (
	"flag"
	"log"
	"os"

	"github.com/ixpectus/declarate/formatter"
)

var flagSourceDir = flag.String(
	"source_dir", "", "tests directory",
)

var flagTargetDir = flag.String(
	"target_dir", "", "tests directory",
)

func main() {
	flag.Parse()
	if *flagSourceDir == "" {
		log.Println("source directory empty, pass -source_dir flag")
		return
	}
	if *flagTargetDir == "" {
		log.Println("target directory empty, pass -target_dir flag")
		return
	}
	c := formatter.New(*flagSourceDir, *flagTargetDir)
	if err := c.Format(); err != nil {
		log.Printf("convert failed, %s", err)
		os.Exit(1)
		return
	}
	log.Printf("convert success")
}
