package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/db"
	"ponziworld/backend/routes"
)

func TestNextDayEndpoint(t *testing.T) {
	// Create test dependencies
	deps := CreateTestDependencies("bank")
	defer CleanupTestDependencies(deps)

	// Ensure database indexes are created before running tests
	if err := db.EnsureAllIndexes(deps.DatabaseConfig); err != nil {
		t.Fatalf("Failed to ensure database indexes: %v", err)
	}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, deps)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create admin user for testing with unique username
	timestamp := time.Now().Unix()
	adminUsername := fmt.Sprintf("testadmin_%d", timestamp)
	adminToken, err := CreateAdminUserForTest(deps.DatabaseConfig, adminUsername, "password123", "TestAdminBank")
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

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
		// Create a regular (non-admin) user with unique username
		timestamp := time.Now().Unix()
		regularUsername := fmt.Sprintf("regularuser_%d", timestamp)
		regularToken, err := CreateRegularUserForTest(deps.DatabaseConfig, regularUsername, "password123", "RegularBank")
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

func TestCurrentDayEndpoint(t *testing.T) {
	// Create test dependencies
	deps := CreateTestDependencies("bank")
	defer CleanupTestDependencies(deps)

	// Ensure database indexes are created before running tests
	if err := db.EnsureAllIndexes(deps.DatabaseConfig); err != nil {
		t.Fatalf("Failed to ensure database indexes: %v", err)
	}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, deps)
	server := httptest.NewServer(mux)
	defer server.Close()

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
		// Create admin user for testing with unique username
		timestamp := time.Now().Unix()
		adminUsername := fmt.Sprintf("testadmin2_%d", timestamp)
		adminToken, err := CreateAdminUserForTest(deps.DatabaseConfig, adminUsername, "password123", "TestAdminBank2")
		if err != nil {
			t.Fatalf("Failed to create admin user: %v", err)
		}

		// Create a game state with day 5 by calling nextDay API endpoint
		for i := 0; i < 5; i++ {
			req, err := http.NewRequest("POST", server.URL+"/api/nextDay", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Failed to advance to day %d, status: %d", i+1, resp.StatusCode)
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
