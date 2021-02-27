package domain

type Choice int

const (
	ROCK     Choice = 0
	PAPER    Choice = 1
	SCISSORS Choice = 2
)

func (this Choice) Compare(another Choice) int {
	if this == ROCK && another == SCISSORS ||
		this == SCISSORS && another == ROCK {

	}
	if this == another {
		return 0
	}
	if this > another {
		return 1
	}
	return -1
}

// func main() {
// 	PAPER.Compare(SCISSORS)
// 	choices :=
// 	choices = []Choise {
// 		PAPER,
// 		PAPER,
// 		ROCK,
// 	}
// 	max(choices) => PAPER
// }
