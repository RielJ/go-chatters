package database

import (
	"database/sql"
	"fmt"
)

func InitDB(db *sql.DB) {
	createUserTable(db)
	createConversationTable(db)
	createMessageTable(db)
	createConversationUserTable(db)
}

// Create User Table
func createUserTable(db *sql.DB) {
	fmt.Println("Creating User Table")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT,
		first_name TEXT,
		last_name TEXT,
		email TEXT,
		password TEXT
	)`)
	if err != nil {
		fmt.Println("Error creating user table")
		panic(err)
	}
}

// Create Conversation Table
func createConversationTable(db *sql.DB) {
	fmt.Println("Creating Conversation Table")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS conversations (
		id SERIAL PRIMARY KEY,
		name TEXT
	)`)
	if err != nil {
		panic(err)
	}
}

// Create Message Table
func createMessageTable(db *sql.DB) {
	fmt.Println("Creating Message Table")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		user_id INTEGER,
		message TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		conversation_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(conversation_id) REFERENCES conversations(id)
	)`)
	if err != nil {
		panic(err)
	}
}

// Create Conversation User Table
func createConversationUserTable(db *sql.DB) {
	fmt.Println("Creating Conversation User Table")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS conversation_users (
		id SERIAL PRIMARY KEY,
		conversation_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY(conversation_id) REFERENCES conversations(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)`)
	if err != nil {
		panic(err)
	}
}
