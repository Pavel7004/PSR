package game

import (
	"github.com/pavel/PSR/pkg/domain"
)

type State interface {
	AddPlayer(*domain.Player) error
	Choose(id string, choice domain.Choice) error
	GetLeader() (string, error)
	GetPlayerScore(string) (uint64, error)
	IncPlayerScore(string) error
}
