package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/routes"
)

func TestGameService_NextDay(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("game")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()

	t.Run("should create initial day 0 and increment to day 1", func(t *testing.T) {
		day, err := container.ServiceContainer.Game.NextDay(ctx)
		if err != nil {
			t.Fatalf("Failed to advance to next day: %v", err)
		}

		if day != 1 {
			t.Errorf("Expected day to be 1, got %d", day)
		}
	})

	t.Run("should increment existing day", func(t *testing.T) {
		day, err := container.ServiceContainer.Game.NextDay(ctx)
		if err != nil {
			t.Fatalf("Failed to advance to next day: %v", err)
		}

		if day != 2 {
			t.Errorf("Expected day to be 2, got %d", day)
		}
	})
}

func TestNextDayEndpoint(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("game")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, container)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should reject non-admin users", func(t *testing.T) {
		// Create a regular (non-admin) user with unique username
		timestamp := time.Now().Unix()
		regularUsername := fmt.Sprintf("regularuser_%d", timestamp)
		regularToken, err := CreateRegularUserForTest(container, regularUsername, "password123", "RegularBank")
		if err != nil {
			t.Fatal("Failed to create regular user:", err)
		}

		req, err := http.NewRequest("POST", server.URL+"/api/nextDay", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+regularToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
		}
	})
}

func TestGameService_CurrentDay(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("game")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()

	t.Run("should return day 0 when no game state exists", func(t *testing.T) {
		currentDay, err := container.ServiceContainer.Game.GetCurrentDay(ctx)
		if err != nil {
			t.Fatalf("Failed to get current day: %v", err)
		}

		if currentDay != 0 {
			t.Errorf("Expected currentDay to be 0, got %d", currentDay)
		}
	})

	t.Run("should return current day when game state exists", func(t *testing.T) {
		// Advance the game to day 5 by calling NextDay service directly
		var finalDay int
		for i := range 5 {
			day, err := container.ServiceContainer.Game.NextDay(ctx)
			if err != nil {
				t.Fatalf("Failed to advance to day %d: %v", i+1, err)
			}
			finalDay = day
		}

		if finalDay != 5 {
			t.Errorf("Expected final day to be 5, got %d", finalDay)
		}

		currentDay, err := container.ServiceContainer.Game.GetCurrentDay(ctx)
		if err != nil {
			t.Fatalf("Failed to get current day: %v", err)
		}

		if currentDay != 5 {
			t.Errorf("Expected currentDay to be 5, got %d", currentDay)
		}
	})
}
