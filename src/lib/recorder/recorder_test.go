package recorder

import (
	"testing"
	"time"
)

func TestRecorder_StartRecord(t *testing.T) {
	rec := NewRecorder("https://listen.moe/stream")

	err := rec.StartRecord()
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Millisecond * 500)
	if rec.isRecording != true {
		t.Fatalf("rec.isRecording != true")
	}

	select {
	case _, ok := <-rec.TrackCh():
		t.Fatalf("got something from rec.TrackCh(): ok=%v", ok)
	default:
	}
}
