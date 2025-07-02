package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"ponziworld/backend/db"
	"ponziworld/backend/routes"
	"ponziworld/backend/services"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestNextDayEndpoint(t *testing.T) {
	// Ensure database indexes are created before running tests
	if err := db.EnsureAllIndexes(); err != nil {
		t.Fatalf("Failed to ensure database indexes: %v", err)
	}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Clean up game collection before test
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("ponziworld")
	db.Collection("game").DeleteMany(ctx, bson.M{})

	// Create admin user for testing
	adminToken, err := CreateAdminUserForTest("testadmin", "password123", "TestAdminBank")
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}
	defer CleanupTestData("testadmin", "TestAdminBank")

	t.Run("should create initial day 0 and increment to day 1", func(t *testing.T) {
		req, err := http.NewRequest("POST", server.URL+"/api/nextDay", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// Read the error response for debugging
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status code %d, got %d. Response: %s", http.StatusOK, resp.StatusCode, string(bodyBytes))
			return
		}

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response["currentDay"] != 1 {
			t.Errorf("Expected currentDay to be 1, got %d", response["currentDay"])
		}
	})

	t.Run("should increment existing day", func(t *testing.T) {
		req, err := http.NewRequest("POST", server.URL+"/api/nextDay", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			return
		}

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response["currentDay"] != 2 {
			t.Errorf("Expected currentDay to be 2, got %d", response["currentDay"])
		}
	})

	t.Run("should reject non-POST methods", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL+"/api/nextDay", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})

	t.Run("should reject non-admin users", func(t *testing.T) {
		// Create a regular (non-admin) user
		regularToken, err := CreateRegularUserForTest("regularuser", "password123", "RegularBank")
		if err != nil {
			t.Fatal("Failed to create regular user:", err)
		}
		defer CleanupTestData("regularuser", "RegularBank")

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

func TestGameService(t *testing.T) {
	// Ensure database indexes are created before running tests
	if err := db.EnsureAllIndexes(); err != nil {
		t.Fatalf("Failed to ensure database indexes: %v", err)
	}

	// Clean up game collection before test
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	database := client.Database("ponziworld")
	database.Collection("game").DeleteMany(ctx, bson.M{})

	serviceManager := services.NewServiceManager(database)

	t.Run("should return 0 for initial day", func(t *testing.T) {
		day, err := serviceManager.Game.GetCurrentDay(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if day != 0 {
			t.Errorf("Expected initial day to be 0, got %d", day)
		}
	})

	t.Run("should increment day correctly", func(t *testing.T) {
		day, err := serviceManager.Game.NextDay(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if day != 1 {
			t.Errorf("Expected day to be 1 after increment, got %d", day)
		}

		// Verify with GetCurrentDay
		currentDay, err := serviceManager.Game.GetCurrentDay(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if currentDay != 1 {
			t.Errorf("Expected current day to be 1, got %d", currentDay)
		}
	})

	t.Run("should increment day multiple times", func(t *testing.T) {
		// Increment to day 2
		day, err := serviceManager.Game.NextDay(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if day != 2 {
			t.Errorf("Expected day to be 2, got %d", day)
		}

		// Increment to day 3
		day, err = serviceManager.Game.NextDay(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if day != 3 {
			t.Errorf("Expected day to be 3, got %d", day)
		}
	})
}
