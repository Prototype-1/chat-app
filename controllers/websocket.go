package controllers

import (
	"log"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
	"chat-app/models"
)

var (
	clients = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.Message)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request)bool {
			return true
		},
	}
	mu sync.Mutex
)

//Handles WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error: %v\n", err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}
		broadcast <- msg
	}
}

//Broadcasts messages to all  clients
func HandleMessage() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Write error: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}