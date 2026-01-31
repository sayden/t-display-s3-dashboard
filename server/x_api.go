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

type Tweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

type XAPIResponse struct {
	Data []Tweet `json:"data"`
}

func fetchTimeline(userID string) {
	// 1. Check Rate Limit
	if !canFetch() {
		log.Println("Skipping X API fetch due to rate limits.")
		return
	}

	// 2. Fetch from X API
	token := os.Getenv("X_BEARER_TOKEN")
	if token == "" {
		log.Println("X_BEARER_TOKEN is not set")
		return
	}

	// Using the reverse chronological home timeline as a proxy for the feed
	url := fmt.Sprintf("https://api.twitter.com/2/users/%s/timelines/reverse_chronological", userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch from X API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("X API returned status: %d", resp.StatusCode)
		// If 429, maybe update last fetch time to avoid hammering?
		// For now, only update on success or maybe just let the 8h interval handle it.
		// Actually, if we fail, we probably shouldn't block for 8 hours unless it was a rate limit error.
		// But let's keep it simple: update 'last_fetch' only on attempt to ensure we don't exceed quota.
		updateLastFetchTime()
		return
	}

	var apiResp XAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Printf("Failed to decode X API response: %v", err)
		return
	}

	// 3. Store in SQLite
	saveTweets(apiResp.Data)
	updateLastFetchTime()
	log.Printf("Successfully fetched %d tweets from X API", len(apiResp.Data))
}

func canFetch() bool {
	var lastFetchStr string
	err := db.QueryRow("SELECT value FROM meta WHERE key = 'last_fetch_time'").Scan(&lastFetchStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return true // Never fetched
		}
		log.Printf("Error check last fetch time: %v", err)
		return false // Fail safe
	}

	lastFetch, err := time.Parse(time.RFC3339, lastFetchStr)
	if err != nil {
		return true // Invalid date, retry
	}

	// 8 hours interval (approx 100 requests / month)
	// 30 days * 24h = 720h. 100 reqs -> 1 per 7.2h. Using 8h to be safe.
	return time.Since(lastFetch) > 8*time.Hour
}

func updateLastFetchTime() {
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec("INSERT OR REPLACE INTO meta (key, value) VALUES ('last_fetch_time', ?)", now)
	if err != nil {
		log.Printf("Failed to update last fetch time: %v", err)
	}
}

func saveTweets(tweets []Tweet) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return
	}

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO tweets (id, text, author_id, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return
	}
	defer stmt.Close()

	for _, t := range tweets {
		// API returns ISO8601 strings usually, but json decoder might handle if struct has string
		// Wait, my struct has time.Time. Default json decoder might fail if format isn't exact standard
		// X API v2 uses ISO 8601. Go's json decoder handles RFC3339 which is compatible.
		_, err := stmt.Exec(t.ID, t.Text, t.AuthorID, t.CreatedAt)
		if err != nil {
			log.Printf("Failed to insert tweet %s: %v", t.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
	}
}
