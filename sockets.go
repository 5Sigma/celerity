package celerity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sony/sonyflake"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
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

	newline     = []byte{'\n'}
	space       = []byte{' '}
	idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})
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

// Get returns a value at a JSON path in the data.
func (ce *ChannelEvent) Get(path string) gjson.Result {
	return gjson.GetBytes(ce.Data, path)
}

// Extract unmarshals the JSON into a struct
func (ce *ChannelEvent) Extract(obj interface{}) error {
	return json.Unmarshal(ce.Data, &obj)
}

// ExtractAt Unmarshals JSON at a path into a struct.
func (ce *ChannelEvent) ExtractAt(path string, obj interface{}) error {
	raw := gjson.GetBytes(ce.Data, path).Raw
	return json.Unmarshal([]byte(raw), &obj)
}

// Channel segments communication into a single context
type Channel struct {
	Clients    map[*SocketClient]bool
	connect    chan *SocketClient
	disconnect chan *SocketClient
	message    chan *SocketMessage
	Handler    ChannelHandler
	upgrader   websocket.Upgrader
	Rooms      map[string]*ChannelRoom
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
		Rooms:      map[string]*ChannelRoom{},
	}
}

// Room will return a room for the given name or it will initialize a new room
// for that name
func (ch *Channel) Room(name string) *ChannelRoom {
	if val, ok := ch.Rooms[name]; ok {
		return val
	}
	r := NewChannelRoom(name)
	ch.Rooms[name] = r
	return r
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
				for _, r := range client.Rooms {
					r.Remove(client)
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

// Broadcast send an event to all clients in the channel
func (ch *Channel) Broadcast(msg interface{}) {
	for c, connected := range ch.Clients {
		if connected {
			c.Send(msg)
		}
	}
}

// BroadcastRaw send an event to all clients in the channel
func (ch *Channel) BroadcastRaw(msg []byte) {
	for c, connected := range ch.Clients {
		if connected {
			c.SendRaw(msg)
		}
	}
}

// SocketClient Description
type SocketClient struct {
	Context Context
	ID      uint64
	ch      *Channel
	send    chan []byte
	conn    *websocket.Conn
	Rooms   []*ChannelRoom
}

// NewSocketClient creates a new client to control websocket connections
func NewSocketClient(c Context, ch *Channel, conn *websocket.Conn) *SocketClient {
	id, _ := idGenerator.NextID()
	client := &SocketClient{
		ch:      ch,
		ID:      id,
		Context: c,
		conn:    conn,
		send:    make(chan []byte),
		Rooms:   []*ChannelRoom{},
	}
	ch.connect <- client
	go client.readLoop()
	go client.writeLoop()
	return client
}

// Channel returns the channel the socket belogns to
func (c *SocketClient) Channel() *Channel {
	return c.ch
}

// Send sends data to the client
func (c *SocketClient) Send(msg interface{}) {
	b, err := json.Marshal(msg)
	if err != nil {
		c.Context.Log.Error(err.Error())
		return
	}
	c.SendRaw(b)
}

// SendString sends a string to the client
func (c *SocketClient) SendString(msg string) {
	c.SendRaw([]byte(msg))
}

// SendRaw sends bytes to the client
func (c *SocketClient) SendRaw(msg []byte) {
	c.send <- append(msg)
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

// ChannelRoom is a group of socket clients and can be used to partition the
// clients into arbitrary groups to make it easy to send messages back and
// fourth.
type ChannelRoom struct {
	Name    string
	Clients []*SocketClient
	mx      *sync.Mutex
}

// NewChannelRoom Undescribed
func NewChannelRoom(name string, clients ...*SocketClient) *ChannelRoom {
	return &ChannelRoom{
		Name:    name,
		Clients: clients,
		mx:      &sync.Mutex{},
	}
}

// Add adds a client to the room
func (r *ChannelRoom) Add(c *SocketClient) {
	r.mx.Lock()
	r.Clients = append(r.Clients, c)
	c.Rooms = append(c.Rooms, r)
	r.mx.Unlock()
}

// Remove adds a client to the room
func (r *ChannelRoom) Remove(c *SocketClient) {
	r.mx.Lock()
	for i := range r.Clients {
		if c == r.Clients[i] {
			r.Clients[i] = r.Clients[len(r.Clients)-1]
			r.Clients = r.Clients[:len(r.Clients)-1]
			break
		}
	}
	for i := range c.Rooms {
		if r == c.Rooms[i] {
			c.Rooms[i] = c.Rooms[len(c.Rooms)-1]
			c.Rooms = c.Rooms[:len(c.Rooms)-1]
			break
		}
	}
	r.mx.Unlock()
}

// Broadcast send an event to all clients in the channel
func (r *ChannelRoom) Broadcast(msg interface{}) {
	for _, c := range r.Clients {
		c.Send(msg)
	}
}

// BroadcastRaw send an event to all clients in the channel
func (r *ChannelRoom) BroadcastRaw(msg []byte) {
	for _, c := range r.Clients {
		c.SendRaw(msg)
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
