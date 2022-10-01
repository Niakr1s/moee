package moe

import "testing"

func Test_moeWs_Connect(t *testing.T) {
	w := &moeWs{}
	err := w.Connect()
	if err != nil {
		t.Fatalf("couldn't connect to moe ws server: %v", err)
	}
}
