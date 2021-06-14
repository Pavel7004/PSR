package winner_definer

import (
	"github.com/pavel/PSR/pkg/domain"
)

type WinnerDefiner struct{}

func (wd *WinnerDefiner) GetWinners(playersChoices []PlayerChoice) []string {
	const numberOfChoices = 3
	count := [numberOfChoices]int{0}
	for _, choice := range playersChoices {
		count[int(choice.Input)]++
	}
	missing := -1
	missingKol := 0
	for i := 0; i < numberOfChoices; i++ {
		if count[i] == 0 {
			missing = i
			missingKol++
		}
	}
	if missingKol != 1 {
		return nil
	}
	var another domain.Choice
	switch missing {
	case 0:
		another = domain.SCISSORS
	case 1:
		another = domain.ROCK
	case 2:
		another = domain.PAPER
	}
	ret := []string{}
	for _, choice := range playersChoices {
		if choice.Input.Compare(another) == 0 {
			ret = append(ret, choice.PlayerID)
		}
	}
	return ret
}
