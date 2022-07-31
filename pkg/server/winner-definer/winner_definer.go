package winner_definer

import "github.com/pavel/PSR/pkg/domain"

type WinnerDefiner struct{}

func (wd *WinnerDefiner) GetWinners(playersChoices map[string]domain.Choice) []string {
	count := map[domain.Choice]int{
		domain.ROCK:     0,
		domain.PAPER:    0,
		domain.SCISSORS: 0,
	}
	for _, choice := range playersChoices {
		count[choice]++
	}
	missing := -1
	for key, value := range count {
		if value == 0 {
			if missing != -1 {
				return []string{}
			}
			missing = int(key)
		}
	}
	if missing == -1 {
		return []string{}
	}
	winningPiece := domain.Choice((missing + 2) % 3)
	winners := []string{}
	for id, choice := range playersChoices {
		if choice.Compare(winningPiece) == 0 {
			winners = append(winners, id)
		}
	}
	return winners
}
