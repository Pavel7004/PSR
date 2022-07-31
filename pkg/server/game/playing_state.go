package game

import (
	"sync"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/rs/zerolog/log"
)

type PlayingState struct {
	room *Game
	mtx  *sync.Mutex
}

func NewPlayingState(r *Game) *PlayingState {
	return &PlayingState{
		room: r,
		mtx:  new(sync.Mutex),
	}
}

func (s *PlayingState) AddPlayer(player *domain.Player) error {
	return ErrGameAlreadyStarted
}

func (s *PlayingState) Choose(id string, choice domain.Choice) error {
	if !s.room.HasPlayer(id) {
		return ErrPlayerNotPresent
	}

	s.mtx.Lock()
	s.room.combinations[id] = choice

	last := len(s.room.combinations) == len(s.room.players)
	s.mtx.Unlock()

	if last {
		winners := s.room.winnerDefiner.GetWinners(s.room.combinations)
		log.Info().Msgf("Winners: %v", winners)

		err := s.room.observer.Publish("winners", winners)
		s.room.combinations = make(map[string]domain.Choice, len(s.room.players))

		if err != nil {
			log.Error().Err(err).Msg("Failed to publish event \"winners\"")
			return err
		}
	}

	return nil
}

func (s *PlayingState) GetLeader() (string, error) {
	return s.room.scoremanager.GetLeadingPlayerName(), nil
}

func (s *PlayingState) GetPlayerScore(name string) (uint64, error) {
	score, err := s.room.scoremanager.GetPlayerScore(name)
	if err != nil {
		return 0, ErrPlayerNotPresent
	}

	return score, nil
}

func (s *PlayingState) IncPlayerScore(name string) error {
	if err := s.room.scoremanager.IncrementPlayerScore(name); err != nil {
		log.Error().Err(err).Msgf("Player \"%s\" not found by score_manager", name)
		return ErrPlayerNotPresent
	}

	return nil
}
