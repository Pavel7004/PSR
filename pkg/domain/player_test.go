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
			want: &Player{id: "Test_player"},
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
				id: tt.fields.ID,
			}
			if got := p.GetID(); got != tt.want {
				t.Errorf("Player.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}
