package moe

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	metadata "github.com/niakr1s/moee/src/lib/metatdata"
	"github.com/niakr1s/moee/src/lib/recorder"
	"github.com/raitonoberu/ytmusic"
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

			lyrics, err := GetLyrics(trackInfo.Data.Song.SuggestedFileName())
			if err != nil {
				log.Printf("err while GetLyrics: %v", err)
				continue
			}

			err = metadata.WriteLyrics(savedPath, lyrics)
			if err != nil {
				log.Printf("err while WriteLyrics: %v", err)
				continue
			}
		}
	}
}

// returns full saved filepath
func WriteTrack(dir string, extension string, track recorder.Track, trackInfo TrackInfo) (string, error) {
	song := trackInfo.Data.Song
	path := filepath.Join(dir, song.SuggestedFileName()+extension)
	err := os.WriteFile(path, track.Raw.Bytes(), 0666)
	if err != nil {
		return "", fmt.Errorf("WriteTrack() error: %v", err)
	}
	return path, nil
}

func GetLyrics(query string) (string, error) {
	searchResult, err := ytmusic.Search(query).Next()
	if err != nil {
		return "", err
	}
	if len(searchResult.Tracks) == 0 {
		return "", fmt.Errorf("got zero tracks")
	}
	track := searchResult.Tracks[0]
	videoId := track.VideoID

	return ytmusic.GetLyrics(videoId)
}
