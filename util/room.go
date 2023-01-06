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
	send chan *Message
}

// NewClient creates a new client
func CreateClient(conn *websocket.Conn) *Client {
	return &Client{

		conn: conn,
		send: make(chan *Message),
	}
}

func (c *Client) Leave(room *Room) {
	room.leave <- c
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

// Room represents a chat room
type Room struct {
	id string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	join chan *Client

	// leave requests from clients.
	leave chan *Client

	//Check if room is running
	isrunning bool
}

// List of rooms
var rmms = map[string]*Room{}

func GetRooms() map[string]*Room {
	return rmms
}

func GetRomsIds() []string {
	var ids []string
	for id := range rmms {
		ids = append(ids, id)
	}
	return ids
}

func GetRoomById(id string) *Room {
	return rmms[id]
}

// NewRoom creates a new room
func CreateRoom() *Room {
	id := len(rmms) + 1

	r := &Room{
		id:        fmt.Sprintf("%d", id),
		broadcast: make(chan *Message),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		clients:   make(map[*Client]bool),
		isrunning: false,
	}
	rmms[r.id] = r
	return r
}

// Run the room's event loop
func (ro *Room) Run() {
	for {
		select {
		case client := <-ro.join:
			ro.clients[client] = true
		case client := <-ro.leave:
			if _, ok := ro.clients[client]; ok {
				delete(ro.clients, client)
				close(client.send)
			}
		case message := <-ro.broadcast:
			for client := range ro.clients {

				sendMessage(message, client.GetConn())
			}
		}
	}
}

func (ro *Room) GetClients() map[*Client]bool {
	return ro.clients
}

func (ro *Room) updateBroadcast(message *Message) {
	ro.broadcast <- message
	return
}

func (ro *Room) SetStatus(status bool) {
	ro.isrunning = status
}

func (ro *Room) GetStatus() bool {
	return ro.isrunning
}

func (ro *Room) updateRegister(client *Client) {
	ro.join <- client
	return
}

func (ro *Room) GetId() string {
	return ro.id
}
