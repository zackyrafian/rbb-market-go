package main

import (
	"log"
	"net/http"
	"rbb-market-go/internal/db"
	"rbb-market-go/internal/websocket"
)

func main() {
	db.Connect()
	defer db.Close()

	hub := websocket.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", websocket.ServeWS(hub))

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
