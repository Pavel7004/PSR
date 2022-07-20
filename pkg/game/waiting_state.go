package room

import (
	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/score-manager"
	. "github.com/pavel/PSR/pkg/winner-definer"
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
	log.Info().Msgf("Player %s added to the room", player.GetID())

	if len(s.room.players) == s.room.Config.MaxPlayerCount {
		s.room.state = NewPlayingState(s.room)
		s.room.scoremanager = NewScoreManager(s.room.players)
		log.Info().Msg("Room started")

		err := s.room.observer.Publish("room_started", struct{}{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to publish event \"room_started\"")
			return err
		}
	}

	s.room.stepMtx.Unlock()
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
