package room

import (
	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/winner-definer"
)

type State interface {
	AddPlayer(*domain.Player) error
	Choose(*PlayerChoice) error
}
