package winner_definer

import (
	"github.com/pavel/PSR/pkg/domain"
)

type WinnerDefiner struct{}

type playerChoice struct {
	playerId string
	input    domain.Choice
}

func (wd *WinnerDefiner) GetWinners(playersChoices []playerChoice) []string {
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
		count[int(choice.input)]++
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
	switch missing {
	case -1: // draw, all players lost
		return nil
	case 0: // Paper is missing
		for _, choice := range playersChoices {
			if choice.input.Compare(domain.ROCK) == 0 {
				ret = append(ret, choice.playerId)
			}
		}
	case 1: // Scissors is missing
		ret := []string{}
		for _, choice := range playersChoices {
			if choice.input.Compare(domain.PAPER) == 0 {
				ret = append(ret, choice.playerId)
			}
		}
	case 2: // Rock is missing
		ret := []string{}
		for _, choice := range playersChoices {
			if choice.input.Compare(domain.SCISSORS) == 0 {
				ret = append(ret, choice.playerId)
			}
		}
	}
	return ret
}
