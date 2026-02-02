package main

import (
	"encoding/json"
	"testing"
)

func TestGetActivities(t *testing.T) {
	t.Skip()

	i := intervals{}
	activities, err := i.GetActivities()
	if err != nil {
		t.Errorf("Failed to get activities: %v", err)
	}

	byt, _ := json.MarshalIndent(activities, "", "  ")
	t.Logf("Fitness: %s", string(byt))
}

func TestGetFitness(t *testing.T) {
	t.Skip()

	i := intervals{}
	fitness, err := i.GetFitness("2026-02-02")
	if err != nil {
		t.Errorf("Failed to get fitness: %v", err)
	}

	byt, _ := json.MarshalIndent(fitness, "", "  ")
	t.Logf("Fitness: %s", string(byt))
}
