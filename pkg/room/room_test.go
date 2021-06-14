package room

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/winner-definer"
)

//! testify

func TestRoom_AddPlayer(t *testing.T) {
	type fields struct {
		config        RoomConfig
		players       []*domain.Player
		combinations  []PlayerChoice
		state         State
		stopCh        chan struct{}
		chooseCh      chan PlayerChoice
		stepMtx       *sync.Mutex
		winnerDefiner *WinnerDefiner
	}
	testRoom := Room{
		RoomConfig{
			5 * time.Minute,
			5,
			5,
			false,
		},
		make([]*domain.Player, 0, 5),
		[]PlayerChoice{},
		nil,
		nil,
		new(sync.Mutex),
		&WinnerDefiner{},
	}
	testRoom.state = NewPlayingState(&testRoom)
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


	initTestRoomFn := func(config RoomConfig, isRunning bool, players []*domain.Player) *Room {
				room := NewRoom(config)
				if isRunning {
					room.state = NewPlayingState(room)
				}
				for i := 0; i < len(players); i++ {
					room.AddPlayer(players[i])
				}
				return room
	}
	tests := []struct {
		name                string
		roomConfig          RoomConfig
		isRunning           bool
		players             []*domain.Player
		initFn              func(config RoomConfig, isRunning bool, players []*domain.Player) *Room
		args                args
		expectedRoomStarted bool
		wantErr             bool
	}{
		{
			name: "Test adding player in active game",
			roomConfig: RoomConfig{
				2 * time.Second,
				2,
				5,
				false,
			},
			isRunning: true,
			players: []*domain.Player{
				{
					ID: "Existing player 1",
				},
				{
					ID: "Existing player 2",
				},
			},
			initFn: initTestRoomFn,
			args: args{
				&domain.Player{
					ID: "testPlayer",
				},
			},
			wantErr:             true,
		},
		{
			name: "Adding player",
			roomConfig: RoomConfig{
				2 * time.Second,
				2,
				5,
				false,
			},
			isRunning: false,
			players: []*domain.Player{},
			initFn: initTestRoomFn,
			args: args{
				&domain.Player{
					ID: "new Player1",
				},
			},
			wantErr:             false,
		},
		{
			name: "Adding last player",
			roomConfig: RoomConfig{
				2 * time.Second,
				2,
				5,
				false,
			},
			isRunning: false,
			players: []*domain.Player{
				{
					ID: "Existing player 1",
				},
			},
			initFn: initTestRoomFn,
			args: args{
				&domain.Player{
					ID: "New Player1",
				},
			},
			wantErr:             false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := tt.initFn(tt.roomConfig, tt.isRunning, tt.players)
			if err := room.AddPlayer(tt.args.player); (err != nil) != tt.wantErr {
				t.Errorf("Room.AddPlayer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(testPlayers[tt.name], room.players) {
				t.Errorf("Room.AddPlayer() error adding players got = %v, expected = %v", room.players, testPlayers[tt.name])
			}
			if (reflect.TypeOf(room.state).String() == "PlayingState") && tt.expectedRoomStarted {
				t.Errorf("Room.AddPlayer() error %v: unexpected room state", tt.name)
			}
		})
	}
}
