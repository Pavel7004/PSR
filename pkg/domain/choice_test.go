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
			name: "test rock vs paper",
			this: ROCK,
			args: args{another: PAPER},
			want: -1,
		},
		{
			name: "test rock vs scissors",
			this: ROCK,
			args: args{another: SCISSORS},
			want: 1,
		},
		{
			name: "test rock vs rock",
			this: ROCK,
			args: args{another: ROCK},
			want: 0,
		},
		{
			name: "test paper vs rock",
			this: PAPER,
			args: args{another: ROCK},
			want: 1,
		},
		{
			name: "test paper vs paper",
			this: PAPER,
			args: args{another: PAPER},
			want: 0,
		},
		{
			name: "test paper vs scissors",
			this: PAPER,
			args: args{another: SCISSORS},
			want: -1,
		},
		{
			name: "test scissors vs rock",
			this: SCISSORS,
			args: args{another: ROCK},
			want: -1,
		},
		{
			name: "test scissors vs paper",
			this: SCISSORS,
			args: args{another: PAPER},
			want: 1,
		},
		{
			name: "test scissors vs scissors",
			this: SCISSORS,
			args: args{another: SCISSORS},
			want: 0,
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
