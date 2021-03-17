package winner_definer

import (
	"testing"
	"time"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/room"
)

func TestGetWinners(t *testing.T) {
	room := room.NewRoom(room.RoomConfig{
		StepTimeout:    5 * time.Second,
		MaxPlayerCount: 3,
		MaxScore:       5,
		OnlyComputer:   false,
	})
	players := []*domain.Player{
		domain.NewPlayer("1"),
		domain.NewPlayer("2"),
		domain.NewPlayer("3"),
	}
	for _, player := range players {
		room.AddPlayer(player)
	}
	room.AddPlayer(domain.NewPlayer("4"))
	time.Sleep(2 * time.Second)
}
