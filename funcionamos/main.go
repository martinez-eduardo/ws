package main
 
import (
	"fmt"
	"log"
	"net/http"
 
	"github.com/gorilla/websocket"
)
 
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
 
 
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))
 
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
 
	}
}
 
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!DOCTYPE HTML>

<html>
   <head>
      
      <script type = "text/javascript">
         function WebSocketTest() {
            
            if ("WebSocket" in window) {
               alert("WebSocket is supported by your Browser!");
               
               // Let us open a web socket
               var ws = new WebSocket("ws://localhost:8080/ws");
				
               ws.onopen = function() {
                  
                  // Web Socket is connected, send data using send()
                  ws.send("HOLA DESDE FIREFOX");
                  alert("Mensaje enviado");
               };
				
               ws.onmessage = function (evt) { 
                  var received_msg = evt.data;
                  alert(received_msg);
               };
				
               ws.onclose = function() { 
                  
                  // websocket is closed.
                  alert("Connection is closed..."); 
               };
            } else {
              
               // The browser doesn't support WebSocket
               alert("WebSocket NOT supported by your Browser!");
            }
         }
      </script>
		
   </head>
   
   <body>
      <div id = "sse">
         <a href = "javascript:WebSocketTest()">Run WebSocket</a>
      </div>
      
   </body>
</html>`)
}
 
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
 
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("HOLA DESDE GOLANG"))
	if err != nil {
		log.Println(err)
	}
 
	reader(ws)
}
 
func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}
 
func main() {
	fmt.Println("Hello World")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}