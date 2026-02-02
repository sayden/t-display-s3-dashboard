package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// handleTamagotchi returns the current dog state with sprite
func handleTamagotchi(w http.ResponseWriter, r *http.Request) {
	dog, err := GetDog()
	if err != nil {
		log.Printf("Error getting dog: %v", err)
		http.Error(w, "Failed to get dog state", http.StatusInternalServerError)
		return
	}

	// Update stats based on time elapsed
	UpdateStats(dog)

	// Save updated state
	if err := SaveDog(dog); err != nil {
		log.Printf("Error saving dog: %v", err)
	}

	// Determine visual state and get sprite
	state := GetState(dog)
	image, width, height := GetSprite(state)

	response := GameResponse{
		Dog:            *dog,
		State:          state,
		NeedsAttention: CheckAttention(dog),
		Image:          image,
		ImgWidth:       width,
		ImgHeight:      height,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] GET /api/tamagotchi -> %s (State:%s, H:%d, Ha:%d, Hy:%d, He:%d)",
		time.Now().Format("15:04:05"), dog.Name, state,
		dog.Hunger, dog.Happiness, dog.Hygiene, dog.Health)
}

// handleFeed handles feeding the dog
func handleFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	feedType := r.URL.Query().Get("type")
	if feedType == "" {
		feedType = ActionFeedMeal
	}

	dog, err := GetDog()
	if err != nil {
		http.Error(w, "Failed to get dog", http.StatusInternalServerError)
		return
	}

	UpdateStats(dog)
	message := Feed(dog, feedType)

	if err := SaveDog(dog); err != nil {
		http.Error(w, "Failed to save dog", http.StatusInternalServerError)
		return
	}

	sendGameResponse(w, dog, message)
	log.Printf("[%s] POST /api/tamagotchi/feed?type=%s -> %s",
		time.Now().Format("15:04:05"), feedType, message)
}

// handlePlay handles playing with the dog
func handlePlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dog, err := GetDog()
	if err != nil {
		http.Error(w, "Failed to get dog", http.StatusInternalServerError)
		return
	}

	UpdateStats(dog)
	message := Play(dog)

	if err := SaveDog(dog); err != nil {
		http.Error(w, "Failed to save dog", http.StatusInternalServerError)
		return
	}

	sendGameResponse(w, dog, message)
	log.Printf("[%s] POST /api/tamagotchi/play -> %s",
		time.Now().Format("15:04:05"), message)
}

// handleClean handles cleaning the dog or poop
func handleClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cleanType := r.URL.Query().Get("type")
	if cleanType == "" {
		cleanType = ActionCleanPoop
	}

	dog, err := GetDog()
	if err != nil {
		http.Error(w, "Failed to get dog", http.StatusInternalServerError)
		return
	}

	UpdateStats(dog)
	message := Clean(dog, cleanType)

	if err := SaveDog(dog); err != nil {
		http.Error(w, "Failed to save dog", http.StatusInternalServerError)
		return
	}

	sendGameResponse(w, dog, message)
	log.Printf("[%s] POST /api/tamagotchi/clean?type=%s -> %s",
		time.Now().Format("15:04:05"), cleanType, message)
}

// handleDiscipline handles disciplining the dog
func handleDiscipline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	actionType := r.URL.Query().Get("type")
	if actionType == "" {
		actionType = ActionPraise
	}

	dog, err := GetDog()
	if err != nil {
		http.Error(w, "Failed to get dog", http.StatusInternalServerError)
		return
	}

	UpdateStats(dog)
	message := Discipline(dog, actionType)

	if err := SaveDog(dog); err != nil {
		http.Error(w, "Failed to save dog", http.StatusInternalServerError)
		return
	}

	sendGameResponse(w, dog, message)
	log.Printf("[%s] POST /api/tamagotchi/discipline?type=%s -> %s",
		time.Now().Format("15:04:05"), actionType, message)
}

// handleCure handles giving medicine to the dog
func handleCure(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dog, err := GetDog()
	if err != nil {
		http.Error(w, "Failed to get dog", http.StatusInternalServerError)
		return
	}

	UpdateStats(dog)
	message := Cure(dog)

	if err := SaveDog(dog); err != nil {
		http.Error(w, "Failed to save dog", http.StatusInternalServerError)
		return
	}

	sendGameResponse(w, dog, message)
	log.Printf("[%s] POST /api/tamagotchi/cure -> %s",
		time.Now().Format("15:04:05"), message)
}

// handleReset resets the dog to a new game
func handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Buddy"
	}

	dog, err := ResetDog(name)
	if err != nil {
		http.Error(w, "Failed to reset dog", http.StatusInternalServerError)
		return
	}

	sendGameResponse(w, dog, "Welcome your new friend: "+name+"!")
	log.Printf("[%s] POST /api/tamagotchi/reset -> New dog: %s",
		time.Now().Format("15:04:05"), name)
}

// sendGameResponse sends a full game response with sprite
func sendGameResponse(w http.ResponseWriter, dog *Dog, message string) {
	state := GetState(dog)
	image, width, height := GetSprite(state)

	response := GameResponse{
		Dog:            *dog,
		State:          state,
		NeedsAttention: CheckAttention(dog),
		Message:        message,
		Image:          image,
		ImgWidth:       width,
		ImgHeight:      height,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
