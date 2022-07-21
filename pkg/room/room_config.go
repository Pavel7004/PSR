package room

import "time"

type RoomConfig struct {
	Name           string
	StepTimeout    time.Duration
	MaxPlayerCount int
	MaxScore       uint64
	OnlyComputer   bool
}
