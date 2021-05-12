package room

import "github.com/pavel/PSR/pkg/domain"

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
