package main

import (
	"encoding/json"
	"net/http"
)

func historyHandler(w http.ResponseWriter, r *http.Request) {
	user1 := r.URL.Query().Get("user1")
	user2 := r.URL.Query().Get("user2")

	if user1 == "" || user2 == "" {
		writeJSONError(w, "Missing user1 or user2", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT id, sender, receiver, content, created_at
		FROM messages
		WHERE (sender = ? AND receiver = ?)
		   OR (sender = ? AND receiver = ?)
		ORDER BY created_at ASC
	`, user1, user2, user2, user1)
	if err != nil {
		writeJSONError(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var history []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Receiver, &msg.Content, &msg.CreatedAt); err == nil {
			history = append(history, msg)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
