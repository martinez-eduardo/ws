//package main
package WebSockets

import (
	//"fmt"

	"log"
	"net/http" 
	"github.com/gorilla/websocket"
)

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
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			ws.Close()
			delete(clients, ws)
			return
		}
		client=ws  // <- hay que corregir agregandole programacion orientada a objetos para tener multiples instancias
		funcion(string(p))
	}
}


/*TODO ESTE CONTENIDO DEBERIA ESTAR EN OTRA INSTANCIA*/

var clients = make(map[*websocket.Conn]bool)
var client *websocket.Conn

func Send(mensaje string) error {
	err := client.WriteMessage(websocket.TextMessage, []byte(mensaje))
	return err
}

func Broadcast(mensaje string) {
	for ws := range clients {
		err := ws.WriteMessage(websocket.TextMessage, []byte(mensaje))
		if err != nil {
			ws.Close()
			delete(clients, ws)
		}
	}
}
