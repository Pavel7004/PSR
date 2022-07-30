package game

import (
	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/server/winner-definer"
	"github.com/rs/zerolog/log"
)

type PlayingState struct {
	room *Game
}

func NewPlayingState(r *Game) *PlayingState {
	return &PlayingState{
		room: r,
	}
}

func (s *PlayingState) AddPlayer(player *domain.Player) error {
	return ErrGameAlreadyStarted
}

func (s *PlayingState) Choose(choice *PlayerChoice) error {
	if !s.room.HasPlayer(choice.PlayerID) {
		return ErrPlayerNotPresent
	}

	s.room.combinations = append(s.room.combinations, *choice)
	if len(s.room.combinations) == len(s.room.players) {
		winners := s.room.winnerDefiner.GetWinners(s.room.combinations)
		log.Info().Msgf("Winners: %v", winners)

		err := s.room.observer.Publish("winners", winners)
		s.room.combinations = make([]PlayerChoice, 0, len(s.room.combinations))

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
