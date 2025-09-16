package main

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./chatapp.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO messages (sender, receiver, content) VALUES ('debug_sender','debug_receiver','hello from init')")
if err != nil {
    log.Println("❌ Debug insert failed:", err)
} else {
    log.Println("✅ Debug insert worked")
}


	// Print the exact DB path being used
	absPath, _ := filepath.Abs("./chatapp.db")
	log.Println("📂 Using DB file at:", absPath)

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender TEXT NOT NULL,
		receiver TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createMessagesTable)
	if err != nil {
		log.Fatal("Error creating messages table:", err)
	}
}
