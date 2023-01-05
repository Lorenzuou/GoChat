package rooms

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Client represents a chat client
type Client struct {

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// NewClient creates a new client
func createClient(conn *websocket.Conn) *Client {
	return &Client{

		conn: conn,
		send: make(chan []byte, 256),
	}
}

// Room represents a chat room
type Room struct {
	id string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type RoomManager struct {
	rooms map[string]*Room
}

func (rm *RoomManager) GetRoom(id string) *Room {
	return rm.rooms[id]
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}

}

// NewRoom creates a new room
func (rms *RoomManager) CreateRoom() *Room {
	id := len(rms.rooms) + 1

	r := &Room{
		id:         fmt.Sprintf("%d", id),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	rms.rooms[r.id] = r
	return r
}

// Run the room's event loop
func (ro *Room) Run() {
	for {
		select {
		case client := <-ro.register:
			ro.clients[client] = true
		case client := <-ro.unregister:
			if _, ok := ro.clients[client]; ok {
				delete(ro.clients, client)
				close(client.send)
			}
		case message := <-ro.broadcast:
			for client := range ro.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(ro.clients, client)
				}
			}
		}
	}
}
func (ro *Room) GetClients() map[*Client]bool {
	return ro.clients
}

func (ro *Room) updateBroadcast(message []byte) {
	ro.broadcast <- message
	return
}

func (ro *Room) GetId() string {
	return ro.id
}
