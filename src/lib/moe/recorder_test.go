package moe

import "testing"

func TestGetLyrics(t *testing.T) {
	got, err := GetLyrics("茅野愛衣 - Rising in revolt")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got lyrics with len=%d", len(got))
}
