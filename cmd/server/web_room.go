package main

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/pavel/PSR/pkg/domain"
	"github.com/pavel/PSR/pkg/room"
	"github.com/pavel/PSR/pkg/subscribe"
)

type WebRoom struct {
	name         string
	room         *room.Room
	connections  map[string]*websocket.Conn
	roomStateSub *subscribe.Subscriber
	winnersSub   *subscribe.Subscriber
}

func NewWebRoom(name string, maxPlayers int) *WebRoom {
	p := subscribe.NewPublisher()
	rm := room.NewRoom(
		room.RoomConfig{
			StepTimeout:    5 * time.Second,
			MaxPlayerCount: maxPlayers,
			MaxScore:       5,
			OnlyComputer:   false,
		},
		p,
	)
	subRoomState := subscribe.NewSubscriber(0)
	p.Subscribe(subRoomState, "room_started")
	subWinners := subscribe.NewSubscriber(0)
	p.Subscribe(subWinners, "winners")
	return &WebRoom{
		name:         name,
		room:         rm,
		connections:  make(map[string]*websocket.Conn, maxPlayers),
		roomStateSub: subRoomState,
		winnersSub:   subWinners,
	}
}

func (r *WebRoom) StartGame() {
	r.roomStateSub.Receive()
	startMsg := []byte("Игра началась")
	for _, conn := range r.connections {
		conn.WriteMessage(websocket.TextMessage, startMsg)
	}
	time.Sleep(2 * time.Second)
	r.CloseConnections()
}

func (r *WebRoom) CloseConnections() {
	for _, conn := range r.connections {
		conn.Close()
	}
}

func (r *WebRoom) AddPlayer(id string) {
	r.room.AddPlayer(domain.NewPlayer(id))
}
