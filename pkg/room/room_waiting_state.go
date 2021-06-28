package room

import (
	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/winner-definer"
	"github.com/rs/zerolog/log"
)

type WaitingState struct {
	room *Room
}

func NewWaitingState(r *Room) *WaitingState {
	return &WaitingState{
		room: r,
	}
}

func (s *WaitingState) AddPlayer(player *domain.Player) error {
	s.room.stepMtx.Lock()
	s.room.players = append(s.room.players, player)
	log.Info().Msgf("Player %s added to the room", player.GetID())
	if len(s.room.players) == s.room.config.MaxPlayerCount {
		s.room.state = NewPlayingState(s.room)
		log.Info().Msg("Room started")
		s.room.observer.Publish("room_started", struct{}{})
	}
	s.room.stepMtx.Unlock()
	return nil
}

func (s *WaitingState) Choose(choice *PlayerChoice) error {
	return ErrGameNotStarted
}
