package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"ponziworld/backend/db"
	"ponziworld/backend/models"
	"ponziworld/backend/routes"
)

// TestFullUserWorkflow tests the complete end-to-end user workflow
func TestFullUserWorkflow(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()
	var defaultCapital int64 = 1000

	// Generate unique username for this test
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("workflow_%d", timestamp)
	testBankName := fmt.Sprintf("Workflow Bank %d", timestamp)
	testPassword := "testpassword123"
	var authToken string

	t.Run("Create User", func(t *testing.T) {
		// arrange
		createUserData := map[string]string{
			"username": testUsername,
			"password": testPassword,
			"bankName": testBankName,
		}
		jsonData, err := json.Marshal(createUserData)
		if err != nil {
			t.Fatalf("Failed to marshal create user data: %v", err)
		}

		// act
		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))

		// assert
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})
	
	t.Run("Login User", func(t *testing.T) {
		// Test user login
		loginData := map[string]string{
			"username": testUsername,
			"password": testPassword,
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to login user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Parse the response to get the JWT token
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
	
	t.Run("Get User Details", func(t *testing.T) {
			// GET /api/user is removed; expect Method Not Allowed
			req, err := http.NewRequest("GET", server.URL+"/api/user", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+authToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to fetch user after login: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d from /api/user, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
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

	// Cleanup: Remove test user from database
	t.Cleanup(func() {
		client, ctx, cancel := db.ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)
		collection := client.Database("ponziworld").Collection("users")
		_, err := collection.DeleteOne(ctx, bson.M{"username": testUsername})
		if err != nil {
			t.Logf("Failed to cleanup test user: %v", err)
		}
	})
}
