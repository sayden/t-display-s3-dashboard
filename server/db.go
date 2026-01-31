package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./dashboard.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTables()
}

func createTables() {
	queryTweets := `
	CREATE TABLE IF NOT EXISTS tweets (
		id TEXT PRIMARY KEY,
		text TEXT,
		author_id TEXT,
		created_at DATETIME,
		sent_to_board INTEGER DEFAULT 0
	);`

	if _, err := db.Exec(queryTweets); err != nil {
		log.Fatalf("Failed to create tweets table: %v", err)
	}

	queryMeta := `
	CREATE TABLE IF NOT EXISTS meta (
		key TEXT PRIMARY KEY,
		value TEXT
	);`

	if _, err := db.Exec(queryMeta); err != nil {
		log.Fatalf("Failed to create meta table: %v", err)
	}
}
