package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// handleMessage returns a placeholder message
func handleMessage(w http.ResponseWriter, r *http.Request) {
	msg := Message{
		Message: "Hello from Go Server! 🚀",
	}

	json.NewEncoder(w).Encode(msg)
	log.Printf("[%s] GET /api/message -> %s", time.Now().Format("15:04:05"), msg.Message)
}
