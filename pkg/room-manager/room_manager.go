package roommanager

import (
	"errors"
	"sync"

	"github.com/pavel/PSR/pkg/room"
	"github.com/rs/zerolog/log"
)

var (
	ErrRoomAlreadyExists  = errors.New("Room with this name already exists.")
	ErrRoomDontExist      = errors.New("Room with this name aren't present.")
	ErrInvalidPlayerCount = errors.New("Room can't have less then 1 player.")
	ErrInvalidMaxScore    = errors.New("Room score must be more than 0.")
	ErrTimeoutTooShort    = errors.New("Room round time must be more than 5 seconds.")
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
	err := rm.CheckRoomConfig(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Bad room config.")
		return err
	}

	room, err := room.NewRoom(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Can't create new room.")
		return err
	}

	rm.mtx.Lock()
	rm.rooms[cfg.Name] = room
	rm.mtx.Unlock()

	go room.Main()

	return nil
}

func (rm *RoomManager) CheckRoomConfig(cfg *room.RoomConfig) error {
	rm.mtx.Lock()
	_, exist := rm.rooms[cfg.Name]
	rm.mtx.Unlock()
	if exist {
		return ErrRoomAlreadyExists
	}
	if cfg.MaxPlayerCount < 1 {
		return ErrInvalidPlayerCount
	}
	if cfg.MaxScore < 1 {
		return ErrInvalidMaxScore
	}
	if cfg.RoundTimeout.Seconds() < 5.0 {
		return ErrTimeoutTooShort
	}
	return nil
}

func (rm *RoomManager) GetRoomByID(name string) (*room.Room, error) {
	room, exist := rm.rooms[name]
	if !exist {
		return nil, ErrRoomDontExist
	}

	return room, nil
}
