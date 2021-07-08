package scoremanager

import "github.com/pavel/PSR/pkg/domain"

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

func (sm *ScoreManager) GetPlayerScore(name string) int {
	return sm.playersScores[name]
}

func (sm *ScoreManager) IncrementPlayerScore(name string) {
	sm.playersScores[name]++
}

func (sm *ScoreManager) ResetPlayersScores() {
	for name := range sm.playersScores {
		sm.playersScores[name] = 0
	}
}
