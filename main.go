package main

import (
	"bufio"
	"fmt"
	"net"
)

// Client represents a connected client
type Client struct {
	Name   string
	Reader *bufio.Reader
	Writer *bufio.Writer
}

// Send broadcasts a message to all clients
func (c *Client) Send(message string) {
	for _, client := range clients {
		client.Writer.WriteString(fmt.Sprintf("%s: %s\n", c.Name, message))
		client.Writer.Flush()
	}
}

var (
	clients   []*Client
	addClient = make(chan *Client)
	delClient = make(chan *Client)
	messages  = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server started, listening on port 8080")

	// Start a Goroutine to handle client connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			client := &Client{
				Reader: bufio.NewReader(conn),
				Writer: bufio.NewWriter(conn),
			}
			addClient <- client
			go func() {
				for {
					line, _, err := client.Reader.ReadLine()
					if err != nil {
						delClient <- client
						return
					}
					messages <- string(line)
				}
			}()
		}
	}()

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
