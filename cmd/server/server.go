package main

import (
	"net/http"
	"text/template"
	"time"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/room"
	"github.com/pavel/PSR/pkg/subscribe"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	r := chi.NewRouter()
	//TODO: Create a room

	sub := subscribe.NewSubscriber(0)
	p := subscribe.NewPublisher()
	p.Subscribe(sub, "winners")
	rm := room.NewRoom(
		room.RoomConfig{
			5 * time.Second,
			2,
			5,
			false,
		},
		p,
	)
	p.Subscribe(sub, "room_started")

	conns := make(map[*websocket.Conn]string, 2)
	go func() {
		sub.Receive()
		startMsg := []byte("Игра началась")
		for conn := range conns {
			conn.WriteMessage(websocket.TextMessage, startMsg)
		}

		time.Sleep(2 * time.Second)
		for conn := range conns {
			conn.Close()
		}
	}()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("templates/index.html")
		tmpl.Execute(w, nil)
	})

	r.Get("/echo", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		player := domain.NewPlayer(id)
		log.Info().Msgf("[server] Player %s added to the room", id)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Upgrade socket error", 500)
			return
		}
		conns[conn] = id
		rm.AddPlayer(player)

		// conn.WriteMessage(websocket.TextMessage, []byte("Игра началась"))
		// _, msg, err := conn.ReadMessage()
		// if err != nil {
		// 	http.Error(w, "Reading message error", 500)
		// 	return
		// }

		// rm.Choose(&winner_definer.PlayerChoice{
		// 	PlayerID: id,
		// 	Input: ,
		// })
		// log.Info().Msgf("[client] new message: %s", string(msg))
	})

	r.Get("/choice", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("choice"))
	})

	log.Info().Msg("Server started")
	http.ListenAndServe(":3000", r)
}
