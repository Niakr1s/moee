package moe

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const moeStream = "https://listen.moe/stream"

type moeWs struct {
}

func (w *moeWs) Connect() error {
	u := url.URL{Scheme: "wss", Host: "listen.moe", Path: "/gateway_v2"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	return nil
}
