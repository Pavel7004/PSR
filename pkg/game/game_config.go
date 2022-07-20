package room

import "time"

type GameConfig struct {
	Name           string
	StepTimeout    time.Duration
	MaxPlayerCount int
	MaxScore       uint64
	OnlyComputer   bool
}
