package room

import (
	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/winner-definer"
	"github.com/rs/zerolog/log"
)

type PlayingState struct {
	room *Room
}

func NewPlayingState(r *Room) *PlayingState {
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
		s.room.observer.Publish("winners", winners)
		s.room.combinations = make([]PlayerChoice, 0, s.room.config.MaxPlayerCount)
	}
	return nil
}
