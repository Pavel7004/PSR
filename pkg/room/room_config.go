package room

import "time"

type RoomConfig struct {
	StepTimeout    time.Duration
	MaxPlayerCount int
	MaxScore       int
	OnlyComputer   bool
}
