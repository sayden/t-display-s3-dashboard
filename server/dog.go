package main

import (
	"time"
)

// Dog represents the virtual pet state
type Dog struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Hunger     int       `json:"hunger"`     // 0-100, 100 = full
	Happiness  int       `json:"happiness"`  // 0-100, 100 = very happy
	Hygiene    int       `json:"hygiene"`    // 0-100, 100 = clean
	Discipline int       `json:"discipline"` // 0-100, 100 = well-trained
	Weight     float64   `json:"weight"`     // 1.0-10.0, 5.0 = normal
	Health     int       `json:"health"`     // 0-100, 100 = healthy
	IsSick     bool      `json:"is_sick"`
	PoopCount  int       `json:"poop_count"`  // Number of poops to clean
	LastUpdate time.Time `json:"last_update"` // For stat decay calculation
	CreatedAt  time.Time `json:"created_at"`
}

// GameResponse is the JSON response for the tamagotchi endpoint
type GameResponse struct {
	Dog            Dog    `json:"dog"`
	State          string `json:"state"`             // Visual state: "happy", "normal", "sick", etc.
	NeedsAttention bool   `json:"needs_attention"`   // Attention call active
	Message        string `json:"message,omitempty"` // Action feedback message
	Image          string `json:"image"`             // Base64 RGB565 sprite data
	ImgWidth       int    `json:"img_width"`
	ImgHeight      int    `json:"img_height"`
}

// Action types for game interactions
const (
	ActionFeedMeal  = "meal"
	ActionFeedSnack = "snack"
	ActionCleanBath = "bath"
	ActionCleanPoop = "poop"
	ActionScold     = "scold"
	ActionPraise    = "praise"
)

// Stat decay rates (per hour)
const (
	HungerDecayPerHour    = 5.0
	HappinessDecayPerHour = 3.0
	HygieneDecayPerHour   = 4.0
	HealthDecayWhenSick   = 10.0 // Per hour when sick
)

// Action effects
const (
	FeedMealHunger   = 20
	FeedMealWeight   = 0.5
	FeedSnackHunger  = 10
	FeedSnackWeight  = 0.2
	PlayHappiness    = 15
	PlayWeight       = -0.3
	BathHygiene      = 40
	CleanPoopHygiene = 10
	ScoldDiscipline  = 10
	PraiseDiscipline = 5
	MedicineHealth   = 30
)

// Weight limits
const (
	MinWeight    = 1.0
	MaxWeight    = 10.0
	NormalWeight = 5.0
)

// NewDog creates a new dog with default stats
func NewDog(name string) *Dog {
	now := time.Now().UTC()
	return &Dog{
		Name:       name,
		Hunger:     80,
		Happiness:  80,
		Hygiene:    80,
		Discipline: 50,
		Weight:     NormalWeight,
		Health:     100,
		IsSick:     false,
		PoopCount:  0,
		LastUpdate: now,
		CreatedAt:  now,
	}
}

// Clamp ensures a value stays within min and max bounds
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat ensures a float value stays within min and max bounds
func ClampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
