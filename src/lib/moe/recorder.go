package moe

import (
	"log"

	"github.com/niakr1s/moee/src/lib/recorder"
)

type Recorder struct {
	rec *recorder.Recorder
	ws  *MoeWs

	dir string
}

func NewRecorder(dir string) *Recorder {
	return &Recorder{
		rec: recorder.NewRecorder("https://listen.moe/stream"),
		ws:  &MoeWs{},
		dir: dir,
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

	trackInfoCh := rec.ws.trackInfoCh
	trackCh := rec.rec.TrackCh()

	var prevTrackInfo TrackInfo
	var currentTrackInfo TrackInfo

	for {
		select {
		// trackInfo usually comes faster then track
		case info := <-trackInfoCh:
			log.Printf("got track info: %v", info)
			prevTrackInfo = currentTrackInfo
			currentTrackInfo = info

		case track := <-trackCh:
			log.Printf("got track: %v", track)
			err := WriteTrack(rec.dir, prevTrackInfo)
			if err != nil {
				log.Printf("err while WriteTrack(rec.dir, prevTrackInfo): %v", err)
				continue
			}
		}
	}
}

func WriteTrack(dir string, trackInfo TrackInfo) error {
	return nil
}
