package main

import (
	"log"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client represents a connected client
type Client struct {
	Name string
	Conn *websocket.Conn
}

// Send broadcasts a message to all clients
func (c *Client) Send(message string) {
	for _, client := range clients {
		client.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s: %s", c.Name, message)))
	}
}

var (
	clients   []*Client
	addClient = make(chan *Client)
	delClient = make(chan *Client)
	messages  = make(chan string)
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New client connected")
		conn, _ := upgrader.Upgrade(w, r, nil)
		// if err != nil {
		// 	fmt.Println("Error upgrading connection %v", err)
		// 	log.Println(err)
		// 	return 
		// }

		client := &Client{Conn: conn}
		addClient <- client
		fmt.Println("Client added")
		go func() {
			for {
				_, message, err := conn.ReadMessage()
				fmt.Println("Message received")
				fmt.Println(string(message))
				if err != nil {
					delClient <- client
					return
				}
				messages <- string(message)
			}
		}()
	})
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))

	// router.ListenAndServe(":8080", nil)
}

func handleClients() {
	for {
		select {
		case client := <-addClient:
			clients = append(clients, client)
			client.Name = fmt.Sprintf("User %d", len(clients))
			client.Send("Welcome to the chat!")
		case client := <-delClient:
			for i, c := range clients {
				if c == client {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			client.Send("Goodbye!")
		case message := <-messages:
			for _, client := range clients {
				client.Send(message)
			}
		}
	}
}
