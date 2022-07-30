package game

import (
	"errors"
	"sync"

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
	combinations  []PlayerChoice
	state         State
	observer      *subscribe.Publisher
	stepMtx       *sync.Mutex
	winnerDefiner *WinnerDefiner
	scoremanager  *ScoreManager
}

func NewGame(playerCount int, obs *subscribe.Publisher) *Game {
	game := &Game{
		players:       make([]*domain.Player, 0, playerCount),
		combinations:  make([]PlayerChoice, 0, playerCount),
		state:         nil,
		observer:      obs,
		stepMtx:       new(sync.Mutex),
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

func (game *Game) GetPlayerCount() int {
	return len(game.players)
}

func (game *Game) Start() {
	game.state = NewPlayingState(game)
	game.scoremanager = NewScoreManager(game.players)
}

func (game *Game) ResetCombinations() {
	game.combinations = make([]PlayerChoice, 0, len(game.combinations))
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
