package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// handleMessage returns a placeholder message
func handleMessage(w http.ResponseWriter, r *http.Request) {
	// 2. Fetch from X API
	username := os.Getenv("X_USERNAME")
	if username == "" {
		log.Println("X_USERNAME is not set")
		return
	}

	// Trigger generic fetch (rate limited internally)
	// User ID from request: 110010111101011
	go fetchTimeline(username)

	// Get the most recent unsent tweet
	var tweetText string
	var tweetID string
	var authorID string

	var tweetCreatedAt time.Time
	// Select most recent unsent tweet including timestamp
	row := db.QueryRow("SELECT id, text, author_id, created_at FROM tweets WHERE sent_to_board = 0 ORDER BY created_at DESC LIMIT 1")
	// Scan created_at. Depending on driver, it might come as time.Time or string.
	// go-sqlite3 normally handles time.Time if DSN has parseTime=true, or we scan into interface/string.
	// db.go just used "sqlite3" driver without flags, so it's likely a string.
	var createdAtStr string
	err := row.Scan(&tweetID, &tweetText, &authorID, &createdAtStr)
	if err == nil {
		tweetCreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	}

	var msg Message
	if err != nil {
		if err == sql.ErrNoRows {
			msg = Message{
				Message: "No new tweets available.",
			}
		} else {
			log.Printf("Error querying tweet: %v", err)
			msg = Message{
				Message: "Error fetching tweet.",
			}
		}
	} else {
		// Mark as sent
		_, err := db.Exec("UPDATE tweets SET sent_to_board = 1 WHERE id = ?", tweetID)
		if err != nil {
			log.Printf("Failed to mark tweet %s as sent: %v", tweetID, err)
		}

		// Parse and format time
		timeFormatted := tweetCreatedAt.Format("15:04 02/01")

		msg = Message{
			Author:  authorID,
			Text:    tweetText,
			Time:    timeFormatted,
			Message: fmt.Sprintf("[%s]: %s", authorID, tweetText), // Keep backward compatibility
		}
	}

	json.NewEncoder(w).Encode(msg)
	log.Printf("[%s] GET /api/message -> %s", time.Now().Format("15:04:05"), msg.Message)
}
