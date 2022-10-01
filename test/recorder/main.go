package main

import (
	"log"

	"github.com/niakr1s/moee/src/lib/recorder"
)

func main() {
	for {
		log.Printf("start recording")
		r := recorder.NewRecorder("https://listen.moe/stream")
		r.DiscardFirstTrack = false

		err := r.StartRecord()
		if err != nil {
			log.Printf("record error: %v", err)
		}
		for {
			track := <-r.TrackCh()
			log.Printf("got track %v", track)
			onNewTrack(track)
		}
	}
}

func onNewTrack(track recorder.Track) {

}
