package room

import (
	"errors"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/rs/zerolog/log"
)

var (
	ErrGameAlreadyStarted = errors.New("The game is already started!")
)

type Room struct {
	config  RoomConfig
	players []*domain.Player
	active  bool
}

func NewRoom(config RoomConfig) *Room {
	return &Room{
		config:  config,
		active:  false,
		players: make([]*domain.Player, 0, config.MaxPlayerCount),
	}
}

func (room *Room) IsActive() bool {
	return room.active
}

func (room *Room) AddPlayer(player *domain.Player) error {
	if room.IsActive() {
		return ErrGameAlreadyStarted
	}
	room.players = append(room.players, player)
	log.Info().Msgf("Player %s added to the room", player.ID)

	if len(room.players) == room.config.MaxPlayerCount {
		room.active = true
		go room.Run()
	}
	return nil
}

func (room *Room) Run() {
	log.Info().Msg("Room started")
}
