package game

import (
	"errors"
	"sync"

	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/score-manager"
	"github.com/pavel/PSR/pkg/subscribe"
	. "github.com/pavel/PSR/pkg/winner-definer"
)

var (
	ErrGameAlreadyStarted = errors.New("The game is already started!")
	ErrGameNotStarted     = errors.New("The game isn't started!")
	ErrPlayerNotPresent   = errors.New("The player isn't exist in the room")
)

type Game struct {
	players       []*domain.Player
	combinations  []PlayerChoice
	state         State
	observer      *subscribe.Publisher
	stepMtx       *sync.Mutex
	winnerDefiner *WinnerDefiner
	scoremanager  *ScoreManager
}

func NewGame(playerCount int, obs *subscribe.Publisher) *Game {
	room := Game{
		players:       make([]*domain.Player, 0, playerCount),
		combinations:  []PlayerChoice{},
		state:         nil,
		observer:      obs,
		stepMtx:       new(sync.Mutex),
		winnerDefiner: new(WinnerDefiner),
		scoremanager:  nil,
	}
	room.state = NewWaitingState(&room)
	return &room
}

func (room *Game) HasPlayer(playerName string) bool {
	for _, pl := range room.players {
		if pl.GetID() == playerName {
			return true
		}
	}
	return false
}

func (room *Game) AddPlayer(player *domain.Player) error {
	return room.state.AddPlayer(player)
}

func (room *Game) Choose(choice *PlayerChoice) error {
	return room.state.Choose(choice)
}

func (room *Game) GetLeader() (string, error) {
	return room.state.GetLeader()
}

func (room *Game) GetPlayerScore(name string) (uint64, error) {
	return room.state.GetPlayerScore(name)
}

func (room *Game) IncPlayerScore(name string) error {
	return room.state.IncPlayerScore(name)
}
