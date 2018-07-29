package celerity

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSockets(t *testing.T) {
	server := New()

	server.Router.Root.Channel("welcome", "/welcome",
		func(client *SocketClient, e ChannelEvent) {
			if e.Event == ChannelEvents.Connect {
				client.SendString("hello client")
			}
		})

	server.Channel("echo", "/echo", func(client *SocketClient, e ChannelEvent) {
		client.SendRaw(e.Data)
	})

	server.Channel("broadcast", "/broadcast",
		func(client *SocketClient, e ChannelEvent) {
			client.Channel().Broadcast(map[string]string{"test": "test"})
		})

	server.Channel("raw", "/broadcast-raw",
		func(client *SocketClient, e ChannelEvent) {
			client.Channel().BroadcastRaw([]byte("raw broadcast"))
		})

	server.Channel("empty", "/empty",
		func(client *SocketClient, e ChannelEvent) {})

	server.GET("/send", func(c Context) Response {
		c.Server.Channels["empty"].BroadcastRaw([]byte("test through another request"))
		return c.R("ok")
	})

	ts := httptest.NewServer(server)
	defer ts.Close()
	tsURL, err := url.Parse(ts.URL)
	tsURL.Scheme = "ws"
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Run("welcome", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(tsURL.String()+"/welcome", nil)
		if err != nil {
			t.Fatal("dial:", err)
		}
		defer c.Close()

		_, message, err := c.ReadMessage()
		if err != nil {
			t.Errorf("read: %s", err)
			return
		}
		if expected := "hello client"; string(message) != expected {
			t.Errorf("Recieved: '%s' wanted '%s'", string(message), expected)
		}
	})

	t.Run("echo", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(tsURL.String()+"/echo", nil)
		if err != nil {
			t.Fatal("dial:", err)
		}
		defer c.Close()

		err = c.WriteMessage(websocket.TextMessage, []byte("this is a test"))
		if err != nil {
			t.Fatal(err.Error())
		}
		_, message, err := c.ReadMessage()
		_, message, err = c.ReadMessage()
		if err != nil {
			t.Fatalf("read: %s", err)
			return
		}

		if expected := "this is a test"; string(message) != expected {
			t.Errorf("Recieved: '%s' wanted '%s'", string(message), expected)
		}
	})

	t.Run("broadcast", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(tsURL.String()+"/broadcast", nil)
		if err != nil {
			t.Fatal("dial:", err)
		}
		defer c.Close()

		_, message, err := c.ReadMessage()
		if err != nil {
			t.Errorf("read: %s", err)
			return
		}
		expected := `{"test":"test"}`
		if string(message) != expected {
			t.Errorf("Recieved: '%s' wanted '%s'", string(message), expected)
		}
	})

	t.Run("broadcast-raw", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(tsURL.String()+"/broadcast-raw", nil)
		if err != nil {
			t.Fatal("dial:", err)
		}
		defer c.Close()

		_, message, err := c.ReadMessage()
		if err != nil {
			t.Errorf("read: %s", err)
			return
		}
		expected := `raw broadcast`
		if string(message) != expected {
			t.Errorf("Recieved: '%s' wanted '%s'", string(message), expected)
		}
	})

	t.Run("through endpoint", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(tsURL.String()+"/empty", nil)
		if err != nil {
			t.Fatal("dial:", err)
		}
		defer c.Close()

		http.Get(ts.URL + "/send")

		_, message, err := c.ReadMessage()
		if err != nil {
			t.Errorf("read: %s", err)
			return
		}
		expected := `test through another request`
		if string(message) != expected {
			t.Errorf("Recieved: '%s' wanted '%s'", string(message), expected)
		}
	})
}

func TestRoomRemove(t *testing.T) {
	c1 := &SocketClient{ID: 1}
	c2 := &SocketClient{ID: 2}
	c3 := &SocketClient{ID: 3}
	room := NewChannelRoom("test-room", c1, c2, c3)
	c1.Rooms = []*ChannelRoom{room}
	c2.Rooms = []*ChannelRoom{room}
	c3.Rooms = []*ChannelRoom{room}

	room.Remove(c1)
	if l := len(room.Clients); l != 2 {
		t.Errorf("room should have 2 clients, has %d", l)
	}
	if l := len(c1.Rooms); l != 0 {
		t.Errorf("client1 should have 0 rooms, has %d", l)
	}
}

func TestRoomAdd(t *testing.T) {
	c1 := &SocketClient{ID: 1}
	room := NewChannelRoom("test-room")
	room.Add(c1)
	if len(room.Clients) != 1 {
		t.Errorf("room should have 1 client, has %d", len(room.Clients))
	}
}

func emptyChannelHandler(client *SocketClient, e ChannelEvent) {}

func TestCreateRoomFromChannel(t *testing.T) {
	ch := NewChannel(emptyChannelHandler)
	ch.Room("test")
	if l := len(ch.Rooms); l != 1 {
		t.Errorf("channel should have 1 room, has %d", l)
	}
	ch.Room("test")
	if l := len(ch.Rooms); l != 1 {
		t.Errorf("channel should have 1 room, has %d", l)
	}
}

func TestChannelEvents(t *testing.T) {
	evt := ChannelEvent{
		Event: ChannelEvents.Message,
		Data:  []byte(`{"user": { "name": "Alice" }}`),
	}
	t.Run("get", func(t *testing.T) {
		name := evt.Get("user.name").String()
		if name != "Alice" {
			t.Errorf("value at user.name should be 'Alice', got %s", name)
		}
	})
	t.Run("extract", func(t *testing.T) {
		v := struct {
			User struct {
				Name string `json:"name"`
			} `json:"user"`
		}{}
		evt.Extract(&v)
		if v.User.Name != "Alice" {
			t.Errorf("value at user.name should be 'Alice', got %s", v.User.Name)
		}
	})
	t.Run("extractat", func(t *testing.T) {
		v := struct {
			Name string `json:"name"`
		}{}
		evt.ExtractAt("user", &v)
		if v.Name != "Alice" {
			t.Errorf("value at user.name should be 'Alice', got %s", v.Name)
		}
	})
}
