package main

import (
	"log" // Added log import
	"math/rand"
	"time"
)

// UpdateStats applies time-based stat decay since last update
func UpdateStats(dog *Dog) {
	now := time.Now().UTC() // Changed to UTC
	elapsed := now.Sub(dog.LastUpdate)
	hours := elapsed.Hours()

	log.Printf("UpdateStats: LastUpdate=%v ScannedNow=%v Hours=%.4f",
		dog.LastUpdate, now, hours) // Added log statement

	if hours < 0.01 { // Less than ~36 seconds, skip update
		// Force update LastUpdate to prevent minor drifts if called frequently
		// But only if we are moving forward in time
		if hours > 0 {
			// dog.LastUpdate = now // Maybe better to not update if we didn't decay?
			// No, let's leave it to accumulate time
		}
		return
	}

	// Apply hunger decay
	hungerDecay := int(HungerDecayPerHour * hours)
	oldHunger := dog.Hunger // Added oldHunger
	dog.Hunger = Clamp(dog.Hunger-hungerDecay, 0, 100)

	log.Printf("Decay: Hunger %d -> %d (Decay: %d)", oldHunger, dog.Hunger, hungerDecay) // Added log statement

	// Apply happiness decay
	happinessDecay := int(HappinessDecayPerHour * hours)
	dog.Happiness = Clamp(dog.Happiness-happinessDecay, 0, 100)

	// Apply hygiene decay
	hygieneDecay := int(HygieneDecayPerHour * hours)
	dog.Hygiene = Clamp(dog.Hygiene-hygieneDecay, 0, 100)

	// Health decay when sick
	if dog.IsSick {
		healthDecay := int(HealthDecayWhenSick * hours)
		dog.Health = Clamp(dog.Health-healthDecay, 10, 100) // Min 10, dog can't die
	}

	// Random poop generation (roughly once every 2-4 hours)
	poopChance := hours / 3.0
	if rand.Float64() < poopChance {
		dog.PoopCount++
		dog.Hygiene = Clamp(dog.Hygiene-5, 0, 100)
	}

	// Check for sickness based on poor stats
	if !dog.IsSick {
		checkSickness(dog)
	}

	// Slow health recovery when not sick and stats are good
	if !dog.IsSick && dog.Hunger > 50 && dog.Hygiene > 50 && dog.Happiness > 30 {
		healthRegen := int(2.0 * hours)
		dog.Health = Clamp(dog.Health+healthRegen, 0, 100)
	}

	dog.LastUpdate = now
}

// checkSickness determines if the dog should become sick based on stats
func checkSickness(dog *Dog) {
	// Poor conditions increase sickness chance
	sicknessRisk := 0.0

	if dog.Hunger < 20 {
		sicknessRisk += 0.1
	}
	if dog.Hygiene < 20 {
		sicknessRisk += 0.15
	}
	if dog.PoopCount >= 3 {
		sicknessRisk += 0.1
	}
	if dog.Health < 50 {
		sicknessRisk += 0.1
	}

	if rand.Float64() < sicknessRisk {
		dog.IsSick = true
		dog.Happiness = Clamp(dog.Happiness-10, 0, 100)
	}
}

// Feed the dog with meal or snack
func Feed(dog *Dog, feedType string) string {
	switch feedType {
	case ActionFeedMeal:
		dog.Hunger = Clamp(dog.Hunger+FeedMealHunger, 0, 100)
		dog.Weight = ClampFloat(dog.Weight+FeedMealWeight, MinWeight, MaxWeight)
		return "Buddy enjoyed a tasty meal!"
	case ActionFeedSnack:
		dog.Hunger = Clamp(dog.Hunger+FeedSnackHunger, 0, 100)
		dog.Weight = ClampFloat(dog.Weight+FeedSnackWeight, MinWeight, MaxWeight)
		return "Buddy loved the snack!"
	default:
		dog.Hunger = Clamp(dog.Hunger+FeedSnackHunger, 0, 100)
		return "Buddy ate something."
	}
}

// Play with the dog
func Play(dog *Dog) string {
	if dog.IsSick {
		dog.Happiness = Clamp(dog.Happiness+5, 0, 100) // Less effect when sick
		return "Buddy tried to play but doesn't feel well..."
	}

	dog.Happiness = Clamp(dog.Happiness+PlayHappiness, 0, 100)
	dog.Weight = ClampFloat(dog.Weight+PlayWeight, MinWeight, MaxWeight)
	dog.Hunger = Clamp(dog.Hunger-5, 0, 100) // Playing makes dog hungry

	return "Buddy had fun playing!"
}

// Clean the dog (bath or poop cleanup)
func Clean(dog *Dog, cleanType string) string {
	switch cleanType {
	case ActionCleanBath:
		dog.Hygiene = Clamp(dog.Hygiene+BathHygiene, 0, 100)
		dog.Happiness = Clamp(dog.Happiness-5, 0, 100) // Dogs often don't love baths
		return "Buddy is squeaky clean!"
	case ActionCleanPoop:
		if dog.PoopCount > 0 {
			dog.PoopCount--
			dog.Hygiene = Clamp(dog.Hygiene+CleanPoopHygiene, 0, 100)
			return "You cleaned up after Buddy."
		}
		return "Nothing to clean up!"
	default:
		return "Nothing happened."
	}
}

// Discipline the dog (scold or praise)
func Discipline(dog *Dog, actionType string) string {
	switch actionType {
	case ActionScold:
		dog.Discipline = Clamp(dog.Discipline+ScoldDiscipline, 0, 100)
		dog.Happiness = Clamp(dog.Happiness-5, 0, 100)
		return "You scolded Buddy."
	case ActionPraise:
		dog.Discipline = Clamp(dog.Discipline+PraiseDiscipline, 0, 100)
		dog.Happiness = Clamp(dog.Happiness+3, 0, 100)
		return "Good boy, Buddy!"
	default:
		return "Nothing happened."
	}
}

// Cure gives medicine to the dog
func Cure(dog *Dog) string {
	if !dog.IsSick {
		return "Buddy is already healthy!"
	}

	dog.Health = Clamp(dog.Health+MedicineHealth, 0, 100)

	// Cure sickness if health is above threshold
	if dog.Health >= 50 {
		dog.IsSick = false
		return "Buddy feels much better now!"
	}

	return "The medicine helped a little."
}

// GetState determines the visual state of the dog for sprite selection
func GetState(dog *Dog) string {
	// Priority-based state selection (higher priority first)
	if dog.IsSick || dog.Health < 30 {
		return "sick"
	}
	if dog.Hygiene < 20 || dog.PoopCount >= 3 {
		return "dirty"
	}
	if dog.Hunger < 20 {
		return "hungry"
	}
	if dog.Happiness < 20 {
		return "sad"
	}
	if dog.Happiness > 70 && dog.Health > 70 {
		return "happy"
	}
	return "normal"
}

// CheckAttention determines if the dog needs attention
func CheckAttention(dog *Dog) bool {
	return dog.Hunger < 30 ||
		dog.Happiness < 30 ||
		dog.Hygiene < 30 ||
		dog.IsSick ||
		dog.PoopCount >= 2
}
