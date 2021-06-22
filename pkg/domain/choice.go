package domain

import (
	"errors"
	"strings"
)

type Choice int

const (
	ROCK     Choice = 0
	PAPER    Choice = 1
	SCISSORS Choice = 2
)

func (this Choice) Compare(another Choice) int {
	const choicesCount = 3
	if this == another {
		return 0
	}
	thisInt := int(this)
	anotherInt := int(another)
	winningChoice := (thisInt + 1) % choicesCount
	if winningChoice == anotherInt {
		return -1
	}
	return 1
}

func GetChoiceByName(name string) (Choice, error) {
	switch strings.ToLower(name) {
	case "rock":
		return ROCK, nil
	case "paper":
		return PAPER, nil
	case "scissors":
		return SCISSORS, nil
	}
	return 0, errors.New("Not such choice")
}
