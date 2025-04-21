package websocket

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDParam := r.URL.Query().Get("user_id")
		if userIDParam == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDParam)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		client := &Client{
			ID:   userID,
			Conn: conn,
			Send: make(chan []byte),
			Hub:  hub,
		}

		hub.Register <- client

		go client.readPump()
		go client.writePump()
	}
}
