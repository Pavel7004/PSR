package room

import (
	"github.com/pavel/PSR/pkg/domain"
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
	log.Info().Msgf("Player %s added to the room", player.ID)
	if len(s.room.players) == s.room.config.MaxPlayerCount {
		go s.room.Run()
	}
	s.room.stepMtx.Unlock()
	return nil
}
