package game

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/subscribe"
	. "github.com/pavel/PSR/pkg/winner-definer"
)

func TestGame_AddPlayer(t *testing.T) {
	type args struct {
		player *domain.Player
	}
	testPlayers := map[string][]*domain.Player{
		"Test adding player in active game": {},
		"Adding player": {
			domain.NewPlayer("Player 1"),
		},
		"Adding last player": {
			domain.NewPlayer("Player 1"),
			domain.NewPlayer("Player 2"),
		},
	}
	initTestGameFn := func(config *GameConfig, isRunning bool, players []*domain.Player) *Game {
		pub := subscribe.NewPublisher()
		room := NewGame(config, pub)
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
		roomConfig          *GameConfig
		isRunning           bool
		players             []*domain.Player
		initFn              func(config *GameConfig, isRunning bool, players []*domain.Player) *Game
		args                args
		expectedRoomStarted bool
		wantErr             bool
	}{
		{
			name:                "Test adding player in active game",
			roomConfig:          &GameConfig{"", 2 * time.Second, 2, 5, false},
			isRunning:           true,
			players:             []*domain.Player{domain.NewPlayer("Existing player 1"), domain.NewPlayer("Existing player 2")},
			initFn:              initTestGameFn,
			args:                args{domain.NewPlayer("testPlayer")},
			expectedRoomStarted: false,
			wantErr:             true,
		},
		{
			name:                "Adding player",
			roomConfig:          &GameConfig{"", 2 * time.Second, 2, 5, false},
			isRunning:           false,
			players:             []*domain.Player{},
			initFn:              initTestGameFn,
			args:                args{testPlayers["Adding player"][0]},
			expectedRoomStarted: false,
			wantErr:             false,
		},
		{
			name:                "Adding last player",
			roomConfig:          &GameConfig{"", 2 * time.Second, 2, 5, false},
			isRunning:           false,
			players:             testPlayers["Adding last player"],
			initFn:              initTestGameFn,
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

func TestGame_Choose(t *testing.T) {
	initTestRoomFn := func(config *GameConfig, isRunning bool, existingChoices []PlayerChoice) *Game {
		pub := subscribe.NewPublisher()
		room := NewGame(config, pub)
		for i := 0; i < config.MaxPlayerCount; i++ {
			room.players = append(room.players, domain.NewPlayer(fmt.Sprintf("TestPlayer%d", i+1)))
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
		config           *GameConfig
		initFn           func(config *GameConfig, isRunning bool, existingChoices []PlayerChoice) *Game
		initCombinations []PlayerChoice
		isStarted        bool
		winners          []string
		args             args
		wantErr          bool
	}{
		{
			name:             "Try to choose in non-started room",
			config:           &GameConfig{"", 2 * time.Second, 2, 5, false},
			initFn:           initTestRoomFn,
			initCombinations: []PlayerChoice{},
			isStarted:        false,
			winners:          nil,
			args:             args{&PlayerChoice{PlayerID: "TestPlayer1", Input: 0}},
			wantErr:          true,
		},
		{
			name:             "Choose in playing state",
			config:           &GameConfig{"", 2 * time.Second, 2, 5, false},
			initFn:           initTestRoomFn,
			initCombinations: []PlayerChoice{},
			isStarted:        true,
			winners:          nil,
			args:             args{&PlayerChoice{PlayerID: "TestPlayer1", Input: 0}},
			wantErr:          false,
		},
		{
			name:             "Choose with player not present",
			config:           &GameConfig{"", 2 * time.Second, 2, 5, false},
			initFn:           initTestRoomFn,
			initCombinations: []PlayerChoice{},
			isStarted:        true,
			winners:          nil,
			args:             args{&PlayerChoice{PlayerID: "Player1", Input: 0}},
			wantErr:          true,
		},
		{
			name:             "Last player choose",
			config:           &GameConfig{"", 2 * time.Second, 2, 5, false},
			initFn:           initTestRoomFn,
			initCombinations: []PlayerChoice{{PlayerID: "TestPlayer1", Input: 0}},
			isStarted:        true,
			winners:          []string{},
			args:             args{&PlayerChoice{PlayerID: "TestPlayer2", Input: 0}},
			wantErr:          false,
		},
		{
			name:             "Last player choose, check winners",
			config:           &GameConfig{"", 2 * time.Second, 2, 5, false},
			initFn:           initTestRoomFn,
			initCombinations: []PlayerChoice{{PlayerID: "TestPlayer1", Input: 0}},
			isStarted:        true,
			winners:          []string{"TestPlayer2"},
			args:             args{&PlayerChoice{PlayerID: "TestPlayer2", Input: 1}},
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
				msg, ok := sub.Receive().([]string)
				if !ok {
					t.Errorf("Room.Choose() wrong winners message type: got = %T, expected = %T", msg, tt.winners)
				}
				if !reflect.DeepEqual(msg, tt.winners) {
					t.Errorf("Room.Choose() wrong winners got = %v, expected = %v", msg, tt.winners)
				}
			}
		})
	}
}

func TestGame_HasPlayer(t *testing.T) {
	initTestRoomFn := func(config *GameConfig, players []*domain.Player) *Game {
		pub := subscribe.NewPublisher()
		room := NewGame(config, pub)
		room.players = players
		return room
	}
	type args struct {
		playerName string
	}
	tests := []struct {
		name    string
		config  *GameConfig
		initFn  func(config *GameConfig, players []*domain.Player) *Game
		players []*domain.Player
		args    args
		want    bool
	}{
		{
			name:    "Has player",
			config:  &GameConfig{"", 2 * time.Second, 2, 5, false},
			initFn:  initTestRoomFn,
			players: []*domain.Player{domain.NewPlayer("Player1")},
			args:    args{playerName: "Player1"},
			want:    true,
		},
		{
			name:    "player in not in the room",
			config:  &GameConfig{"", 2 * time.Second, 2, 5, false},
			initFn:  initTestRoomFn,
			players: []*domain.Player{domain.NewPlayer("Player1")},
			args:    args{playerName: "Player2"},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := tt.initFn(tt.config, tt.players)
			if got := room.HasPlayer(tt.args.playerName); got != tt.want {
				t.Errorf("Room.HasPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}
