package server

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// ChannelEvents are the type of event that occured.
	ChannelEvents = struct {
		Connect    ChannelEvent
		Disconnect ChannelEvent
	}{
		Connect:    "connect",
		Disconnect: "disconnect",
		Message:    "message",
	}
)

// ChannelEventType is an event type for websocket connections. Such as joining,
// leaving, message arrival, etc.
type ChannelEventType string

// SocketMessage can be used by channels to pass around a client client
// reference and a message
type SocketMessage struct {
	Client  SocketClient
	Message []byte
}

// ChannelHandler is the handling function for incomming messages into the
// channel
type ChannelHandler func(*SocketClient, ChannelEvent)

// ChannelEvent is an event that can occure on the socket
type ChannelEvent struct {
	Event ChannelEventType
	Data  []byte
}

// Channel segments communication into a single context
type Channel struct {
	Clients    map[*SocketClient]bool
	connect    chan *SocketClient
	disconnect chan *SocketClient
	message    chan *SocketMessage
	Handler    ChannelHandler
}

// Open run the channel and begin processing messages
func (ch *Channel) Open() {
	for {
		select {
		case client := <-ch.connect:
			h.clients[client] = true
			evt := ChannelEvent{
				Event: ChannelEvents.Join,
				Data:  []byte{},
			}
			ch.Handler(client, evt)
		case client := <-ch.Leave:
			evt := ChannelEvent{
				Event: ChannelEvents.Leave,
				Data:  []byte{},
			}
			ch.Handler(client, evt)
			if _, ok := ch.clients[client]; ok {
				delete(ch.clients, client)
				close(client.send)
			}

		case msg := <-ch.message:
			evt := ChannelEvent{
				Event: ChannelEvent.Message,
				Data:  msg,
			}
			ch.Handler(client, evt)
		}
	}
}

// SocketClient Description
type SocketClient struct {
	Context Context
	ch      *Channel
	send    chan []byte
	conn    *websocket.Conn
}

// ReadLoop stats the the oscket
func (c *SocketClient) readLoop() {
	defer func() {
		c.ch.disconnect <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadLine(time.Now().Add(60 * second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, _ := c.conn.ReadMessage()
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.ch.message <- ClientMessage{c, message}
	}
}
