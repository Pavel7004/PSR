package game

import (
	"sync"

	"github.com/pavel/PSR/pkg/domain"

	scoremanager "github.com/pavel/PSR/pkg/server/score-manager"
	. "github.com/pavel/PSR/pkg/server/winner-definer"
	"github.com/rs/zerolog/log"
)

type WaitingState struct {
	game *Game
	mtx  *sync.Mutex
}

func NewWaitingState(g *Game) *WaitingState {
	return &WaitingState{
		game: g,
		mtx:  new(sync.Mutex),
	}
}

func (s *WaitingState) AddPlayer(player *domain.Player) error {
	s.mtx.Lock()
	s.game.players = append(s.game.players, player)
	if len(s.game.players) == s.game.maxPlayers {
		s.game.state = NewPlayingState(s.game)
		s.game.scoremanager = scoremanager.NewScoreManager(s.game.players)
		s.game.combinations = make([]PlayerChoice, 0, len(s.game.combinations))

		err := s.game.observer.Publish("room_started", struct{}{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to publish event \"room_started\"")
		}
	}
	s.mtx.Unlock()

	log.Info().Msgf("Player %s added to the room", player.GetID())

	return nil
}

func (s *WaitingState) Choose(choice *PlayerChoice) error {
	return ErrGameNotStarted
}

func (s *WaitingState) GetLeader() (string, error) {
	return "", ErrGameNotStarted
}

func (s *WaitingState) GetPlayerScore(name string) (uint64, error) {
	return 0, ErrGameNotStarted
}

func (s *WaitingState) IncPlayerScore(name string) error {
	return ErrGameNotStarted
}
