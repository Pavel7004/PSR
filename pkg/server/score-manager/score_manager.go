package scoremanager

import (
	"errors"

	"github.com/pavel/PSR/pkg/domain"
)

var (
	ErrPlayerNotFound = errors.New("Player isn't present in current score_manager")
)

type ScoreManager struct {
	playersScores map[string]uint64
}

func NewScoreManager(players []*domain.Player) *ScoreManager {
	newScores := make(map[string]uint64, len(players))
	for _, name := range players {
		newScores[name.GetID()] = 0
	}
	return &ScoreManager{
		playersScores: newScores,
	}
}

func (sm *ScoreManager) GetPlayerScore(name string) (uint64, error) {
	val, ok := sm.playersScores[name]
	if !ok {
		return 0, ErrPlayerNotFound
	}
	return val, nil
}

func (sm *ScoreManager) IncrementPlayerScore(name string) error {
	val, ok := sm.playersScores[name]
	if !ok {
		return ErrPlayerNotFound
	}
	sm.playersScores[name] = val + 1
	return nil
}

func (sm *ScoreManager) GetLeadingPlayerName() string {
	var (
		maxScore uint64
		maxName  string
	)
	for name, score := range sm.playersScores {
		if score > maxScore {
			maxName = name
			maxScore = score
		}
	}
	return maxName
}

func (sm *ScoreManager) ResetPlayersScores() {
	for name := range sm.playersScores {
		sm.playersScores[name] = 0
	}
}
