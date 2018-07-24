package celerity

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// ChannelEvents are the type of event that occured.
	ChannelEvents = struct {
		Connect    ChannelEventType
		Disconnect ChannelEventType
		Message    ChannelEventType
	}{
		Connect:    "connect",
		Disconnect: "disconnect",
		Message:    "message",
	}

	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// ChannelEventType is an event type for websocket connections. Such as joining,
// leaving, message arrival, etc.
type ChannelEventType string

// SocketMessage can be used by channels to pass around a client client
// reference and a message
type SocketMessage struct {
	Client  *SocketClient
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
	upgrader   websocket.Upgrader
}

// NewChannel Undescribed
func NewChannel(h ChannelHandler) *Channel {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return &Channel{
		Clients:    map[*SocketClient]bool{},
		connect:    make(chan *SocketClient),
		disconnect: make(chan *SocketClient),
		message:    make(chan *SocketMessage),
		Handler:    h,
		upgrader:   upgrader,
	}
}

// Open run the channel and begin processing messages
func (ch *Channel) Open() {
	go func() {
		for {
			select {
			case client := <-ch.connect:
				ch.Clients[client] = true
				evt := ChannelEvent{
					Event: ChannelEvents.Connect,
					Data:  []byte{},
				}
				ch.Handler(client, evt)
			case client := <-ch.disconnect:
				evt := ChannelEvent{
					Event: ChannelEvents.Disconnect,
					Data:  []byte{},
				}
				ch.Handler(client, evt)
				if _, ok := ch.Clients[client]; ok {
					delete(ch.Clients, client)
					close(client.send)
				}

			case msg := <-ch.message:
				evt := ChannelEvent{
					Event: ChannelEvents.Message,
					Data:  msg.Message,
				}
				ch.Handler(msg.Client, evt)
			}
		}
	}()
}

// SocketClient Description
type SocketClient struct {
	Context Context
	ch      *Channel
	send    chan []byte
	conn    *websocket.Conn
}

// NewSocketClient creates a new client to control websocket connections
func NewSocketClient(c Context, ch *Channel, conn *websocket.Conn) *SocketClient {
	client := &SocketClient{
		ch:      ch,
		Context: c,
		conn:    conn,
		send:    make(chan []byte),
	}
	ch.connect <- client
	go client.readLoop()
	go client.writeLoop()
	return client
}

// Send sends data to the client
func (c *SocketClient) Send(data []byte) {
	c.send <- data
}

func (c *SocketClient) readLoop() {
	defer func() {
		c.ch.disconnect <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Context.Log.Errorf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if len(message) > 0 {
			c.ch.message <- &SocketMessage{c, message}
		}
	}
}

func (c *SocketClient) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ChannelRoute is a server route that serves a websocket connection
type ChannelRoute struct {
	Path    RoutePath
	Channel *Channel
}

// Match checks the incoming URI path against the route
func (r *ChannelRoute) Match(method string, path string) (bool, string) {
	ok, xtra := r.Path.Match(path)
	return (ok && xtra == ""), xtra
}

// Handle performs no logic for a websocket connection
func (r *ChannelRoute) Handle(c Context) Response {
	conn, err := r.Channel.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Log.Error(err)
	}
	NewSocketClient(c, r.Channel, conn)
	return Response{Handled: true}
}

// RoutePath returns the channels route path
func (r *ChannelRoute) RoutePath() RoutePath {
	return r.Path
}
