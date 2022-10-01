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
	SaveOnlyWithLyrics bool

	*recorder.Recorder
	ws *MoeWs

	dir string
}

func NewRecorder(dir string) *Recorder {
	return &Recorder{
		Recorder: recorder.NewRecorder("https://listen.moe/stream"),
		ws:       &MoeWs{},
		dir:      dir,
	}
}

func (rec *Recorder) Start() error {
	err := rec.ws.Connect()
	if err != nil {
		return err
	}

	err = rec.StartRecord()
	if err != nil {
		return err
	}

	trackInfoCh := rec.ws.trackInfoCh
	trackCh := rec.TrackCh()

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

			lyrics, err := GetLyrics(trackInfo.Data.Song.SuggestedFileName())
			if err != nil {
				log.Printf("err while GetLyrics: %v", err)
				lyrics = ""
				if rec.SaveOnlyWithLyrics {
					log.Printf("skipping saving track, because SaveOnlyWithLyrics=%v", rec.SaveOnlyWithLyrics)
					continue
				}
			} else {
				log.Printf("got lyrics with len=%d", len(lyrics))
			}

			savedPath, err := WriteTrack(rec.dir, track.Extension, track, trackInfo)
			if err != nil {
				log.Printf("err while WriteTrack: %v", err)
				continue
			}
			log.Printf("saved track with info %s to %s", trackInfo, savedPath)

			err = metadata.WriteClean(savedPath)
			if err != nil {
				log.Printf("err while WriteClean: %v", err)
			}
			err = metadata.WriteArtist(savedPath, trackInfo.Data.Song.Artist())
			if err != nil {
				log.Printf("err while WriteArtist: %v", err)
			}
			err = metadata.WriteTitle(savedPath, trackInfo.Data.Song.Title)
			if err != nil {
				log.Printf("err while WriteTitle: %v", err)
			}

			if lyrics != "" {
				err = metadata.WriteLyrics(savedPath, lyrics)
				if err != nil {
					log.Printf("err while WriteLyrics: %v", err)
					continue
				}
				log.Printf("wrote lyrics with len=%d", len(lyrics))
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
