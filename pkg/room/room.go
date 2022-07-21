package room

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/game"
	"github.com/pavel/PSR/pkg/subscribe"
	winner_definer "github.com/pavel/PSR/pkg/winner-definer"
	"github.com/rs/zerolog/log"
)

type Room struct {
	game               *game.Game
	cfg                *RoomConfig
	mtx                *sync.Mutex
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

func NewRoom(cfg *RoomConfig) *Room {
	p := subscribe.NewPublisher()
	game := game.NewGame(cfg.MaxPlayerCount, p)
	subRoomState := subscribe.NewSubscriber(0)
	p.Subscribe(subRoomState, "room_started")
	subWinners := subscribe.NewSubscriber(0)
	p.Subscribe(subWinners, "winners")
	return &Room{
		game:               game,
		mtx:                new(sync.Mutex),
		playerToConnection: make(map[string]*websocket.Conn, cfg.MaxPlayerCount),
		roomStateSub:       subRoomState,
		winnersSub:         subWinners,
	}
}

func (r *Room) RoundProcess() {
	winners, ok := r.winnersSub.Receive().([]string)
	if !ok {
		log.Error().Msgf("Received wrong winners type, got = %T, expected = []string", winners)
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
		if err := r.game.IncPlayerScore(name); err != nil {
			log.Error().Err(err).Msgf("Incrementing score for player \"%s\" error", name)
		}
	}
	for name, conn := range r.playerToConnection {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(messages[getMessage(name)])); err != nil {
			log.Error().Err(err).Msgf("Error sending winner signal to player \"%s\"", name)
		}
	}
}

func (r *Room) Main() {
	r.roomStateSub.Receive()
	startMsg := []byte("Игра началась")
	for _, conn := range r.playerToConnection {
		if err := conn.WriteMessage(websocket.TextMessage, startMsg); err != nil {
			log.Error().Err(err).Msgf("Error sending message \"%s\"", startMsg)
		}
	}
	for {
		r.RoundProcess()
		leadingPlayerName, err := r.game.GetLeader()
		if err != nil {
			log.Warn().Err(err).Msg("Error Getting max score")
			break
		}
		leadingPlayerScore, err := r.game.GetPlayerScore(leadingPlayerName)
		if leadingPlayerScore == r.cfg.MaxScore {
			conn := r.playerToConnection[leadingPlayerName]
			if err := conn.WriteMessage(websocket.TextMessage, []byte("Score win")); err != nil {
				log.Error().Err(err).Msgf("Error sending message to player \"%s\"", leadingPlayerName)
			}
			break
		}
	}
	r.CloseConnections()
}

func (r *Room) CloseConnections() {
	for id, conn := range r.playerToConnection {
		if err := conn.Close(); err != nil {
			log.Warn().Err(err).Msgf("Player \"%s\": closing connection error", id)
		}
		log.Info().Msgf("Player \"%s\": connection closed", id)
	}
}

func (r *Room) AddPlayer(id string, conn *websocket.Conn) {
	if err := r.game.AddPlayer(domain.NewPlayer(id)); err != nil {
		log.Error().Err(err).Msgf("Error adding player \"%s\"", id)
		return
	}
	r.mtx.Lock()
	r.playerToConnection[id] = conn
	r.mtx.Unlock()
	log.Info().Msgf("Player %s: connection established", id)
	go r.listenConn(id, conn)
}

func (r *Room) listenConn(id string, conn *websocket.Conn) {
	for {
		tMsg, msg, err := conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			break
		}
		if err != nil {
			log.Warn().Err(err).Msgf("Reading message from %s error", id)
			err = conn.Close()
			if err != nil {
				log.Warn().Err(err).Msgf("Error closing connection for player %s", id)
			}
			break
		}
		if tMsg != websocket.TextMessage {
			log.Warn().Msgf("Message type isn't text, got = %v, expected = %v", tMsg, websocket.TextMessage)
			continue
		}
		log.Info().Msgf("Got message from %s: %v", id, string(msg))
		choice, err := domain.GetChoiceByName(string(msg))
		if err != nil {
			log.Warn().Err(err).Msgf("Player %s: invalid choice %v", id, string(msg))
			continue
		}
		r.game.Choose(&winner_definer.PlayerChoice{
			PlayerID: id,
			Input:    choice,
		})
	}
}
