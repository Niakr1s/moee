package moe

import (
	"fmt"
	"time"
)

const (
	TRACK_UPDATE         = "TRACK_UPDATE"
	TRACK_UPDATE_REQUEST = "TRACK_UPDATE_REQUEST"
)

type TrackInfo struct {
	Op   int `json:"op"`
	Data struct {
		Song      Song      `json:"song"`
		StartTime time.Time `json:"startTime"`

		LastPlayed []Song `json:"lastPlayed"`
		Listeners  int    `json:"listeners"`
	} `json:"d"`
	Type string `json:"t"` // TRACK_UPDATE or TRACK_UPDATE_REQUEST
}

func (s TrackInfo) String() string {
	return fmt.Sprintf("Type: %s, StartTime: %v, Song: [%v]", s.Type, s.Data.StartTime, s.Data.Song)
}

type Song struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Artists []struct {
		ID         int         `json:"id"`
		Name       string      `json:"name"`
		NameRomaji interface{} `json:"nameRomaji"`
		Image      string      `json:"image"`
	} `json:"artists"`
	Duration int `json:"duration"`
}

func (s Song) String() string {
	res := fmt.Sprintf("ID: %d, Title: %s, Artist: %s, Duration: %v", s.ID, s.Title, s.Artist(), s.Duration)
	return res
}

func (s Song) Artist() string {
	if len(s.Artists) == 0 {
		return ""
	}
	return s.Artists[0].Name
}

// returns filename in format Artist - Title (without suffix)
func (s Song) SuggestedFileName() string {
	return fmt.Sprintf("%s - %s", s.Artist(), s.Title)
}
