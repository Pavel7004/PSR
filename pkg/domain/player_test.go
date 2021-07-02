package domain

import (
	"reflect"
	"testing"
)

func TestNewPlayer(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want *Player
	}{
		{
			name: "Generate new player",
			args: args{"Test_player"},
			want: &Player{id: "Test_player", score: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlayer(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_GetID(t *testing.T) {
	type fields struct {
		ID string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Get player ID",
			fields: fields{"Test_player_ID"},
			want:   "Test_player_ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				id:    tt.fields.ID,
				score: 0,
			}
			if got := p.GetID(); got != tt.want {
				t.Errorf("Player.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_GetScore(t *testing.T) {
	type fields struct {
		id    string
		score int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "get zero",
			fields: fields{id: "testPlayer", score: 0},
			want:   0,
		},
		{
			name:   "get non-zero",
			fields: fields{id: "testPlayer", score: 7},
			want:   7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				id:    tt.fields.id,
				score: tt.fields.score,
			}
			if got := p.GetScore(); got != tt.want {
				t.Errorf("Player.GetScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_IncrementScore(t *testing.T) {
	type fields struct {
		id    string
		score int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Player
	}{
		{
			name:   "increment zero",
			fields: fields{id: "testPlayer", score: 0},
			want:   &Player{id: "testPlayer", score: 1},
		},
		{
			name:   "increment non-zero",
			fields: fields{id: "testPlayer", score: 6},
			want:   &Player{id: "testPlayer", score: 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				id:    tt.fields.id,
				score: tt.fields.score,
			}
			p.IncrementScore()
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("Player.IncrementScore() = %v, want %v", p.score, tt.want.score)
			}
		})
	}
}

func TestPlayer_ResetScore(t *testing.T) {
	type fields struct {
		id    string
		score int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Player
	}{
		{
			name:   "Reset zero score",
			fields: fields{id: "testPlayer", score: 0},
			want:   &Player{id: "testPlayer", score: 0},
		},
		{
			name:   "Reset non-zero score",
			fields: fields{id: "testPlayer", score: 8},
			want:   &Player{id: "testPlayer", score: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				id:    tt.fields.id,
				score: tt.fields.score,
			}
			p.ResetScore()
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("Player.ResetScore() = %v, want %v", p.score, tt.want.score)
			}
		})
	}
}
