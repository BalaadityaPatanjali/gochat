package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	if creds.Username == "" || creds.Password == "" {
		writeJSONError(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)",
		creds.Username, string(hashedPassword))
	if err != nil {
		writeJSONError(w, "User already exists or DB error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	var storedPass string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", creds.Username).Scan(&storedPass)
	if err == sql.ErrNoRows {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		writeJSONError(w, "DB error", http.StatusInternalServerError)
		return
	}

	// Compare password hash
	if bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(creds.Password)) != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(creds.Username)
	if err != nil {
		writeJSONError(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// helper for JSON errors
func writeJSONError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
