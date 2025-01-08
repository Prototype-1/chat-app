package main

import (
	"fmt"
	"log"
	"net/http"
	"chat-app/routes"
	"chat-app/controllers"
)

func main() {
	routes.RegisterRoutes()
	go controllers.HandleMessage()

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}