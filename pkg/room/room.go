package room

import (
	"errors"
	"sync"

	"github.com/pavel/PSR/pkg/domain"
	. "github.com/pavel/PSR/pkg/winner-definer"
	"github.com/rs/zerolog/log"
)

var (
	ErrGameAlreadyStarted = errors.New("The game is already started!")
)

type Room struct {
	config        RoomConfig
	players       []*domain.Player
	combinations  []PlayerChoice
	state         State
	stopCh        chan struct{}
	chooseCh      chan PlayerChoice
	stepMtx       *sync.Mutex
	winnerDefiner *WinnerDefiner
}

func NewRoom(config RoomConfig) *Room {
	room := Room{
		config:        config,
		players:       make([]*domain.Player, 0, config.MaxPlayerCount),
		combinations:  []PlayerChoice{},
		state:         nil,
		stopCh:        make(chan struct{}),
		chooseCh:      make(chan PlayerChoice),
		stepMtx:       new(sync.Mutex),
		winnerDefiner: &WinnerDefiner{},
	}
	room.state = NewWaitingState(&room)
	return &room
}

func (room *Room) AddPlayer(player *domain.Player) error {
	return room.state.AddPlayer(player)
}

func (room *Room) Run() {
	log.Info().Msg("Room started")
GAME_LOOP:
	for {
		select {
		case <-room.stopCh:
			log.Info().Msg("Room stopped")
			break GAME_LOOP
		case choice := <-room.chooseCh:
			room.combinations = append(room.combinations, choice)
			if len(room.combinations) == len(room.players) {
				winners := room.winnerDefiner.GetWinners(room.combinations)
				log.Info().Msgf("Winners: %v", winners)
				room.combinations = make([]PlayerChoice, 0, room.config.MaxPlayerCount)
				// break GAME_LOOP
			}
		}
	}
}

func (room *Room) Choose(choice PlayerChoice) {
	room.chooseCh <- choice
}
