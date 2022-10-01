package moe

import (
	"log"

	"github.com/niakr1s/moee/src/lib/recorder"
)

type Recorder struct {
	rec *recorder.Recorder
	ws  *MoeWs
}

func NewRecorder() *Recorder {
	return &Recorder{
		rec: recorder.NewRecorder("https://listen.moe/stream"),
		ws:  &MoeWs{},
	}
}

func (rec *Recorder) Start() error {
	err := rec.ws.Connect()
	if err != nil {
		return err
	}

	err = rec.rec.StartRecord()
	if err != nil {
		return err
	}

	trackInfoCh := rec.ws.wsTrackInfoCh
	trackCh := rec.rec.TrackCh()

	for {
		select {
		case info := <-trackInfoCh:
			log.Printf("info: %v", info)
		case track := <-trackCh:
			log.Printf("track: %v", track)
			// in wsTrack we have already info about new song, so let's copy song's info about
		}
	}
}
