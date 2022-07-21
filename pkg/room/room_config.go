package room

import "time"

type RoomConfig struct {
	Name           string
	RoundTimeout   time.Duration
	MaxPlayerCount int
	MaxScore       uint64
}
