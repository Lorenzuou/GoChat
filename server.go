package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	rm "webchat.com/webchat/util"
)

// this function maps routes to functions
func mapRoutes() {
	mux["/"] = homePage
	mux["/rooms"] = roomsPage
	mux["/createRoom"] = createRoom
}

func mapRoutesWithParams() {
	muxP["/ws"] = wsEndpoint
	muxP["/room"] = roomPage

}

var rooms = rm.NewRoomManager()

func homePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("public/index.html")
	t.Execute(w, nil)

}

func roomsPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("public/rooms.gohtml")
	t.Execute(w, nil)

}

func createRoom(w http.ResponseWriter, r *http.Request) {

	room := rooms.CreateRoom()
	fmt.Println("room id: ", room.GetId())
	go room.Run()
	fmt.Println("room created")
	//send client to room with id as query param
	http.Redirect(w, r, "/room?id="+room.GetId(), http.StatusSeeOther)

}

func roomPage(w http.ResponseWriter, r *http.Request, params ...string) {

	// if rooms.GetRoom(id) == nil {
	// 	http.Redirect(w, r, "/rooms", http.StatusSeeOther)
	// }

	// t, _ := template.ParseFiles("public/room.gohtml")
	// t.Execute(w, id)

}

func wsEndpoint(w http.ResponseWriter, r *http.Request, params ...string) {
	fmt.Println("wsEndpoint")
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

// maps routes to functions
var mux map[string]func(http.ResponseWriter, *http.Request)
var muxP map[string]func(http.ResponseWriter, *http.Request, ...string)

func main() {
	fmt.Println("Web Chat usando Golang v0.01")

	server := http.Server{
		Addr:         ":8080",
		Handler:      &myHandler{},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	muxP = make(map[string]func(http.ResponseWriter, *http.Request, ...string))
	mapRoutes()
	mapRoutesWithParams()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func treatParameters(url string) string {
	params := strings.Split(url, "?")
	paramslist := strings.Split(params[1], "&")
	//map each parameter to a map
	fmt.Println(paramslist)
	os.Exit(1)
	return params[1]

}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	params := strings.Split(r.URL.String(), "?")
	if h, ok := mux[params[0]]; ok {
		h(w, r)
		return
	} else if h, ok := muxP[params[0]]; ok {
		h(w, r, treatParameters(params[1]))
		return
	}

	http.NotFound(w, r)

}
