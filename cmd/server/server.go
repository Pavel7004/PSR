package main

import (
	"net/http"
	"text/template"
	"time"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/pavel/PSR/pkg/server/room"
	roommanager "github.com/pavel/PSR/pkg/server/room-manager"
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
	server := &http.Server{
		Addr:              ":3000",
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	rm := roommanager.New()

	err := rm.CreateRoom(&room.RoomConfig{
		Name:           "test",
		RoundTimeout:   5 * time.Second,
		MaxPlayerCount: 3,
		MaxScore:       7,
	})
	if err != nil {
		panic(err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse \"templates/index.html\"")
		}

		if err := tmpl.Execute(w, nil); err != nil {
			log.Error().Err(err).Msg("Failed to execute \"templates/index.html\"")
		}
	})

	r.Get("/game", func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("id")
		roomID := r.URL.Query().Get("roomID")

		tmpl, err := template.ParseFiles("templates/game.html")
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse \"templates/game.html\"")
		}

		err = tmpl.Execute(w, struct {
			ID     string
			RoomID string
		}{
			ID:     ID,
			RoomID: roomID,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to execute \"templates/game.html\"")
		}
	})

	r.Get("/echo", func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("id")
		roomID := r.URL.Query().Get("roomID")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Upgrade socket error", 500)
			return
		}

		room, err := rm.GetRoomByID(roomID)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to get room by id %q", roomID)
			return
		}

		room.AddPlayer(ID, conn)
	})

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	log.Info().Msg("Server started on port 3000")
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
