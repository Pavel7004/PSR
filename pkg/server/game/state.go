package game

import (
	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/server/winner-definer"
)

type State interface {
	AddPlayer(*domain.Player) error
	Choose(*PlayerChoice) error
	GetLeader() (string, error)
	GetPlayerScore(string) (uint64, error)
	IncPlayerScore(string) error
}
