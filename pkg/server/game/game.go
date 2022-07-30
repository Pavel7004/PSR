package game

import (
	"errors"

	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/server/score-manager"
	"github.com/pavel/PSR/pkg/server/subscribe"
	. "github.com/pavel/PSR/pkg/server/winner-definer"
)

var (
	ErrGameAlreadyStarted = errors.New("The game is already started!")
	ErrGameNotStarted     = errors.New("The game isn't started!")
	ErrPlayerNotPresent   = errors.New("The player isn't exist in the room")
)

type Game struct {
	players       []*domain.Player
	maxPlayers    int
	combinations  []PlayerChoice
	state         State
	observer      *subscribe.Publisher
	winnerDefiner *WinnerDefiner
	scoremanager  *ScoreManager
}

func NewGame(playerCount int, obs *subscribe.Publisher) *Game {
	game := &Game{
		players:       make([]*domain.Player, 0, playerCount),
		maxPlayers:    playerCount,
		combinations:  nil,
		state:         nil,
		observer:      obs,
		winnerDefiner: new(WinnerDefiner),
		scoremanager:  nil,
	}
	game.state = NewWaitingState(game)
	return game
}

func (g *Game) HasPlayer(playerName string) bool {
	for _, pl := range g.players {
		if pl.GetID() == playerName {
			return true
		}
	}
	return false
}

func (g *Game) AddPlayer(player *domain.Player) error {
	return g.state.AddPlayer(player)
}

func (g *Game) Choose(choice *PlayerChoice) error {
	return g.state.Choose(choice)
}

func (g *Game) GetLeader() (string, error) {
	return g.state.GetLeader()
}

func (g *Game) GetPlayerScore(id string) (uint64, error) {
	return g.state.GetPlayerScore(id)
}

func (g *Game) IncPlayerScore(id string) error {
	return g.state.IncPlayerScore(id)
}
