package main

import "time"

// Represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // hide in JSON responses
}

// Represents a chat message
type Message struct {
    ID        int       `json:"id"`
    Sender    string    `json:"sender"`
    Receiver  string    `json:"receiver"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"timestamp"`
}


