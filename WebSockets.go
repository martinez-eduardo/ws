//package main
package WebSockets

import (
	//"fmt"

	"log"
	"net/http" 
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
type convert func(mensaje string)
var funcion = func(mensaje string) { log.Println(mensaje) }

func Reception(fn convert) {
	funcion = fn
}

func Read(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.Close()
		return
	}
	clients[ws] = true
	message(ws)
}


func message(client *websocket.Conn) {
	//send(client, "Bienvenido al WebSocket Golang")
	for {
		_, p, err := client.ReadMessage()
		if err != nil {
			client.Close()
			delete(clients, client)
			return
		}
		funcion(string(p))
		//onbroadcast(string(p))
	}
}

func Send(client *websocket.Conn, mensaje string) error {
	err := client.WriteMessage(websocket.TextMessage, []byte(mensaje))
	return err
}

func Broadcast(mensaje string) {
	for client := range clients {
		err := Send(client, mensaje)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
