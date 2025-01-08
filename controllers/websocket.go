package controllers

import (
	"log"
	"net/http"
	"sync"
	"time"
	"chat-app/models"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan models.Message)      
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mu sync.Mutex
)

// Handles incoming WebSocket requests
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	leaveProcessed := make(chan bool)

	defer func() {
		mu.Lock()
		if username, exists := clients[conn]; exists && username != "" {
			log.Printf("Broadcasting leave message for user: %s", username)
			broadcast <- models.Message{
				Username: "Server",
				Action:   "broadcast",
				Content:  username + " has left the chat!",
			}
		}
		delete(clients, conn)
		mu.Unlock()

		log.Println("Closing connection")
		conn.Close()
		close(leaveProcessed)
	}()

	mu.Lock()
	clients[conn] = ""
	mu.Unlock()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error: %v\n", err)
			return
		}

		switch msg.Action {
		case "join":
			mu.Lock()
			clients[conn] = msg.Username
			mu.Unlock()
			log.Printf("User joined: %s", msg.Username)
			broadcast <- models.Message{
				Username: "Server",
				Action:   "broadcast",
				Content:  msg.Username + " has joined the chat!",
			}
		case "leave":
			log.Printf("Received leave request from user: %s", msg.Username)
			broadcast <- models.Message{
				Username: msg.Username,
				Action:   "broadcast",
				Content:  msg.Username + " has left the chat!",
			}

			time.Sleep(500 * time.Millisecond) 
			return 
		case "message":
			broadcast <- models.Message{
				Username: msg.Username,
				Action:   "message",
				Content:  msg.Content,
			}
		}
	}
}

//Handles broadcast to multiple clients
func HandleMessage() {
	for {
		msg := <-broadcast
		log.Printf("Broadcasting message: %+v", msg)
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
