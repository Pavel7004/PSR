package main

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/room"
	"github.com/pavel/PSR/pkg/subscribe"
	winner_definer "github.com/pavel/PSR/pkg/winner-definer"
	"github.com/rs/zerolog/log"
)

var (
	ErrWrongWinnersType = errors.New("Received wrong winners type")
)

type WebRoom struct {
	name               string
	room               *room.Room
	connectionToPlayer map[*websocket.Conn]string
	playerToConnection map[string]*websocket.Conn
	roomStateSub       *subscribe.Subscriber
	winnersSub         *subscribe.Subscriber
	err                chan error
}

func NewWebRoom(name string, maxPlayers int) *WebRoom {
	p := subscribe.NewPublisher()
	rm := room.NewRoom(
		room.RoomConfig{
			StepTimeout:    5 * time.Second,
			MaxPlayerCount: maxPlayers,
			MaxScore:       5,
			OnlyComputer:   false,
		},
		p,
	)
	subRoomState := subscribe.NewSubscriber(0)
	p.Subscribe(subRoomState, "room_started")
	subWinners := subscribe.NewSubscriber(0)
	p.Subscribe(subWinners, "winners")
	return &WebRoom{
		name:               name,
		room:               rm,
		connectionToPlayer: make(map[*websocket.Conn]string, maxPlayers),
		playerToConnection: make(map[string]*websocket.Conn, maxPlayers),
		roomStateSub:       subRoomState,
		winnersSub:         subWinners,
		err:                make(chan error, 1),
	}
}

func (r *WebRoom) GameProcess() {
	r.roomStateSub.Receive()
	startMsg := []byte("Игра началась")
	for conn := range r.connectionToPlayer {
		err := conn.WriteMessage(websocket.TextMessage, startMsg)
		if err != nil {
			log.Error().Err(err).Msgf("[WebRoom:%s] Error sending message \"%s\"", r.name, startMsg)
			r.err <- err
		}
	}
	winners, err := r.winnersSub.Receive().([]string)
	if !err {
		log.Error().Msgf("[WebRoom:%s] Received wrong winners type, got = %T, expected = []string", r.name, winners)
		r.err <- ErrWrongWinnersType
		return
	}
	for _, name := range winners {
		if r.room.HasPlayer(name) {
			err := r.playerToConnection[name].WriteMessage(websocket.TextMessage, []byte("win"))
			if err != nil {
				log.Error().Err(err).Msgf("[WebRoom:%s] Error sending winner signal to player \"%s\"", r.name, name)
			}
		}
	}
	r.err <- nil
	r.CloseConnections()
}

func (r *WebRoom) CloseConnections() {
	for conn, id := range r.connectionToPlayer {
		conn.Close()
		log.Info().Msgf("[WebRoom:%s] Player %s: connection closed", r.name, id)
	}
}

func (r *WebRoom) AddPlayer(id string, conn *websocket.Conn) {
	err := r.room.AddPlayer(domain.NewPlayer(id))
	if err != nil {
		log.Error().Err(err).Msgf("[WebRoom:%s] Error adding player \"%s\"", r.name, id)
		r.err <- err
		return
	}
	r.connectionToPlayer[conn] = id
	r.playerToConnection[id] = conn
	log.Info().Msgf("[WebRoom:%s] Player %s: connection established", r.name, id)
	go r.listenConn(conn)
}

func (r *WebRoom) listenConn(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			log.Warn().Err(err).Msgf("[WebRoom:%s] reading message from %s error", r.name, r.connectionToPlayer[conn])
			r.err <- err
			break
		}
		log.Info().Msgf("[WebRoom:%s] Got message from %s: %v", r.name, r.connectionToPlayer[conn], string(msg))
		choice, err := domain.GetChoiceByName(string(msg))
		if err != nil {
			log.Warn().Err(err).Msgf("[WebRoom:%s] player %s: invalid choice %v", r.name, r.connectionToPlayer[conn], string(msg))
			continue
		}
		r.room.Choose(&winner_definer.PlayerChoice{
			PlayerID: r.connectionToPlayer[conn],
			Input:    choice,
		})
	}
}

func (r *WebRoom) isErrorOccured() error {
	return <-r.err
}
