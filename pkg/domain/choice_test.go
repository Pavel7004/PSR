package domain

import "testing"

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
		{
			"test scissors vs rock",
			SCISSORS,
			args{
				another: ROCK,
			},
			-1,
		},
		{
			"test scissors vs paper",
			SCISSORS,
			args{
				another: PAPER,
			},
			1,
		},
		{
			"test scissors vs scissors",
			SCISSORS,
			args{
				another: SCISSORS,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Compare(tt.args.another); got != tt.want {
				t.Errorf("Choice.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
