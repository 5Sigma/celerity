package server

import "golang.org/x/net/websocket"

var (
	RoomEvents = struct {
		Join  RoomEvent
		Leave RoomEvent
	}{
		Join:    "join",
		Leave:   "leave",
		Message: "message",
	}
)

type RoomEventType string

// RoomHandler is the handling function for incomming messages into the room
type RoomHandler func(*SocketClient, RoomEvent)

type RoomEvent struct {
	Event RoomEventType
	Data  []byte
}

// Room segments communication into a single context
type Room struct {
	clients map[*SocketClient]bool
	join    chan *SocketClient
	leave   chan *SocketClient
	message chan struct{SocketClient; []data}
	Context Context
	Handler RoomHandler
}

// SocketClient Description
type SocketClient struct {
	room *Room
	send chan []byte
	conn *websocket.Conn
}

// Open run the channel and begin processing messages
func (room *Room) Open() {
	for {
		select {
		case client := <-room.join:
			h.clients[client] = true
			evt := RoomEvent{
				Event: RoomEvents.Join,
				Data:  []byte{},
			}
			room.Handler(client, evt)
		case client := <-room.Leave:
			evt := RoomEvent{
				Event: RoomEvents.Leave,
				Data:  []byte{},
			}
			room.Handler(client, evt)
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.send)
			}

		case msg := <-room.message:
			evt := RoomEvent{
				Event: RoomEvent.Message,
				Data:  msg,
			}
			room.Handler(client, evt)
		}
	}
}
