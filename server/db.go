package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the SQLite database
func InitDB() error {
	// Get user data directory for persistent storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	dbDir := filepath.Join(homeDir, ".tamagotchi")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Printf("Warning: Could not create data directory: %v", err)
		dbDir = "."
	}

	dbPath := filepath.Join(dbDir, "dog.db")
	log.Printf("Database path: %s", dbPath)

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Create dogs table
	schema := `
	CREATE TABLE IF NOT EXISTS dogs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		hunger INTEGER DEFAULT 80,
		happiness INTEGER DEFAULT 80,
		hygiene INTEGER DEFAULT 80,
		discipline INTEGER DEFAULT 50,
		weight REAL DEFAULT 5.0,
		health INTEGER DEFAULT 100,
		is_sick BOOLEAN DEFAULT FALSE,
		poop_count INTEGER DEFAULT 0,
		last_update DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	// Create intervals cache table
	intervalsSchema := `
	CREATE TABLE IF NOT EXISTS intervals_cache (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		ctl REAL,
		atl REAL,
		ramp_rate REAL,
		fatigue REAL,
		stress REAL,
		activities_json TEXT,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(intervalsSchema)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// GetDog retrieves the current dog from the database, creating one if none exists
func GetDog() (*Dog, error) {
	dog := &Dog{}

	row := db.QueryRow(`
		SELECT id, name, hunger, happiness, hygiene, discipline, weight, health, 
		       is_sick, poop_count, last_update, created_at
		FROM dogs
		ORDER BY id DESC
		LIMIT 1
	`)

	var lastUpdate, createdAt string
	err := row.Scan(
		&dog.ID, &dog.Name, &dog.Hunger, &dog.Happiness, &dog.Hygiene,
		&dog.Discipline, &dog.Weight, &dog.Health, &dog.IsSick, &dog.PoopCount,
		&lastUpdate, &createdAt,
	)

	if err == sql.ErrNoRows {
		log.Println("GetDog: No dog found, creating new one")
		dog = NewDog("Buddy")
		err = SaveDog(dog)
		if err != nil {
			return nil, err
		}
		return dog, nil
	}

	if err != nil {
		return nil, err
	}

	// Parse timestamps - try RFC3339 first, fallback to old format for migration
	t, err := time.Parse(time.RFC3339, lastUpdate)
	if err != nil {
		// Try old format
		t, _ = time.Parse("2006-01-02 15:04:05", lastUpdate)
	}
	dog.LastUpdate = t

	t, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		t, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	}
	dog.CreatedAt = t

	log.Printf("GetDog: ID=%d Name=%s Hunger=%d LastUpdate=%v",
		dog.ID, dog.Name, dog.Hunger, dog.LastUpdate)

	return dog, nil
}

// SaveDog persists the dog state to the database
func SaveDog(dog *Dog) error {
	log.Printf("SaveDog: ID=%d Hunger=%d LastUpdate=%v", dog.ID, dog.Hunger, dog.LastUpdate)

	if dog.ID == 0 {
		// Insert new dog
		result, err := db.Exec(`
			INSERT INTO dogs (name, hunger, happiness, hygiene, discipline, weight, 
			                  health, is_sick, poop_count, last_update, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, dog.Name, dog.Hunger, dog.Happiness, dog.Hygiene, dog.Discipline,
			dog.Weight, dog.Health, dog.IsSick, dog.PoopCount,
			dog.LastUpdate.Format(time.RFC3339),
			dog.CreatedAt.Format(time.RFC3339))

		if err != nil {
			return err
		}

		id, _ := result.LastInsertId()
		dog.ID = int(id)
	} else {
		// Update existing dog
		_, err := db.Exec(`
			UPDATE dogs SET 
				name = ?, hunger = ?, happiness = ?, hygiene = ?, discipline = ?,
				weight = ?, health = ?, is_sick = ?, poop_count = ?, last_update = ?
			WHERE id = ?
		`, dog.Name, dog.Hunger, dog.Happiness, dog.Hygiene, dog.Discipline,
			dog.Weight, dog.Health, dog.IsSick, dog.PoopCount,
			dog.LastUpdate.Format(time.RFC3339), dog.ID)

		if err != nil {
			return err
		}
	}

	return nil
}

// ResetDog deletes the current dog and starts fresh
func ResetDog(name string) (*Dog, error) {
	// Delete all dogs
	_, err := db.Exec("DELETE FROM dogs")
	if err != nil {
		return nil, err
	}

	// Create new dog
	dog := NewDog(name)
	err = SaveDog(dog)
	if err != nil {
		return nil, err
	}

	log.Printf("Created new dog: %s", name)
	return dog, nil
}
