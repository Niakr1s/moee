package moe

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type wsMessage struct {
	Op int `json:"op"`
}

type wsWelcomeMsg struct {
	Op   int `json:"op"`
	Data struct {
		Message   string `json:"message"`
		Heartbeat int    `json:"heartbeat"`
	} `json:"d"`
}

type wsSong struct {
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

const (
	TRACK_UPDATE         = "TRACK_UPDATE"
	TRACK_UPDATE_REQUEST = "TRACK_UPDATE_REQUEST"
)

type wsTrackMsg struct {
	Op   int `json:"op"`
	Data struct {
		Song      wsSong    `json:"song"`
		StartTime time.Time `json:"startTime"`

		LastPlayed []wsSong `json:"lastPlayed"`
		Listeners  int      `json:"listeners"`
	} `json:"d"`
	Type string `json:"t"` // TRACK_UPDATE or TRACK_UPDATE_REQUEST
}

type moeWs struct {
	doneCh chan struct{}
	conn   *websocket.Conn

	// needed to close
	wsTrackCh chan wsTrackMsg

	heartbeat int //interval
}

func (w *moeWs) Connect() error {
	w.doneCh = make(chan struct{})
	w.wsTrackCh = make(chan wsTrackMsg)

	u := url.URL{Scheme: "wss", Host: "listen.moe", Path: "/gateway_v2"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	w.conn = c

	// listens to done channel and closes everything
	go func() {
		<-w.doneCh
		close(w.wsTrackCh)
	}()

	// will close w.doneCh on any error
	go func() {
		defer close(w.doneCh)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			msg := wsMessage{}
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("json.Unmarshal(message, &wsMsg):", err)
				return
			}
			log.Printf("recv msg.Op: %d", msg.Op)

			switch msg.Op {
			case 0:
				// it's welcome message
				msg := wsWelcomeMsg{}
				err = json.Unmarshal(message, &msg)
				if err != nil {
					log.Println("json.Unmarshal wsWelcomeMsg:", err)
					return
				}
				log.Printf("got welcome message: %s", msg.Data.Message)
				w.heartbeat = msg.Data.Heartbeat

			case 1:
				// it's track message
				msg := wsTrackMsg{}
				err = json.Unmarshal(message, &msg)
				if err != nil {
					log.Println("json.Unmarshal wsTrackMsg:", err)
					return
				}
				w.wsTrackCh <- msg

			default:
				log.Printf("wsMessage: unknown message Op; %d", msg.Op)
				return
			}

		}
	}()

	return nil
}

func (w *moeWs) close() {
	w.conn.Close() // after this, w.doneCh will be closed automatically
}
