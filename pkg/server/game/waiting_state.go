package game

import (
	"sync"

	"github.com/pavel/PSR/pkg/domain"

	scoremanager "github.com/pavel/PSR/pkg/server/score-manager"
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

	last := len(s.game.players) == s.game.maxPlayers
	s.mtx.Unlock()

	if last {
		s.game.state = NewPlayingState(s.game)
		s.game.scoremanager = scoremanager.NewScoreManager(s.game.players)
		s.game.combinations = make(map[string]domain.Choice, len(s.game.players))

		err := s.game.observer.Publish("room_started", struct{}{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to publish event \"room_started\"")
		}
	}

	log.Info().Msgf("Player %q added to the room", player.GetID())

	return nil
}

func (s *WaitingState) Choose(id string, choice domain.Choice) error {
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
