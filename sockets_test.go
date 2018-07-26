package celerity

import (
	"log"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSockets(t *testing.T) {
	server := New()

	server.Router.Root.Channel("/welcome", func(client *SocketClient, e ChannelEvent) {
		if e.Event == ChannelEvents.Connect {
			client.SendRaw([]byte("hello client"))
		}
	})

	server.Router.Root.Channel("/echo", func(client *SocketClient, e ChannelEvent) {
		println(string(e.Data))
		client.SendRaw(e.Data)
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
			log.Fatal("dial:", err)
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
			log.Fatal("dial:", err)
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
