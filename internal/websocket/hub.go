package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rbb-market-go/internal/db"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("Client connected")

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Println("Client disconnected")
			}

		case message := <-h.Broadcast:
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}

			err := saveMessageToDB(msg.From, msg.To, msg.Content)
			if err != nil {
				log.Println("Error saving message:", err)
			}

			for client := range h.Clients {
				if client.ID == msg.To || client.ID == msg.From {
					client.Send <- message
				}
			}
		}
	}
}

func saveMessageToDB(fromUserID, toUserID int, content string) error {
	query := `INSERT INTO messages ("from", "to", content) VALUES ($1, $2, $3)`

	_, err := db.DB.Exec(context.Background(), query, fromUserID, toUserID, content)
	if err != nil {
		return fmt.Errorf("unable to save message: %v", err)
	}
	return nil
}
