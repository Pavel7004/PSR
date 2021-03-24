package domain

import "testing"

// 	got = PAPER.Compare(ROCK)
// 	if got != 1 {
// 		t.Errorf("PAPER <= ROCK (returned %d, expected 1)", got)
// 	}
// 	got = PAPER.Compare(PAPER)
// 	if got != 0 {
// 		t.Errorf("PAPER != PAPER (returned %d, expected 0)", got)
// 	}
// 	got = PAPER.Compare(SCISSORS)
// 	if got != -1 {
// 		t.Errorf("PAPER >= SCISSORS (returned %d, expected -1)", got)
// 	}
// 	got = SCISSORS.Compare(ROCK)
// 	if got != -1 {
// 		t.Errorf("SCISSORS >= ROCK (returned %d, expected -1)", got)
// 	}
// 	got = SCISSORS.Compare(PAPER)
// 	if got != 1 {
// 		t.Errorf("SCISSORS <= PAPER (returned %d, expected 1)", got)
// 	}
// 	got = SCISSORS.Compare(SCISSORS)
// 	if got != 0 {
// 		t.Errorf("SCISSORS != SCISSORS (returned %d, expected 0)", got)
// 	}
// }

func TestChoice_Compare(t *testing.T) {
	type args struct {
		another Choice
	}
	tests := []struct {
		name string
		this Choice
		args args
		want int
	}{
		{
			"test rock vs paper",
			ROCK,
			args{
				another: PAPER,
			},
			-1,
		},
		{
			"test rock vs scissors",
			ROCK,
			args{
				another: SCISSORS,
			},
			1,
		},
		{
			"test rock vs rock",
			ROCK,
			args{
				another: ROCK,
			},
			0,
		},
		{
			"test paper vs rock",
			PAPER,
			args{
				another: ROCK,
			},
			1,
		},
		{
			"test paper vs paper",
			PAPER,
			args{
				another: PAPER,
			},
			0,
		},
		{
			"test paper vs scissors",
			PAPER,
			args{
				another: SCISSORS,
			},
			-1,
		},
		// TODO: add scissors tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Compare(tt.args.another); got != tt.want {
				t.Errorf("Choice.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
