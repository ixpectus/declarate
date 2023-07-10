package main

import (
	"fmt"
	"time"

	"github.com/ixpectus/declarate/output"
)

func main() {
	bar := output.NewBar(time.Now().Add(10 * time.Second))
	go bar.Start()
	time.Sleep(2 * time.Second)
	bar.Stop()
	time.Sleep(2 * time.Second)
	fmt.Printf("\n>>> %v <<< debug\n", "finish")
}
