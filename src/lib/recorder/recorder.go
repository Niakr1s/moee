package recorder

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/niakr1s/moee/src/lib/recorder/mp3"
	"github.com/niakr1s/moee/src/lib/recorder/vorbis"
)

var client = new(http.Client)

type Recorder struct {
	DiscardFirstTrack bool
	MaxTracks         int

	url            string
	tracksRecorded int

	extractor    Extractor
	currentTrack Track
	isRecording  bool

	trackCh chan Track
}

func NewRecorder(url string) *Recorder {
	r := &Recorder{
		url: url,
	}
	return r
}

func (rec *Recorder) IsRecording() bool {
	return rec.isRecording
}

func (rec *Recorder) TrackCh() <-chan Track {
	return rec.trackCh
}

// Recorder's main function. Recorder will close trackCh on any error.
func (rec *Recorder) StartRecord() error {
	if rec.IsRecording() {
		return fmt.Errorf("already recording")
	}

	rec.isRecording = true
	rec.trackCh = make(chan Track)

	resp, err := rec.doRequest()
	if err != nil {
		return err
	}

	extractor, err := getExtractor(resp)
	if err != nil {
		return err
	}
	rec.extractor = extractor

	go func() {
		defer resp.Body.Close()
		defer rec.onEndRecord()

		err := rec.loop(resp)
		if err != nil {
			log.Printf("err wilhe recording: %v", err)
		}
	}()

	return nil
}

func (rec *Recorder) onEndRecord() {
	close(rec.trackCh)
	rec.isRecording = false
}

func (rec *Recorder) doRequest() (*http.Response, error) {
	req, err := http.NewRequest("GET", rec.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Icy-MetaData", "1") // Request metadata for icecast mp3 streams.
	return client.Do(req)
}

func (rec *Recorder) loop(resp *http.Response) error {
	// Make reader blocking.
	r := NewWaitReader(resp.Body)

	// The first track is always discarded, as streams usually don't start at
	// the exact end of a track, meaning it is almost certainly going to be
	// incomplete.
	discard := rec.DiscardFirstTrack

	rec.currentTrack = newTrack(rec.extractor)
	for {
		var block bytes.Buffer

		isFirstBlock, err := rec.extractor.ReadBlock(r, &block)
		if err != nil {
			// Reconnect, because this error is usually caused by a
			// file corruption or a network error.
			return fmt.Errorf("error reading block: %v", err)
		}

		isStartOfNewTrack := isFirstBlock && rec.currentTrack.Raw.Len() > 0

		if isStartOfNewTrack {
			rec.currentTrack.End = time.Now()
			if !discard {
				rec.trackCh <- rec.currentTrack

				// Stop after the defined number of tracks (if the option was
				// given).
				rec.tracksRecorded++
				if rec.MaxTracks > 0 && rec.tracksRecorded >= rec.MaxTracks {
					log.Printf("Successfully recorded %v tracks, exiting\n", rec.tracksRecorded)
					os.Exit(0)
				}
			} else {
				// See declaration of `discard`.
				discard = false
			}
			rec.currentTrack = newTrack(rec.extractor)
		}

		// Append block to the current file byte buffer.
		rec.currentTrack.Raw.Write(block.Bytes())
	}
}

func newTrack(extractor Extractor) Track {
	return Track{
		Start:     time.Now(),
		Extension: extractor.GetExtension(),
	}
}

func getExtractor(resp *http.Response) (Extractor, error) {
	// Set up extractor depending on content type.
	contentType := resp.Header.Get("content-type")

	log.Printf("Stream type: '%v'\n", contentType)

	switch contentType {
	case "application/ogg", "audio/ogg", "audio/vorbis", "audio/vorbis-config":
		return vorbis.NewExtractor()
	case "audio/mpeg", "audio/MPA", "audio/mpa-robust":
		return mp3.NewExtractor(resp.Header)
	default:
		return nil, fmt.Errorf(`content type '%v' not supported, supported formats: `+
			`Ogg/Vorbis ('application/ogg', 'audio/ogg', 'audio/vorbis', 'audio/vorbis-config'), `+
			`mp3 ('audio/mpeg', 'audio/MPA', 'audio/mpa-robust')`, contentType)
	}
}
