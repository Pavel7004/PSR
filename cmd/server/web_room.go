package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/room"
	"github.com/pavel/PSR/pkg/subscribe"
	winner_definer "github.com/pavel/PSR/pkg/winner-definer"
	"github.com/rs/zerolog/log"
)

type WebRoom struct {
	name               string
	room               *room.Room
	config             *room.RoomConfig
	mtx                *sync.Mutex
	connectionToPlayer map[*websocket.Conn]string
	playerToConnection map[string]*websocket.Conn
	roomStateSub       *subscribe.Subscriber
	winnersSub         *subscribe.Subscriber
}

type winType int

const (
	WIN winType = iota + 1
	LOSE
	TIE
)

func NewWebRoom(name string, config *room.RoomConfig) *WebRoom {
	p := subscribe.NewPublisher()
	rm := room.NewRoom(config, p)
	subRoomState := subscribe.NewSubscriber(0)
	p.Subscribe(subRoomState, "room_started")
	subWinners := subscribe.NewSubscriber(0)
	p.Subscribe(subWinners, "winners")
	return &WebRoom{
		name:               name,
		room:               rm,
		config:             config,
		mtx:                new(sync.Mutex),
		connectionToPlayer: make(map[*websocket.Conn]string, config.MaxPlayerCount),
		playerToConnection: make(map[string]*websocket.Conn, config.MaxPlayerCount),
		roomStateSub:       subRoomState,
		winnersSub:         subWinners,
	}
}

func (r *WebRoom) RoundProcess() {
	winners, err := r.winnersSub.Receive().([]string)
	if !err {
		log.Error().Msgf("[WebRoom:%s] Received wrong winners type, got = %T, expected = []string", r.name, winners)
		return
	}

	messages := map[winType]string{
		WIN:  "You won!",
		LOSE: fmt.Sprintf("You lost! Winners: %s", strings.Join(winners, ", ")),
		TIE:  "Draw, everyone lost!",
	}
	getMessage := func(name string) winType {
		winStatus := LOSE
		for _, winner := range winners {
			if name == winner {
				winStatus = WIN
			}
		}
		if len(winners) == 0 {
			winStatus = TIE
		}
		return winStatus
	}

	for _, name := range winners {
		err := r.room.IncPlayerScore(name)
		if err != nil {
			log.Error().Err(err).Msgf("[WebRoom:%s] Incrementing score for player \"%s\" error", r.name, name)
		}
	}

	for conn, name := range r.connectionToPlayer {
		err := conn.WriteMessage(websocket.TextMessage, []byte(messages[getMessage(name)]))
		if err != nil {
			log.Error().Err(err).Msgf("[WebRoom:%s] Error sending winner signal to player \"%s\"", r.name, name)
		}
	}
}

func (r *WebRoom) Main() {
	r.roomStateSub.Receive()
	startMsg := []byte("Игра началась")
	for conn := range r.connectionToPlayer {
		err := conn.WriteMessage(websocket.TextMessage, startMsg)
		if err != nil {
			log.Error().Err(err).Msgf("[WebRoom:%s] Error sending message \"%s\"", r.name, startMsg)
		}
	}
	for {
		r.RoundProcess()
		leadingPlayer, err := r.room.MaxScore()
		if err != nil {
			log.Warn().Err(err).Msgf("[WebRoom:%s] Error Getting max score", r.name)
			break
		}
		if leadingPlayer.GetScore() == r.config.MaxScore {
			conn := r.playerToConnection[leadingPlayer.GetID()]
			err = conn.WriteMessage(websocket.TextMessage, []byte("Score win"))
			if err != nil {
				log.Error().Err(err).Msgf("[WebRoom:%s] Error sending message to player \"%s\"", r.name, leadingPlayer.GetID())
			}
			break
		}
	}
	r.CloseConnections()
}

func (r *WebRoom) CloseConnections() {
	for conn, id := range r.connectionToPlayer {
		err := conn.Close()
		if err != nil {
			log.Warn().Err(err).Msgf("[WebRoom:%s] Player \"%s\": closing connection error", r.name, id)
		}
		log.Info().Msgf("[WebRoom:%s] Player \"%s\": connection closed", r.name, id)
	}
}

func (r *WebRoom) AddPlayer(id string, conn *websocket.Conn) {
	err := r.room.AddPlayer(domain.NewPlayer(id))
	if err != nil {
		log.Error().Err(err).Msgf("[WebRoom:%s] Error adding player \"%s\"", r.name, id)
		return
	}
	r.mtx.Lock()
	r.connectionToPlayer[conn] = id
	r.playerToConnection[id] = conn
	r.mtx.Unlock()
	log.Info().Msgf("[WebRoom:%s] Player %s: connection established", r.name, id)
	go r.listenConn(conn)
}

func (r *WebRoom) listenConn(conn *websocket.Conn) {
	for {
		tMsg, msg, err := conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			break
		}
		if err != nil {
			log.Warn().Err(err).Msgf("[WebRoom:%s] Reading message from %s error", r.name, r.connectionToPlayer[conn])
			err = conn.Close()
			if err != nil {
				log.Warn().Err(err).Msgf("[WebRoom:%s] Error closing connection for player %s", r.name, r.connectionToPlayer[conn])
			}
			break
		}
		if tMsg != websocket.TextMessage {
			log.Warn().Msgf("[WebRoom:%s] Message type isn't text, got = %v, expected = %v", r.name, tMsg, websocket.TextMessage)
			continue
		}
		log.Info().Msgf("[WebRoom:%s] Got message from %s: %v", r.name, r.connectionToPlayer[conn], string(msg))
		choice, err := domain.GetChoiceByName(string(msg))
		if err != nil {
			log.Warn().Err(err).Msgf("[WebRoom:%s] Player %s: invalid choice %v", r.name, r.connectionToPlayer[conn], string(msg))
			continue
		}
		r.room.Choose(&winner_definer.PlayerChoice{
			PlayerID: r.connectionToPlayer[conn],
			Input:    choice,
		})
	}
}
