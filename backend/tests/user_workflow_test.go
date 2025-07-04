package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/db"
	"ponziworld/backend/models"
	"ponziworld/backend/routes"
)

// TestFullUserWorkflow tests the complete end-to-end player workflow
func TestFullUserWorkflow(t *testing.T) {
	// Create test dependencies
	deps := CreateTestDependencies("bank")
	defer CleanupTestDependencies(deps)
	
	// Ensure database indexes are created before running tests
	if err := db.EnsureAllIndexes(deps.DatabaseConfig); err != nil {
		t.Fatalf("Failed to ensure database indexes: %v", err)
	}
	
	// Setup test server
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, deps)
	server := httptest.NewServer(mux)
	defer server.Close()
	var defaultCapital int64 = 1000

	// Generate unique username for this test
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("workflow_%d", timestamp)
	testBankName := fmt.Sprintf("Workflow Bank %d", timestamp)
	testPassword := "testpassword123"
	var authToken string

	t.Run("Create Player", func(t *testing.T) {
		// arrange
		createUserData := map[string]string{
			"username": testUsername,
			"password": testPassword,
			"bankName": testBankName,
		}
		jsonData, err := json.Marshal(createUserData)
		if err != nil {
			t.Fatalf("Failed to marshal create player data: %v", err)
		}

		// act
		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))

		// assert
		if err != nil {
			t.Fatalf("Failed to create player: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})
	
	t.Run("Login Player", func(t *testing.T) {
		// Test player login
		loginData := map[string]string{
			"username": testUsername,
			"password": testPassword,
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to login player: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Parse the response to get the JWT
		var loginResponse map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		token, ok := loginResponse["token"].(string)
		if !ok {
			t.Fatalf("Login response did not contain token")
		}
		authToken = token
	})
	
	t.Run("Get Player Details", func(t *testing.T) {
			// GET /api/newPlayer is removed; expect Method Not Allowed
			req, err := http.NewRequest("GET", server.URL+"/api/newPlayer", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+authToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to fetch player after login: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d from /api/newPlayer, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
			}
        	// Test /api/bank endpoint
		req, _ = http.NewRequest("GET", server.URL+"/api/bank", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to fetch bank info: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 from /api/bank, got %d", resp.StatusCode)
		}

		var bankResponse models.BankResponse
		if err := json.NewDecoder(resp.Body).Decode(&bankResponse); err != nil {
			t.Fatalf("Failed to decode bank response: %v", err)
		}

		// Verify bank data from /api/bank
		if bankResponse.BankName != testBankName {
			t.Errorf("Expected bank name %s, got %s", testBankName, bankResponse.BankName)
		}
		if bankResponse.ClaimedCapital != defaultCapital {
			t.Errorf("Expected claimed capital %d, got %d", defaultCapital, bankResponse.ClaimedCapital)
		}
		if bankResponse.ActualCapital != defaultCapital {
			t.Errorf("Expected actual capital %d, got %d", defaultCapital, bankResponse.ActualCapital)
		}
		if len(bankResponse.Assets) != 1 {
			t.Errorf("Expected 1 asset, got %d", len(bankResponse.Assets))
		}
		if len(bankResponse.Assets) > 0 {
			if bankResponse.Assets[0].AssetType != "Cash" {
				t.Errorf("Expected asset type 'Cash', got %s", bankResponse.Assets[0].AssetType)
			}
			if bankResponse.Assets[0].Amount != defaultCapital {
				t.Errorf("Expected asset amount %d, got %d", defaultCapital, bankResponse.Assets[0].Amount)
			}
		}
	})
}
