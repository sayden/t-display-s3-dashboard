package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Tamagotchi represents the response for the tamagotchi endpoint (placeholder)
type Tamagotchi struct {
	Name   string `json:"name"`
	Hunger int    `json:"hunger"`
	Happy  int    `json:"happy"`
	Energy int    `json:"energy"`
}

// handleTamagotchi returns placeholder tamagotchi state
func handleTamagotchi(w http.ResponseWriter, r *http.Request) {
	tamagotchi := Tamagotchi{
		Name:   "Pixel",
		Hunger: 75,
		Happy:  80,
		Energy: 90,
	}

	err := json.NewEncoder(w).Encode(tamagotchi)
	if err != nil {
		log.Printf("Error encoding tamagotchi: %v", err)
		http.Error(w, "Failed to encode tamagotchi", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] GET /api/tamagotchi -> %s (H:%d, Ha:%d, E:%d)",
		time.Now().Format("15:04:05"), tamagotchi.Name, tamagotchi.Hunger, tamagotchi.Happy, tamagotchi.Energy)
}
