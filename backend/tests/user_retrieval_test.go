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
	"ponziworld/backend/routes"
)

func TestGetUserEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create a test user first
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("gettest_%d", timestamp)
	testPassword := "testpassword123"
	
	// Create user through API to ensure consistency
	createUserData := map[string]string{
		"username": testUsername,
		"password": testPassword,
		"bankName": "Test Bank",
	}
	jsonData, _ := json.Marshal(createUserData)
	http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))

	// Generate valid token
	validToken, _ := auth.GenerateToken(testUsername)

	t.Cleanup(func() {
		client, ctx, cancel := db.ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)
		collection := client.Database("ponziworld").Collection("users")
		collection.DeleteOne(ctx, bson.M{"username": testUsername})
	})

	t.Run("Valid authenticated request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/user", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		var userResponse map[string]string
		json.NewDecoder(resp.Body).Decode(&userResponse)
		if userResponse["username"] != testUsername {
			t.Errorf("Expected username %s, got %s", testUsername, userResponse["username"])
		}
	})

	t.Run("No Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/user", nil)

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
		req, _ := http.NewRequest("GET", server.URL+"/api/user", nil)
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

	t.Run("Token for non-existent user", func(t *testing.T) {
		// Generate token for user that doesn't exist
		nonExistentToken, _ := auth.GenerateToken("nonexistent_user_12345")
		
		req, _ := http.NewRequest("GET", server.URL+"/api/user", nil)
		req.Header.Set("Authorization", "Bearer "+nonExistentToken)

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

	t.Run("Malformed Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/user", nil)
		req.Header.Set("Authorization", "NotBearer "+validToken)

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
}
