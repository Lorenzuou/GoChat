//test server for testing a websocket on golang application

package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}


func main() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/ws", wsEndpoint)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
    }
    log.Println("Client Successfully Connected...")
    reader(ws)
}

func reader(conn *websocket.Conn) {
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        log.Println(string(p))
        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
    }
}

func writer(conn *websocket.Conn) {
    for {
        time.Sleep(1 * time.Second)
        err := conn.WriteMessage(1, []byte("Hi Client"))
        if err != nil {
            log.Println(err)
            return
        }
    }
}

