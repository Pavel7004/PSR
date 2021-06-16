package main

import (
	"net/http"
	"text/template"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
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

	wRoom := NewWebRoom("test", 2)
	go wRoom.StartGame()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("templates/index.html")
		tmpl.Execute(w, nil)
	})

	r.Get("/echo", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Upgrade socket error", 500)
			return
		}
		wRoom.connections[id] = conn
		wRoom.AddPlayer(id)
		log.Info().Msgf("[server] Player %s added to the room", id)

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
