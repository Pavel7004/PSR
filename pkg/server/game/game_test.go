package game

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/server/subscribe"
)

func TestGame_AddPlayer(t *testing.T) {
	type args struct {
		player *domain.Player
	}
	testPlayers := map[string][]*domain.Player{
		"Test adding player in active game": {
			domain.NewPlayer("Existing player 1"),
			domain.NewPlayer("Existing player 2"),
		},
		"Adding player": {
			domain.NewPlayer("Player 1"),
		},
		"Adding last player": {
			domain.NewPlayer("Player 1"),
			domain.NewPlayer("Player 2"),
		},
	}
	initTestGameFn := func(maxPlayers int, isRunning bool, players []*domain.Player) *Game {
		game := &Game{
			players:       players,
			combinations:  nil,
			state:         nil,
			observer:      subscribe.NewPublisher(),
			winnerDefiner: nil,
			scoremanager:  nil,
		}
		if isRunning {
			game.state = NewPlayingState(game)
		} else {
			game.state = NewWaitingState(game)
		}
		return game
	}
	tests := []struct {
		name                string
		roomCap             int
		isRunning           bool
		players             []*domain.Player
		initFn              func(maxPlayers int, isRunning bool, players []*domain.Player) *Game
		args                args
		expectedRoomStarted bool
		wantErr             bool
	}{
		{
			name:      "Test adding player in active game",
			roomCap:   2,
			isRunning: true,
			players: []*domain.Player{
				testPlayers["Test adding player in active game"][0],
				testPlayers["Test adding player in active game"][1],
			},
			initFn:              initTestGameFn,
			args:                args{domain.NewPlayer("testPlayer")},
			expectedRoomStarted: true,
			wantErr:             true,
		},
		{
			name:                "Adding player",
			roomCap:             2,
			isRunning:           false,
			players:             []*domain.Player{},
			initFn:              initTestGameFn,
			args:                args{testPlayers["Adding player"][0]},
			expectedRoomStarted: false,
			wantErr:             false,
		},
		{
			name:                "Adding last player",
			roomCap:             2,
			isRunning:           false,
			players:             []*domain.Player{testPlayers["Adding last player"][0]},
			initFn:              initTestGameFn,
			args:                args{testPlayers["Adding last player"][1]},
			expectedRoomStarted: true,
			wantErr:             false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := tt.initFn(tt.roomCap, tt.isRunning, tt.players)
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
	initTestRoomFn := func(maxPlayers int, isRunning bool, existingChoices map[string]domain.Choice) *Game {
		game := &Game{
			players:      make([]*domain.Player, 0, maxPlayers),
			combinations: existingChoices,
			observer:     subscribe.NewPublisher(),
		}
		if isRunning {
			game.state = NewPlayingState(game)
		} else {
			game.state = NewWaitingState(game)
		}

		for i := 0; i < maxPlayers; i++ {
			game.players = append(game.players, domain.NewPlayer(fmt.Sprintf("TestPlayer%d", i+1)))
		}

		return game
	}
	type args struct {
		id     string
		choice domain.Choice
	}
	tests := []struct {
		name             string
		roomCap          int
		initFn           func(maxPlayers int, isRunning bool, existingChoices map[string]domain.Choice) *Game
		initCombinations map[string]domain.Choice
		isStarted        bool
		winners          []string
		args             args
		wantErr          bool
	}{
		{
			name:             "Try to choose in non-started room",
			roomCap:          2,
			initFn:           initTestRoomFn,
			initCombinations: map[string]domain.Choice{},
			isStarted:        false,
			winners:          nil,
			args:             args{id: "TestPlayer1", choice: 0},
			wantErr:          true,
		},
		{
			name:             "Choose in playing state",
			roomCap:          2,
			initFn:           initTestRoomFn,
			initCombinations: map[string]domain.Choice{},
			isStarted:        true,
			winners:          nil,
			args:             args{id: "TestPlayer1", choice: 0},
			wantErr:          false,
		},
		{
			name:             "Choose with player not present",
			roomCap:          2,
			initFn:           initTestRoomFn,
			initCombinations: map[string]domain.Choice{},
			isStarted:        true,
			winners:          nil,
			args:             args{id: "Player1", choice: 0},
			wantErr:          true,
		},
		{
			name:             "Last player choose",
			roomCap:          2,
			initFn:           initTestRoomFn,
			initCombinations: map[string]domain.Choice{"TestPlayer1": 0},
			isStarted:        true,
			winners:          []string{},
			args:             args{id: "TestPlayer2", choice: 0},
			wantErr:          false,
		},
		{
			name:             "Last player choose, check winners",
			roomCap:          2,
			initFn:           initTestRoomFn,
			initCombinations: map[string]domain.Choice{"TestPlayer1": 0},
			isStarted:        true,
			winners:          []string{"TestPlayer2"},
			args:             args{id: "TestPlayer2", choice: 1},
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := tt.initFn(tt.roomCap, tt.isStarted, tt.initCombinations)
			sub := subscribe.NewSubscriber(1)
			room.observer.Subscribe(sub, "winners")
			if err := room.Choose(tt.args.id, tt.args.choice); (err != nil) != tt.wantErr {
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
	initTestRoomFn := func(maxPlayers int, players []*domain.Player) *Game {
		game := &Game{
			players:       players,
			combinations:  nil,
			state:         nil,
			observer:      nil,
			winnerDefiner: nil,
			scoremanager:  nil,
		}
		game.state = NewWaitingState(game)

		return game
	}
	type args struct {
		playerName string
	}
	tests := []struct {
		name    string
		roomCap int
		initFn  func(maxPlayers int, players []*domain.Player) *Game
		players []*domain.Player
		args    args
		want    bool
	}{
		{
			name:    "Has player",
			roomCap: 2,
			initFn:  initTestRoomFn,
			players: []*domain.Player{domain.NewPlayer("Player1")},
			args:    args{playerName: "Player1"},
			want:    true,
		},
		{
			name:    "player in not in the room",
			roomCap: 2,
			initFn:  initTestRoomFn,
			players: []*domain.Player{domain.NewPlayer("Player1")},
			args:    args{playerName: "Player2"},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := tt.initFn(tt.roomCap, tt.players)
			if got := room.HasPlayer(tt.args.playerName); got != tt.want {
				t.Errorf("Room.HasPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}
