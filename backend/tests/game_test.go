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

func TestCurrentDayEndpoint(t *testing.T) {
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

	t.Run("should return day 0 when no game state exists", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL+"/api/currentDay", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status code %d, got %d. Response: %s", http.StatusOK, resp.StatusCode, string(bodyBytes))
			return
		}

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response["currentDay"] != 0 {
			t.Errorf("Expected currentDay to be 0, got %d", response["currentDay"])
		}
	})

	t.Run("should return current day when game state exists", func(t *testing.T) {
		// Create a game state with day 5
		serviceManager := services.NewServiceManager(db)
		_, err := serviceManager.Game.NextDay(ctx) // Creates day 1
		if err != nil {
			t.Fatal(err)
		}
		// Advance to day 5
		for i := 0; i < 4; i++ {
			_, err = serviceManager.Game.NextDay(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}

		req, err := http.NewRequest("GET", server.URL+"/api/currentDay", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status code %d, got %d. Response: %s", http.StatusOK, resp.StatusCode, string(bodyBytes))
			return
		}

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response["currentDay"] != 5 {
			t.Errorf("Expected currentDay to be 5, got %d", response["currentDay"])
		}
	})
}
