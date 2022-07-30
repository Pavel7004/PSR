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

func (game *Game) HasPlayer(playerName string) bool {
	for _, pl := range game.players {
		if pl.GetID() == playerName {
			return true
		}
	}
	return false
}

func (game *Game) AddPlayer(player *domain.Player) error {
	return game.state.AddPlayer(player)
}

func (game *Game) Choose(choice *PlayerChoice) error {
	return game.state.Choose(choice)
}

func (game *Game) GetLeader() (string, error) {
	return game.state.GetLeader()
}

func (game *Game) GetPlayerScore(name string) (uint64, error) {
	return game.state.GetPlayerScore(name)
}

func (game *Game) IncPlayerScore(name string) error {
	return game.state.IncPlayerScore(name)
}
