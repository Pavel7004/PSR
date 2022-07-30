package room

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/server/game"
	"github.com/pavel/PSR/pkg/server/subscribe"
	wd "github.com/pavel/PSR/pkg/server/winner-definer"
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

func NewRoom(cfg *RoomConfig) (*Room, error) {
	p := subscribe.NewPublisher()
	game := game.NewGame(cfg.MaxPlayerCount, p)

	subRoomState := subscribe.NewSubscriber(0)
	if err := p.Subscribe(subRoomState, "room_started"); err != nil {
		log.Error().Err(err).Msg("Failed to subscribe to room_started event.")
		return nil, err
	}

	subWinners := subscribe.NewSubscriber(0)
	if err := p.Subscribe(subWinners, "winners"); err != nil {
		log.Error().Err(err).Msg("Failed to subscribe to winners event.")
		return nil, err
	}

	return &Room{
		game:               game,
		cfg:                cfg,
		mtx:                new(sync.Mutex),
		playerToConnection: make(map[string]*websocket.Conn, cfg.MaxPlayerCount),
		roomStateSub:       subRoomState,
		winnersSub:         subWinners,
	}, nil
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
			log.Error().Err(err).Msgf("Incrementing score for player %q error", name)
		}
	}
	for name, conn := range r.playerToConnection {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(messages[getMessage(name)])); err != nil {
			log.Error().Err(err).Msgf("Error sending winner signal to player %q", name)
		}
	}
}

func (r *Room) Main() {
	r.roomStateSub.Receive()
	startMsg := []byte("Игра началась")
	for _, conn := range r.playerToConnection {
		if err := conn.WriteMessage(websocket.TextMessage, startMsg); err != nil {
			log.Error().Err(err).Msgf("Error sending message %q", startMsg)
		}
	}
	for {
		r.RoundProcess()

		leadingPlayerName, err := r.game.GetLeader()
		if err != nil {
			log.Warn().Err(err).Msg("Error getting leader player name")
			break
		}

		leadingPlayerScore, err := r.game.GetPlayerScore(leadingPlayerName)
		if err != nil {
			log.Warn().Err(err).Msgf("Error Getting player %q score", leadingPlayerName)
			break
		}

		if leadingPlayerScore == r.cfg.MaxScore {
			conn := r.playerToConnection[leadingPlayerName]
			if err := conn.WriteMessage(websocket.TextMessage, []byte("Score win")); err != nil {
				log.Error().Err(err).Msgf("Error sending message to player %q", leadingPlayerName)
			}
			break
		}
	}
	r.CloseConnections()
}

func (r *Room) CloseConnections() {
	for id, conn := range r.playerToConnection {
		if err := conn.Close(); err != nil {
			log.Warn().Err(err).Msgf("Player %q: closing connection error", id)
		}
		log.Info().Msgf("Player %q: connection closed", id)
	}
}

func (r *Room) AddPlayer(id string, conn *websocket.Conn) {
	if err := r.game.AddPlayer(domain.NewPlayer(id)); err != nil {
		log.Error().Err(err).Msgf("Error adding player %q", id)
		return
	}

	r.mtx.Lock()
	r.playerToConnection[id] = conn
	r.mtx.Unlock()

	log.Info().Msgf("Player %q: connection established", id)

	go r.listenConn(id, conn)
}

func (r *Room) listenConn(id string, conn *websocket.Conn) {
	for {
		tMsg, msg, err := conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			break
		}
		if err != nil {
			log.Warn().Err(err).Msgf("Reading message from player %q error", id)
			err = conn.Close()
			if err != nil {
				log.Warn().Err(err).Msgf("Error closing connection for player %q", id)
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
			log.Info().Err(err).Msgf("Player %q: invalid choice %v", id, string(msg))
			continue
		}
		err = r.game.Choose(&wd.PlayerChoice{
			PlayerID: id,
			Input:    choice,
		})
		if err != nil {
			log.Info().Err(err).Msgf("Can't accept player %q choice.", id)
			continue
		}
	}
}
