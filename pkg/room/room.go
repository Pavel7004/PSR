package room

import (
	"errors"
	"sync"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/subscribe"
	. "github.com/pavel/PSR/pkg/winner-definer"
)

var (
	ErrGameAlreadyStarted = errors.New("The game is already started!")
	ErrGameNotStarted     = errors.New("The game isn't started!")
	ErrPlayerNotPresent   = errors.New("The player isn't exist in the room")
)

type Room struct {
	config        *RoomConfig
	players       []*domain.Player
	combinations  []PlayerChoice
	state         State
	observer      *subscribe.Publisher
	stepMtx       *sync.Mutex
	winnerDefiner *WinnerDefiner
}

func NewRoom(config *RoomConfig, obs *subscribe.Publisher) *Room {
	room := Room{
		config:        config,
		players:       make([]*domain.Player, 0, config.MaxPlayerCount),
		combinations:  []PlayerChoice{},
		state:         nil,
		observer:      obs,
		stepMtx:       new(sync.Mutex),
		winnerDefiner: &WinnerDefiner{},
	}
	room.state = NewWaitingState(&room)
	return &room
}

func (room *Room) HasPlayer(playerName string) bool {
	for _, pl := range room.players {
		if pl.GetID() == playerName {
			return true
		}
	}
	return false
}

func (room *Room) AddPlayer(player *domain.Player) error {
	return room.state.AddPlayer(player)
}

func (room *Room) Choose(choice *PlayerChoice) error {
	return room.state.Choose(choice)
}

func (room *Room) MaxScore() (*domain.Player, error) {
	return room.state.MaxScore()
}

func (room *Room) IncPlayerScore(name string) error {
	return room.state.IncPlayerScore(name)
}
