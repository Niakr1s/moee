package main

import (
	"log"
	"os"
	"time"

	"github.com/niakr1s/moee/src/lib/moe"
	"github.com/niakr1s/moee/src/lib/util"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("USAGE: moee [savedir]")
	}

	savedir := os.Args[1]
	savedir, err := util.ToAbs(savedir)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(savedir, 0666)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("moee.log")
	if err != nil {
		log.Fatalf("couldn't create log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	StartRecordLoop(savedir)
}

func StartRecordLoop(savedir string) {
	for {
		log.Println("Starting moe recorder.")
		rec := moe.NewRecorder(savedir)
		rec.DiscardFirstTrack = true
		if err := rec.Start(); err != nil {
			log.Println("Error while recording, going to restart moe recorder in 5 seconds.")
			<-time.After(time.Second * 5)
			continue
		}
		return
	}
}
