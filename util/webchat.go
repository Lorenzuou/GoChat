package rooms

import (
	// "bufio"
	// "fmt"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Message defines the structure of a chat message
type Message struct {
	Username string `json:"Username"`
	Message  string `json:"Message"`
	RoomId   string `json:"RoomId"`
}

var messages []Message

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendMessage(message *Message, conn *websocket.Conn) {
	if err := conn.WriteJSON(&message); err != nil {
		log.Printf("error: %v", err)
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request, roomId string) {

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	client := CreateClient(ws)

	room := GetRoomById(roomId)
	//make sure client is not in room
	defer client.Leave(room)

	//if room is not running, start it
	if !room.GetStatus() {
		room.SetStatus(true)
		go room.Run()
	}

	room.updateRegister(client)

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		//convert msg to byte
		// Send the newly received message to the broadcast channel
		fmt.Println("msg: ", msg)
		room.updateBroadcast(&msg)
	}
}

// adapt handleConnections to handle multiple rooms
