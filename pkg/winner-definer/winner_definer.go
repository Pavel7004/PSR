package winner_definer

import (
	"github.com/pavel/PSR/pkg/domain"
)

type WinnerDefiner struct{}

type playerChoise struct {
	playerId string
	choise   domain.Choice
}

func (wd *WinnerDefiner) GetWinners(playersChoises []playerChoise) {
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
	for _, choise := range playersChoises {

	}
}
