package angels

import (
	"errors"
	"testing"

	"github.com/joho/godotenv"
)

const (
	MAIN_PATH = "../../cmd/main/.env"
)

/*
func TestGetTeams(t *testing.T) {
	godotenv.Load(MAIN_PATH)
	BuildCronnables()
	t.Log("COOL")
}
*/

func TestWinConditions(t *testing.T) {
	godotenv.Load(MAIN_PATH)
	id := 153122
	win, _ := CheckWinConditions(id)
	if !win {
		t.Fatalf("Game with id %d should evaluate to a win", id)
	}
}

// Test is only good until the game with id 160571 is played
func TestRescheduleError(t *testing.T) {
	godotenv.Load(MAIN_PATH)
	id := 160571
	_, err := CheckWinConditions(id)
	if !errors.Is(err, RESCHEDULE_ERROR) {
		t.Fatal("Should be throwing an error as this game hasnt started yet")
	}
}
