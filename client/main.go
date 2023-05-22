package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Received: %s", message)
	}
}

func main() {
	done = make(chan interface{})
	interrupt = make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	socketUrl := "ws://localhost:8080/socket"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	for {
		select {
		case <-time.After(10 * time.Millisecond):
			if err := conn.WriteMessage(websocket.TextMessage, []byte("THIS IS A TEST SOCKET BITCH")); err != nil {
				log.Println("Error writing message to socket:", err)
			}
		case <-interrupt:
			log.Println("Interrupt detected, terminating.")
			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Println("Error closing socket:", err)
				return
			}
			select {
			case <-done:
				log.Println("Done")
			case <-time.After(1 * time.Second):
				log.Println("Timeout. Failure. Closing.")
			}
			return
		}
	}

}
