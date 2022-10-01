package moe

import (
	"testing"
	"time"
)

func Test_moeWs_Connect(t *testing.T) {
	w := &moeWs{}
	err := w.Connect()
	if err != nil {
		t.Fatalf("couldn't connect to moe ws server: %v", err)
	}

	msg := <-w.wsTrackCh
	t.Logf("got track message with type %s", msg.Type)

	time.After(time.Millisecond * 300)
	w.close()

	<-w.doneCh
	if _, ok := <-w.wsTrackCh; ok {
		t.Fatalf("got valid track message even after close: %v", err)
	}
}
