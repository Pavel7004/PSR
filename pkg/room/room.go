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
	active        bool
	stopCh        chan struct{}
	chooseCh      chan PlayerChoice
	stepMtx       *sync.Mutex
	winnerDefiner *WinnerDefiner
}

func NewRoom(config RoomConfig) *Room {
	return &Room{
		config:        config,
		players:       make([]*domain.Player, 0, config.MaxPlayerCount),
		combinations:  []PlayerChoice{},
		active:        false,
		stopCh:        make(chan struct{}),
		chooseCh:      make(chan PlayerChoice),
		stepMtx:       new(sync.Mutex),
		winnerDefiner: &WinnerDefiner{},
	}
}

func (room *Room) IsActive() bool {
	return room.active
}

func (room *Room) AddPlayer(player *domain.Player) error {
	if room.IsActive() {
		return ErrGameAlreadyStarted
	}
	room.stepMtx.Lock()
	room.players = append(room.players, player)
	log.Info().Msgf("Player %s added to the room", player.ID)

	if len(room.players) == room.config.MaxPlayerCount {
		room.active = true
		go room.Run()
	}
	room.stepMtx.Unlock()
	return nil
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
				break GAME_LOOP
			}
		}
	}
	winners := room.winnerDefiner.GetWinners(room.combinations)
	log.Info().Msgf("Winners: %v", winners)
}

func (room *Room) Choose(choice PlayerChoice) {
	room.chooseCh <- choice
}
