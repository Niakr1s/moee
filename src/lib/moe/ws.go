package moe

import (
	"encoding/json"
	"fmt"
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
	res := fmt.Sprintf("ID: %d, Title: %s, Duration: %v", s.ID, s.Title, s.Duration)
	if len(s.Artists) > 0 {
		res += fmt.Sprintf(", Artist: %s", s.Artists[0].Name)
	}
	return res
}

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

type MoeWs struct {
	doneCh chan struct{}
	conn   *websocket.Conn

	// needed to close
	trackInfoCh chan TrackInfo
}

func (w *MoeWs) TrackInfoCh() <-chan TrackInfo {
	return w.trackInfoCh
}

func (w *MoeWs) Connect() error {
	w.doneCh = make(chan struct{})
	w.trackInfoCh = make(chan TrackInfo)

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
		close(w.trackInfoCh)
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
				interval := msg.Data.Heartbeat
				w.startHeartbeat(interval)

			case 1:
				// it's track message
				msg := TrackInfo{}
				err = json.Unmarshal(message, &msg)
				if err != nil {
					log.Println("json.Unmarshal trackInfo:", err)
					return
				}
				w.trackInfoCh <- msg

			case 10:
				log.Printf("heartbeat confirmed")

			default:
				log.Printf("wsMessage: unknown message Op; %d", msg.Op)
				return
			}
		}
	}()

	return nil
}

func (w *MoeWs) close() {
	w.conn.Close() // after this, w.doneCh will be closed automatically
}

func (w *MoeWs) startHeartbeat(interval int) {
	go func() {
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(interval)):
				w.sendHeartbeat()
			case <-w.doneCh:
				return
			}
		}
	}()
}

func (w *MoeWs) sendHeartbeat() {
	err := w.conn.WriteJSON(wsMessage{Op: 9})
	if err != nil {
		log.Printf("couldn't send heartbeat: %v", err)
	} else {
		log.Printf("heartbeat send")
	}
}
