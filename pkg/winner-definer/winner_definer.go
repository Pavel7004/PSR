package winner_definer

import (
	"github.com/pavel/PSR/pkg/domain"
)

type WinnerDefiner struct{}

type PlayerChoice struct {
	PlayerID string
	Input    domain.Choice
}

func (wd *WinnerDefiner) GetWinners(playersChoices []PlayerChoice) []string {
	// => [
	// 	{
	// 		id: "Vasya",
	// 		choice: domain.PAPER,
	// 	},
	// 	{
	// 		id: "Petya",
	// 		choice: domain.PAPER,
	// 	},
	// 	{
	// 		id: "Masha",
	// 		choice: domain.ROCK,
	// 	},
	// ]
	// max(...) => the strongest choice
	// => ["Vasya", "Petya"]
	count := [3]int{0}
	for _, choice := range playersChoices {
		count[int(choice.Input)]++
	}
	missing := -1
	missingKol := 0
	for i := 0; i < 3; i++ {
		if count[i] == 0 {
			missing = i
			missingKol++
		}
	}
	if missingKol > 1 {
		return nil
	}
	ret := []string{}
	var another domain.Choice
	switch missing {
	case -1: // draw, all players lost
		return nil
	case 0:
		another = domain.SCISSORS
	case 1:
		another = domain.SCISSORS
	case 2:
		another = domain.PAPER
	}
	for _, choice := range playersChoices {
		if choice.Input.Compare(another) == 0 {
			ret = append(ret, choice.PlayerID)
		}
	}
	return ret
}
