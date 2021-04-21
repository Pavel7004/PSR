package room

import (
	"reflect"
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
				nil,
			},
			want: true,
		},
		{
			name: "Test Unactive",
			fields: fields{
				RoomConfig{},
				nil,
				[]PlayerChoice{},
				false,
				nil,
				nil,
				nil,
				nil,
			},
			want: false,
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
	testPlayers := map[string][]*domain.Player{
		"Test adding player in active game": {},
		"Adding player": {
			&domain.Player{
				ID: "Player 1",
			},
		},
		"Adding last player": {
			&domain.Player{
				ID: "Player 1",
			},
			&domain.Player{
				ID: "Player 2",
			},
			&domain.Player{
				ID: "Player 3",
			},
			&domain.Player{
				ID: "Player 4",
			},
			&domain.Player{
				ID: "Player 5",
			},
		},
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		expectedPlayers     []*domain.Player
		expectedRoomStarted bool
		wantErr             bool
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
					ID: "testPlayer",
				},
			},
			expectedPlayers:     testPlayers["Test adding player in active game"],
			expectedRoomStarted: true,
			wantErr:             true,
		},
		{
			name: "Adding player",
			fields: fields{
				RoomConfig{
					5 * time.Minute,
					5,
					5,
					false,
				},
				make([]*domain.Player, 0, 5),
				[]PlayerChoice{},
				false,
				nil,
				nil,
				new(sync.Mutex),
				&WinnerDefiner{},
			},
			args: args{
				player: testPlayers["Adding player"][0],
			},
			expectedPlayers:     testPlayers["Adding player"],
			expectedRoomStarted: false,
			wantErr:             false,
		},
		{
			name: "Adding last player",
			fields: fields{
				RoomConfig{
					5 * time.Minute,
					5,
					5,
					false,
				},
				[]*domain.Player{
					testPlayers["Adding last player"][0],
					testPlayers["Adding last player"][1],
					testPlayers["Adding last player"][2],
					testPlayers["Adding last player"][3],
				},
				[]PlayerChoice{},
				false,
				make(chan struct{}),
				make(chan PlayerChoice),
				new(sync.Mutex),
				&WinnerDefiner{},
			},
			args: args{
				player: testPlayers["Adding last player"][4],
			},
			expectedPlayers:     testPlayers["Adding last player"],
			expectedRoomStarted: true,
			wantErr:             false,
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
			if !reflect.DeepEqual(testPlayers[tt.name], room.players) {
				t.Errorf("Room.AddPlayer() error adding players got = %v, expected = %v", room.players, testPlayers[tt.name])
			}
			if room.active != tt.expectedRoomStarted {
				t.Errorf("Room.AddPlayer() error %v: unexpected room state", tt.name)
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
