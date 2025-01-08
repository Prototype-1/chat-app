package routes

import (
	"net/http"
	"chat-app/controllers"
)

func RegisterRoutes() {
	http.HandleFunc("/ws", controllers.HandleWebSocket)
}