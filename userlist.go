package main

import (
	"encoding/json"
	"net/http"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err == nil {
			users = append(users, username)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
    "users": users,
    })
}
