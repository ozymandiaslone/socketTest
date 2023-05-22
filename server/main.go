package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // uses default options

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request, converting it into a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Close the connection when the function returns
	defer conn.Close()
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Message received: %s", msg)
		if err := conn.WriteMessage(messageType, msg); err != nil {
			log.Println(err)
			break
		}
	}
}
func main() {
	http.HandleFunc("/socket", socketHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
