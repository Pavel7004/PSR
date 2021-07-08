package scoremanager

import (
	"errors"

	"github.com/pavel/PSR/pkg/domain"
)

var (
	ErrPlayerNotFound = errors.New("Player isn't present in current score_manager")
)

type ScoreManager struct {
	playersScores map[string]int
}

func NewScoreManager(players []*domain.Player) *ScoreManager {
	newScores := make(map[string]int, len(players))
	for _, name := range players {
		newScores[name.GetID()] = 0
	}
	return &ScoreManager{
		playersScores: newScores,
	}
}

func (sm *ScoreManager) GetPlayerScore(name string) (int, error) {
	val, err := sm.playersScores[name]
	if !err {
		return 0, ErrPlayerNotFound
	}
	return val, nil
}

func (sm *ScoreManager) IncrementPlayerScore(name string) error {
	val, err := sm.playersScores[name]
	if !err {
		return ErrPlayerNotFound
	}
	sm.playersScores[name] = val + 1
	return nil
}

func (sm *ScoreManager) GetMaxScore() (string, int) {
	maxScore := -1
	maxName := ""
	for name, score := range sm.playersScores {
		if score > maxScore {
			maxName = name
			maxScore = score
		}
	}
	return maxName, maxScore
}

func (sm *ScoreManager) ResetPlayersScores() {
	for name := range sm.playersScores {
		sm.playersScores[name] = 0
	}
}
