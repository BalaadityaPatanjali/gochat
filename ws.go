package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader with CORS origin check
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        return origin == "http://localhost:8080" || 
               origin == "http://127.0.0.1:8080" || 
               origin == "http://localhost:3000" ||
               origin == "https://myfrontend.com"
    },
}

// Keep track of connected clients
var clients = make(map[string]*websocket.Conn)

// Incoming message format from frontend
type IncomingMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

// WebSocket handler
func handleWS(w http.ResponseWriter, r *http.Request) {
	// --- Authenticate user via JWT ---
	tokenStr := r.URL.Query().Get("token")
	claims, err := validateJWT(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user := claims.Username

	// --- Upgrade to WebSocket ---
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[user] = conn
	log.Println("✅", user, "connected via WebSocket")

	// --- Listen for messages from this client ---
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Read error:", err)
			delete(clients, user)
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(msgBytes, &incoming); err != nil {
			log.Println("❌ Invalid message format:", err)
			continue
		}

		// --- Save message in DB ---
		_, err = db.Exec(
			"INSERT INTO messages (sender, receiver, content) VALUES (?, ?, ?)",
			user, incoming.To, incoming.Content,
		)
		if err != nil {
			log.Printf("❌ DB insert error: %v | sender=%s receiver=%s content=%s\n",
				err, user, incoming.To, incoming.Content)
			continue
		} else {
			log.Printf("💾 Message saved: %s → %s | %s\n", user, incoming.To, incoming.Content)
		}

		// --- Deliver instantly if recipient is online ---
		if recipientConn, ok := clients[incoming.To]; ok {
			outgoing := Message{
				Sender:    user,
				Receiver:  incoming.To,
				Content:   incoming.Content,
				CreatedAt: time.Now(),
			}
			outBytes, _ := json.Marshal(outgoing)
			if err := recipientConn.WriteMessage(websocket.TextMessage, outBytes); err != nil {
				log.Println("❌ Error delivering message:", err)
			} else {
				log.Printf("📩 Delivered to %s\n", incoming.To)
			}
		} else {
			log.Printf("📭 %s is offline, message stored for later\n", incoming.To)
		}
	}
}
