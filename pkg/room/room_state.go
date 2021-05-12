package room

import "github.com/pavel/PSR/pkg/domain"

type State interface {
	AddPlayer(*domain.Player) error
}
