package main

import (
	"fmt"
	"log"

	"github.com/recoilme/pudge"
)

func main() {
	cfg := pudge.DefaultConfig
	cfg.SyncInterval = 0 // disable every second fsync
	db, err := pudge.Open("./db", cfg)
	if err != nil {
		log.Panic(err)
	}
	// s := "1"
	var res string
	db.Get("key", &res)
	fmt.Printf("\n>>> %v <<< debug\n", res)
	// db.Set("key", s)
}
