package domain

import "testing"

// func TestCompare(t *testing.T) {
// 	got := ROCK.Compare(ROCK)
// 	if got != 0 {
// 		t.Errorf("ROCK != ROCK (returned %d, expected 0)", got)
// 	}
// 	got = ROCK.Compare(PAPER)
// 	if got != -1 {
// 		t.Errorf("ROCK >= PAPER (returned %d, expected -1)", got)
// 	}
// 	got = ROCK.Compare(SCISSORS)
// 	if got != 1 {
// 		t.Errorf("ROCK <= SCISSORS (returned %d, expected 1)", got)
// 	}
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
		// {
		// 	name: "test rock",
		// 	this: ROCK,
		// 	args: args{
		// 		another: PAPER,
		// 	},
		// 	want: -1,
		// },

		{
			"test rock",
			ROCK,
			args{
				another: PAPER,
			},
			-1,
		},

		// {
		// 	name: "test rock",
		// 	this: ROCK,
		// 	args: args{
		// 		another: PAPER,
		// 	},
		// 	want: -1,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Compare(tt.args.another); got != tt.want {
				t.Errorf("Choice.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
