package moe

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
			trackInfo := prevTrackInfo
			log.Printf("got track: %v", track)
			savedPath, err := WriteTrack(rec.dir, track.Extension, track, trackInfo)
			if err != nil {
				log.Printf("err while WriteTrack: %v", err)
				continue
			}
			log.Printf("saved track with info %s to %s", trackInfo, savedPath)
		}
	}
}

// returns full saved filepath
func WriteTrack(dir string, extension string, track recorder.Track, trackInfo TrackInfo) (string, error) {
	path := filepath.Join(dir, fmt.Sprintf("%s - %s", trackInfo.Data.Song.Artist(), trackInfo.Data.Song.Title), extension)
	err := os.WriteFile(path, track.Raw.Bytes(), 0666)
	if err != nil {
		return "", fmt.Errorf("WriteTrack() error: %v", err)
	}
	return path, nil
}
