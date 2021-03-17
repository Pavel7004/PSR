package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/room"
	winner_definer "github.com/pavel/PSR/pkg/winner-definer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	r := room.NewRoom(
		room.RoomConfig{
			5 * time.Second,
			2,
			5,
			false,
		},
	)
	players := []string{"Test1", "Test2"}
	for _, name := range players {
		r.AddPlayer(domain.NewPlayer(name))
	}

	for _, name := range players {
		fmt.Printf("%s: ", name)
		var input domain.Choice
		fmt.Scanf("%d", &input)
		r.Choose(winner_definer.PlayerChoice{
			PlayerID: name,
			Input:    input,
		})
	}

	time.Sleep(2 * time.Second)
	// wd := winner_definer.WinnerDefiner{}
	// ret := wd.GetWinners(
	// 	[]winner_definer.PlayerChoice{
	// 		{
	// 			PlayerID: "Test1",
	// 			Input:    domain.PAPER,
	// 		},
	// 		{
	// 			PlayerID: "Test2",
	// 			Input:    domain.PAPER,
	// 		},
	// 		{
	// 			PlayerID: "Test3",
	// 			Input:    domain.PAPER,
	// 		},
	// 	},
	// )
	// fmt.Println(ret)
}
