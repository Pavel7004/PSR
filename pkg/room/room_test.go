package room

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/subscribe"
	. "github.com/pavel/PSR/pkg/winner-definer"
)

func TestRoom_AddPlayer(t *testing.T) {
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
		},
	}
	initTestRoomFn := func(config RoomConfig, isRunning bool, players []*domain.Player) *Room {
		pub := subscribe.NewPublisher()
		room := NewRoom(config, pub)
		if isRunning {
			room.state = NewPlayingState(room)
		}
		for i := 0; i < len(players)-1; i++ {
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
			name:                "Test adding player in active game",
			roomConfig:          RoomConfig{2 * time.Second, 2, 5, false},
			isRunning:           true,
			players:             []*domain.Player{{ID: "Existing player 1"}, {ID: "Existing player 2"}},
			initFn:              initTestRoomFn,
			args:                args{&domain.Player{ID: "testPlayer"}},
			expectedRoomStarted: false,
			wantErr:             true,
		},
		{
			name:                "Adding player",
			roomConfig:          RoomConfig{2 * time.Second, 2, 5, false},
			isRunning:           false,
			players:             []*domain.Player{},
			initFn:              initTestRoomFn,
			args:                args{testPlayers["Adding player"][0]},
			expectedRoomStarted: false,
			wantErr:             false,
		},
		{
			name:                "Adding last player",
			roomConfig:          RoomConfig{2 * time.Second, 2, 5, false},
			isRunning:           false,
			players:             testPlayers["Adding last player"],
			initFn:              initTestRoomFn,
			args:                args{testPlayers["Adding last player"][1]},
			expectedRoomStarted: false,
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

func TestRoom_Choose(t *testing.T) {
	initTestRoomFn := func(config RoomConfig, isRunning bool, existingChoices []PlayerChoice) *Room {
		pub := subscribe.NewPublisher()
		room := NewRoom(config, pub)
		for i := 0; i < config.MaxPlayerCount; i++ {
			room.players = append(room.players, &domain.Player{
				ID: fmt.Sprintf("TestPlayer%d", i+1),
			})
		}
		if isRunning {
			room.state = NewPlayingState(room)
		}
		room.combinations = existingChoices
		return room
	}
	type args struct {
		choice *PlayerChoice
	}
	tests := []struct {
		name             string
		config           RoomConfig
		initFn           func(config RoomConfig, isRunning bool, existingChoices []PlayerChoice) *Room
		initCombinations []PlayerChoice
		isStarted        bool
		winners          []string
		args             args
		wantErr          bool
	}{
		{
			name:             "Try to choose in non-started room",
			config:           RoomConfig{2 * time.Second, 2, 5, false},
			initFn:           initTestRoomFn,
			initCombinations: []PlayerChoice{},
			isStarted:        false,
			winners:          nil,
			args:             args{&PlayerChoice{PlayerID: "TestPlayer1", Input: 0}},
			wantErr:          true,
		},
		{
			name:             "Choose in playing state",
			config:           RoomConfig{2 * time.Second, 2, 5, false},
			initFn:           initTestRoomFn,
			initCombinations: []PlayerChoice{},
			isStarted:        true,
			winners:          nil,
			args:             args{&PlayerChoice{PlayerID: "TestPlayer1", Input: 0}},
			wantErr:          false,
		},
		{
			name:             "Last player choose",
			config:           RoomConfig{2 * time.Second, 2, 5, false},
			initFn:           initTestRoomFn,
			initCombinations: []PlayerChoice{{PlayerID: "TestPlayer1", Input: 0}},
			isStarted:        true,
			winners:          []string{},
			args:             args{&PlayerChoice{PlayerID: "TestPlayer2", Input: 0}},
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := tt.initFn(tt.config, tt.isStarted, tt.initCombinations)
			sub := subscribe.NewSubscriber(1)
			room.observer.Subscribe(sub, "winners")
			if err := room.Choose(tt.args.choice); (err != nil) != tt.wantErr {
				t.Errorf("Room.Choose() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.winners != nil {
				if winners := sub.Receive(); reflect.DeepEqual(winners.([]string), tt.winners) {
					t.Errorf("Room.Choose() wrong winners got = %v, expected = %v", winners, tt.winners)
				}
			}
		})
	}
}
