package handler

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pavel/PSR/pkg/server/room"
	roommanager "github.com/pavel/PSR/pkg/server/room-manager"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	upgrader *websocket.Upgrader
	rm       *roommanager.RoomManager
}

func New() *Handler {
	rm := roommanager.New() // TODO: Add room creation
	err := rm.CreateRoom(&room.RoomConfig{
		Name:           "test",
		RoundTimeout:   5 * time.Second,
		MaxPlayerCount: 3,
		MaxScore:       7,
	})
	if err != nil {
		panic(err)
	}

	return &Handler{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024 * 1024,
			WriteBufferSize: 1024 * 1024,
		},
		rm: rm,
	}
}

func (h *Handler) GetIndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse \"templates/index.html\"")
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Error().Err(err).Msg("Failed to execute \"templates/index.html\"")
	}
}

func (h *Handler) GetGamePage(w http.ResponseWriter, r *http.Request) {
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
}

func (h *Handler) OpenSocketConnection(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	roomID := r.URL.Query().Get("roomID")

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Upgrade socket error", 500)
		return
	}

	room, err := h.rm.GetRoomByID(roomID)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to get room by id %q", roomID)
		return
	}

	room.AddPlayer(ID, conn)
}
