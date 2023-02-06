package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	rm "webchat.com/webchat/util"
)

// this function maps routes to functions
func mapRoutes() {
	mux["/"] = homePage
	mux["/register"] = register
	mux["/rooms"] = roomsPage
	mux["/createRoom"] = createRoom
}

func mapRoutesWithParams() {
	muxP["/ws"] = wsEndpoint
	muxP["/room"] = roomPage

}

var rooms = rm.GetRooms()

var users = map[string]string{
	"1": "1",
	"user2": "password2",
}

var isAuthenticated = false

func homePage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/login.gohtml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if p, ok := users[username]; ok && p == password {
			// Autenticação bem-sucedida, redirecionar para a página protegida
			isAuthenticated = true
			http.Redirect(w, r, "/rooms", http.StatusFound)
			return
		}

		// Autenticação falhou, exibir uma mensagem de erro
		fmt.Fprintln(w, "Usuário ou senha inválidos")
		return
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/register.gohtml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")

		users[username] = password

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func roomsPage(w http.ResponseWriter, r *http.Request) {

	if !isAuthenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	roomsIds := rm.GetRomsIds()

	roomsInterfaces := make([]interface{}, len(roomsIds))
	for i, id := range roomsIds {
		roomsInterfaces[i] = rm.GetRoomById(id).GetRoomInterface()
	}

	t, _ := template.ParseFiles("templates/rooms.gohtml")
	t.Execute(w, roomsInterfaces)

}

func createRoom(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// get post data
	roomName := r.FormValue("name")
	roomDescription := r.FormValue("description")

	room := rm.CreateRoom(roomName, roomDescription)
	fmt.Println("room id: ", room.GetId())
	//send client to room with id as query param
	http.Redirect(w, r, "/room?id="+room.GetId(), http.StatusSeeOther)

}

func roomPage(w http.ResponseWriter, r *http.Request, params map[string]string) {

	if !isAuthenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	id := params["id"]
	if rm.GetRoomById(id) == nil {
		http.Redirect(w, r, "/rooms", http.StatusSeeOther)
	} else {

		roomMap := rm.GetRoomById(id).GetRoomInterface()

		t, _ := template.ParseFiles("templates/room.gohtml")
		t.Execute(w, roomMap)
	}

}

func wsEndpoint(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := params["id"]
	rm.HandleConnections(w, r, id)
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
var muxP map[string]func(http.ResponseWriter, *http.Request, map[string]string)

func main() {
	fmt.Println("Web Chat usando Golang v0.01")

	server := http.Server{
		Addr:         ":8080",
		Handler:      &myHandler{},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	muxP = make(map[string]func(http.ResponseWriter, *http.Request, map[string]string))
	mapRoutes()
	mapRoutesWithParams()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func treatParameters(url string) map[string]string {
	//get all the parameters and return it as a map of string

	paramsList := strings.Split(url, "&")
	params := make(map[string]string)
	for _, param := range paramsList {
		paramList := strings.Split(param, "=")
		params[paramList[0]] = paramList[1]
	}
	return params
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
