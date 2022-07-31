package game

import (
	"sync"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/rs/zerolog/log"
)

type PlayingState struct {
	game *Game
	mtx  *sync.Mutex
}

func NewPlayingState(r *Game) *PlayingState {
	return &PlayingState{
		game: r,
		mtx:  new(sync.Mutex),
	}
}

func (s *PlayingState) AddPlayer(player *domain.Player) error {
	return ErrGameAlreadyStarted
}

func (s *PlayingState) Choose(id string, choice domain.Choice) error {
	if !s.game.HasPlayer(id) {
		return ErrPlayerNotPresent
	}

	s.mtx.Lock()
	s.game.combinations[id] = choice

	last := len(s.game.combinations) == len(s.game.players)
	s.mtx.Unlock()

	if last {
		winners := s.game.winnerDefiner.GetWinners(s.game.combinations)
		log.Info().Msgf("Winners: %v", winners)

		err := s.game.observer.Publish("winners", winners)
		s.game.combinations = make(map[string]domain.Choice, len(s.game.players))

		if err != nil {
			log.Error().Err(err).Msg("Failed to publish event \"winners\"")
			return err
		}
	}

	return nil
}

func (s *PlayingState) GetLeader() (string, error) {
	return s.game.scoremanager.GetLeadingPlayerName(), nil
}

func (s *PlayingState) GetPlayerScore(name string) (uint64, error) {
	score, err := s.game.scoremanager.GetPlayerScore(name)
	if err != nil {
		return 0, ErrPlayerNotPresent
	}

	return score, nil
}

func (s *PlayingState) IncPlayerScore(id string) error {
	if err := s.game.scoremanager.IncrementPlayerScore(id); err != nil {
		log.Error().Err(err).Msgf("Player %q not found by score_manager", id)
		return ErrPlayerNotPresent
	}

	return nil
}
