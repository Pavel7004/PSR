package room

import (
	"sync"
	"testing"
	"time"

	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/winner-definer"
)

func TestRoom_IsActive(t *testing.T) {
	type fields struct {
		config        RoomConfig
		players       []*domain.Player
		combinations  []PlayerChoice
		active        bool
		stopCh        chan struct{}
		chooseCh      chan PlayerChoice
		stepMtx       *sync.Mutex
		winnerDefiner *WinnerDefiner
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test Active",
			fields: fields{
				RoomConfig{},
				nil,
				[]PlayerChoice{},
				true,
				nil,
				nil,
				nil,
				&WinnerDefiner{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				config:        tt.fields.config,
				players:       tt.fields.players,
				combinations:  tt.fields.combinations,
				active:        tt.fields.active,
				stopCh:        tt.fields.stopCh,
				chooseCh:      tt.fields.chooseCh,
				stepMtx:       tt.fields.stepMtx,
				winnerDefiner: tt.fields.winnerDefiner,
			}
			if got := room.IsActive(); got != tt.want {
				t.Errorf("Room.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoom_AddPlayer(t *testing.T) {
	type fields struct {
		config        RoomConfig
		players       []*domain.Player
		combinations  []PlayerChoice
		active        bool
		stopCh        chan struct{}
		chooseCh      chan PlayerChoice
		stepMtx       *sync.Mutex
		winnerDefiner *WinnerDefiner
	}
	type args struct {
		player *domain.Player
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test adding player in active game",
			fields: fields{
				RoomConfig{
					5 * time.Minute,
					5,
					5,
					false,
				},
				make([]*domain.Player, 0, 5),
				[]PlayerChoice{},
				true,
				nil,
				nil,
				new(sync.Mutex),
				&WinnerDefiner{},
			},
			args: args{
				player: &domain.Player{
					ID: "TestPlayer",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				config:        tt.fields.config,
				players:       tt.fields.players,
				combinations:  tt.fields.combinations,
				active:        tt.fields.active,
				stopCh:        tt.fields.stopCh,
				chooseCh:      tt.fields.chooseCh,
				stepMtx:       tt.fields.stepMtx,
				winnerDefiner: tt.fields.winnerDefiner,
			}
			if err := room.AddPlayer(tt.args.player); (err != nil) != tt.wantErr {
				t.Errorf("Room.AddPlayer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoom_Run(t *testing.T) {
	type fields struct {
		config        RoomConfig
		players       []*domain.Player
		combinations  []PlayerChoice
		active        bool
		stopCh        chan struct{}
		chooseCh      chan PlayerChoice
		stepMtx       *sync.Mutex
		winnerDefiner *WinnerDefiner
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				config:        tt.fields.config,
				players:       tt.fields.players,
				combinations:  tt.fields.combinations,
				active:        tt.fields.active,
				stopCh:        tt.fields.stopCh,
				chooseCh:      tt.fields.chooseCh,
				stepMtx:       tt.fields.stepMtx,
				winnerDefiner: tt.fields.winnerDefiner,
			}
			room.Run()
		})
	}
}

func TestRoom_Choose(t *testing.T) {
	type fields struct {
		config        RoomConfig
		players       []*domain.Player
		combinations  []PlayerChoice
		active        bool
		stopCh        chan struct{}
		chooseCh      chan PlayerChoice
		stepMtx       *sync.Mutex
		winnerDefiner *WinnerDefiner
	}
	type args struct {
		choice PlayerChoice
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				config:        tt.fields.config,
				players:       tt.fields.players,
				combinations:  tt.fields.combinations,
				active:        tt.fields.active,
				stopCh:        tt.fields.stopCh,
				chooseCh:      tt.fields.chooseCh,
				stepMtx:       tt.fields.stepMtx,
				winnerDefiner: tt.fields.winnerDefiner,
			}
			room.Choose(tt.args.choice)
		})
	}
}
