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

	"ponziworld/backend/auth"
	"ponziworld/backend/db"
	"ponziworld/backend/models"
	"ponziworld/backend/routes"
)

func TestBankEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("banktest_%d", timestamp)
	testBankName := "Test Bank API"

	// Step 1: Create a test user
	createUserData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
		"bankName": testBankName,
	}
	jsonData, _ := json.Marshal(createUserData)

	resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 for user creation, got %d", resp.StatusCode)
	}

	// Step 2: Login to get JWT token
	loginData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
	}
	jsonData, _ = json.Marshal(loginData)

	resp, err = http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for login, got %d", resp.StatusCode)
	}

	var loginResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	token, ok := loginResponse["token"].(string)
	if !ok {
		t.Fatalf("Login response did not contain token")
	}

	t.Run("Valid bank request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/bank", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to fetch bank: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var bankResponse models.BankResponse
		if err := json.NewDecoder(resp.Body).Decode(&bankResponse); err != nil {
			t.Fatalf("Failed to decode bank response: %v", err)
		}

		// Verify bank data
		if bankResponse.BankName != testBankName {
			t.Errorf("Expected bank name %s, got %s", testBankName, bankResponse.BankName)
		}
		if bankResponse.ClaimedCapital != 1000 {
			t.Errorf("Expected claimed capital 1000, got %d", bankResponse.ClaimedCapital)
		}
		if bankResponse.ActualCapital != 1000 {
			t.Errorf("Expected actual capital 1000, got %d", bankResponse.ActualCapital)
		}
		if len(bankResponse.Assets) != 1 {
			t.Errorf("Expected 1 asset, got %d", len(bankResponse.Assets))
		}
		if len(bankResponse.Assets) > 0 {
			asset := bankResponse.Assets[0]
			if asset.AssetType != "Cash" {
				t.Errorf("Expected asset type 'Cash', got %s", asset.AssetType)
			}
			if asset.Amount != 1000 {
				t.Errorf("Expected asset amount 1000, got %d", asset.Amount)
			}
		}
	})

	t.Run("No Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/bank", nil)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/bank", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Token for non-existent user", func(t *testing.T) {		// Create a token for a user that doesn't exist
		fakeToken, _ := auth.GenerateToken("nonexistentuser")
		req, _ := http.NewRequest("GET", server.URL+"/api/bank", nil)
		req.Header.Set("Authorization", "Bearer "+fakeToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	// Cleanup
	t.Cleanup(func() {
		client, ctx, cancel := db.ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)

		// Clean up users collection
		usersCollection := client.Database("ponziworld").Collection("users")
		usersCollection.DeleteMany(ctx, bson.M{"username": testUsername})

		// Clean up banks collection
		banksCollection := client.Database("ponziworld").Collection("banks")
		banksCollection.DeleteMany(ctx, bson.M{})

		// Clean up assets collection
		assetsCollection := client.Database("ponziworld").Collection("assets")
		assetsCollection.DeleteMany(ctx, bson.M{})
	})
}
