package rooms

import (
	// "bufio"
	// "fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Message defines the structure of a chat message
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	RoomId   string `json:"roomid"`
}

var messages []Message

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		//append to messages
		messages = append(messages, msg)

		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	//get room id from

	// Register new client
	clients[ws] = true
	go func() {
		// Send last 10 messages to new client
		for _, msg := range messages {
			err := ws.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				ws.Close()
				delete(clients, ws)
			}
		}
	}()

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

// adapt handleConnections to handle multiple rooms
