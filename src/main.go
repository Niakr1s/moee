package main

import (
	"log"

	"github.com/niakr1s/moee/src/lib/moe"
)

func main() {
	rec := moe.NewRecorder("d:/tmp")
	rec.DiscardFirstTrack = true
	rec.SaveOnlyWithLyrics = true
	if err := rec.Start(); err != nil {
		log.Fatal(err)
	}
}
