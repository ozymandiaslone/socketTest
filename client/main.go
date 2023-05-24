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
			log.Println("Conn error: ", err)
			log.Println("Unable to read from socket. Returning receiveHandler()...")
			return
		}
		log.Printf("Received: %s", message)
	}
}

func socketLoop(connection *websocket.Conn) {
	defer connection.Close()
	for {
		select {
		case <-time.After(15 * time.Millisecond):
			connection.WriteMessage(websocket.TextMessage, []byte("Test websocket message..."))
		case <-interrupt:
			log.Println("Interrupt detected. Terminating connection and program...")
			// The code to cose a websocket connection normally is fucked, so there's no way I can just remember it
			if err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Println("Error closing socket:", err)
				return
			}
			return
		}
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
	go socketLoop(conn)
	go receiveHandler(conn)

	<-done // just waits for the channel to populate i think
	log.Println("Done.")

}
