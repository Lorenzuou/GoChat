package main

import (
	"fmt"
	"log"
	"net/http"
)

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

func homePage(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./public")).ServeHTTP(w, r)
	// fmt.Fprintf(w, "Home Page")
}

// func wsEndpoint(w http.ResponseWriter, r *http.Request) {
// 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println("Client Successfully Connected...")
// 	reader(ws)
// }

func main() {
	fmt.Println("Web Chat usando Golang v0.01")
	http.HandleFunc("/", homePage)

	// handle websocket connections reading the room id from url

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
