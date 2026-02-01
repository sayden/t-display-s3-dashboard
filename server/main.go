package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Message represents the response for the message endpoint
type Message struct {
	Message string `json:"message"`
}

// CORS middleware to allow requests from any origin
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// handleHealth returns server health status
func handleHealth(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("Error encoding health: %v", err)
		http.Error(w, "Failed to encode health", http.StatusInternalServerError)
		return
	}
}

func main() {
	port := ":8081"

	// Register routes
	http.HandleFunc("/api/message", corsMiddleware(handleMessage))
	http.HandleFunc("/api/weather", corsMiddleware(handleWeather))
	http.HandleFunc("/api/tamagotchi", corsMiddleware(handleTamagotchi))
	http.HandleFunc("/health", corsMiddleware(handleHealth))

	// Print startup info
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Dashboard Server for Lilygo T-Display-S3          ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Server running on http://0.0.0.0%s                    ║\n", port)
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Println("║  Available endpoints:                                      ║")
	fmt.Println("║    GET /api/message    - Returns a text message            ║")
	fmt.Println("║    GET /api/weather    - Returns weather data              ║")
	fmt.Println("║    GET /api/tamagotchi - Returns tamagotchi state          ║")
	fmt.Println("║    GET /health         - Returns server health             ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the server")
	fmt.Println()

	// Start server
	log.Fatal(http.ListenAndServe(port, nil))
}
