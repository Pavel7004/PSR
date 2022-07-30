package game

import (
	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/server/winner-definer"
	"github.com/rs/zerolog/log"
)

type WaitingState struct {
	room *Game
}

func NewWaitingState(r *Game) *WaitingState {
	return &WaitingState{
		room: r,
	}
}

func (s *WaitingState) AddPlayer(player *domain.Player) error {
	s.room.stepMtx.Lock()
	s.room.players = append(s.room.players, player)
	s.room.stepMtx.Unlock()

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
