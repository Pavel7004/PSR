package scoremanager

import (
	"reflect"
	"testing"

	"github.com/pavel/PSR/pkg/domain"
)

func TestNewScoreManager(t *testing.T) {
	type args struct {
		players []*domain.Player
	}
	tests := []struct {
		name string
		args args
		want *ScoreManager
	}{
		{
			name: "Generate new score manager",
			args: args{players: []*domain.Player{domain.NewPlayer("test")}},
			want: &ScoreManager{playersScores: map[string]int{"test": 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScoreManager(tt.args.players); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScoreManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScoreManager_GetPlayerScore(t *testing.T) {
	type fields struct {
		playersScores map[string]int
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Player is not present",
			fields:  fields{playersScores: map[string]int{"test": 3}},
			args:    args{"NotExist"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Get player score",
			fields:  fields{playersScores: map[string]int{"test": 4}},
			args:    args{name: "test"},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &ScoreManager{
				playersScores: tt.fields.playersScores,
			}
			got, err := sm.GetPlayerScore(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScoreManager.GetPlayerScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ScoreManager.GetPlayerScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScoreManager_IncrementPlayerScore(t *testing.T) {
	type fields struct {
		playersScores map[string]int
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal int
		wantErr bool
	}{
		{
			name:    "Player isn't present",
			fields:  fields{playersScores: map[string]int{"test": 1}},
			args:    args{name: "NotExist"},
			wantVal: 0,
			wantErr: true,
		},
		{
			name:    "Increment existing player value",
			fields:  fields{playersScores: map[string]int{"test": 0}},
			args:    args{"test"},
			wantVal: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &ScoreManager{
				playersScores: tt.fields.playersScores,
			}
			if err := sm.IncrementPlayerScore(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ScoreManager.IncrementPlayerScore() error = %v, wantErr %v", err, tt.wantErr)
			}
			if val := sm.playersScores[tt.args.name]; val != tt.wantVal {
				t.Errorf("ScoreManager.IncrementPlayerScore() val = %v, wantVal = %v", val, tt.wantVal)
			}
		})
	}
}

func TestScoreManager_ResetPlayersScores(t *testing.T) {
	type fields struct {
		playersScores map[string]int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Reset all scores",
			fields: fields{
				playersScores: map[string]int{
					"test1": 3,
					"test2": 6,
					"test3": 10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &ScoreManager{
				playersScores: tt.fields.playersScores,
			}
			sm.ResetPlayersScores()
			checkFail := false
			for _, val := range sm.playersScores {
				checkFail = checkFail || (val != 0)
			}
			if checkFail {
				t.Errorf("ScoreManager.IncrementPlayerScore() fail to reset players scores")
			}
		})
	}
}
