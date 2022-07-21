package roommanager

import (
	"errors"
	"sync"

	"github.com/pavel/PSR/pkg/room"
	"github.com/rs/zerolog/log"
)

var (
	ErrRoomAlreadyExists = errors.New("Room with this name already exists.")
)

type RoomManager struct {
	rooms map[string]*room.Room
	mtx   *sync.Mutex
}

func New() *RoomManager {
	return &RoomManager{
		rooms: map[string]*room.Room{},
		mtx:   new(sync.Mutex),
	}
}

func (rm *RoomManager) CreateRoom(cfg *room.RoomConfig) error {
	rm.mtx.Lock()
	if _, exist := rm.rooms[cfg.Name]; exist {
		rm.mtx.Unlock()
		log.Warn().Msgf("RoomManager: room with name %s already exists", cfg.Name)

		return ErrRoomAlreadyExists
	}
	rm.rooms[cfg.Name] = room.NewRoom(cfg)
	rm.mtx.Unlock()

	return nil
}
