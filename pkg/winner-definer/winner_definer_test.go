package winner_definer

import (
	"reflect"
	"testing"

	"github.com/pavel/PSR/pkg/domain"
)

func TestWinnerDefiner_GetWinners(t *testing.T) {
	type args struct {
		playersChoices []PlayerChoice
	}
	tests := []struct {
		name string
		wd   *WinnerDefiner
		args args
		want []string
	}{
		{
			name: "Rocks vs Paper",
			wd:   &WinnerDefiner{},
			args: args{
				[]PlayerChoice{
					{
						"Player 1",
						domain.PAPER,
					},
					{
						"Player 2",
						domain.ROCK,
					},
					{
						"Player 3",
						domain.ROCK,
					},
				},
			},
			want: []string{"Player 1"},
		},
		{
			name: "Scissors vs papers",
			wd:   &WinnerDefiner{},
			args: args{
				[]PlayerChoice{
					{
						"Player 1",
						domain.PAPER,
					},
					{
						"Player 2",
						domain.SCISSORS,
					},
					{
						"Player 3",
						domain.PAPER,
					},
					{
						"Player 4",
						domain.SCISSORS,
					},
				},
			},
			want: []string{"Player 2", "Player 4"},
		},
		{
			name: "Rock vs scissors",
			wd:   &WinnerDefiner{},
			args: args{
				[]PlayerChoice{
					{
						"Player 1",
						domain.ROCK,
					},
					{
						"Player 2",
						domain.SCISSORS,
					},
					{
						"Player 3",
						domain.SCISSORS,
					},
					{
						"Player 4",
						domain.SCISSORS,
					},
				},
			},
			want: []string{"Player 1"},
		},
		{
			name: "Mixed",
			wd:   &WinnerDefiner{},
			args: args{
				[]PlayerChoice{
					{
						"Player 1",
						domain.PAPER,
					},
					{
						"Player 2",
						domain.SCISSORS,
					},
					{
						"Player 3",
						domain.ROCK,
					},
					{
						"Player 4",
						domain.SCISSORS,
					},
				},
			},
			want: nil,
		},
		{
			name: "The same",
			wd:   &WinnerDefiner{},
			args: args{
				[]PlayerChoice{
					{
						"Player 1",
						domain.SCISSORS,
					},
					{
						"Player 2",
						domain.SCISSORS,
					},
					{
						"Player 3",
						domain.SCISSORS,
					},
					{
						"Player 4",
						domain.SCISSORS,
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := &WinnerDefiner{}
			if got := wd.GetWinners(tt.args.playersChoices); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WinnerDefiner.GetWinners() = %v, want %v", got, tt.want)
			}
		})
	}
}
