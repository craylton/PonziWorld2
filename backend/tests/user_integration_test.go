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

// TestUserCreationAndLogin tests the complete user creation and login flow
func TestUserCreationAndLogin(test *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()
	var defaultCapital int64 = 1000

	// Generate unique username for this test
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("testuser_%d", timestamp)
	testBankName := fmt.Sprintf("Test Bank %d", timestamp)
	testPassword := "testpassword123"
	var authToken string

	test.Run("Create User", func(t *testing.T) {
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

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var createdUser models.User
		if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
			t.Fatalf("Failed to decode created user: %v", err)
		}

		// Verify user data
		if createdUser.Username != testUsername {
			t.Errorf("Expected username %s, got %s", testUsername, createdUser.Username)
		}
		if createdUser.BankName != testBankName {
			t.Errorf("Expected bank name %s, got %s", testBankName, createdUser.BankName)
		}
		if createdUser.ClaimedCapital != defaultCapital {
			t.Errorf("Expected claimed capital %d, got %d", defaultCapital, createdUser.ClaimedCapital)
		}
		if createdUser.ActualCapital != defaultCapital {
			t.Errorf("Expected actual capital %d, got %d", defaultCapital, createdUser.ActualCapital)
		}
	})
	test.Run("Login User", func(t *testing.T) {
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
		var loginResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		token, ok := loginResponse["token"].(string)
		if !ok {
			t.Fatalf("Login response did not contain token")
		}
		authToken = token
	})
	test.Run("Get User Details", func(t *testing.T) {
		// Now fetch the user object with authentication
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
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 from /api/user, got %d", resp.StatusCode)
		}
		var loggedInUser models.User
		if err := json.NewDecoder(resp.Body).Decode(&loggedInUser); err != nil {
			t.Fatalf("Failed to decode logged in user: %v", err)
		}

		// Verify user data from /api/user
		if loggedInUser.Username != testUsername {
			t.Errorf("Expected username %s, got %s", testUsername, loggedInUser.Username)
		}
		if loggedInUser.BankName != testBankName {
			t.Errorf("Expected bank name %s, got %s", testBankName, loggedInUser.BankName)
		}
		if loggedInUser.ClaimedCapital != defaultCapital {
			t.Errorf("Expected claimed capital %d, got %d", defaultCapital, loggedInUser.ClaimedCapital)
		}
		if loggedInUser.ActualCapital != defaultCapital {
			t.Errorf("Expected actual capital %d, got %d", defaultCapital, loggedInUser.ActualCapital)
		}
	})

	// Cleanup: Remove test user from database
	test.Cleanup(func() {
		client, ctx, cancel := db.ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)
		collection := client.Database("ponziworld").Collection("users")
		_, err := collection.DeleteOne(ctx, bson.M{"username": testUsername})
		if err != nil {
			test.Logf("Failed to cleanup test user: %v", err)
		}
	})
}

// TestUserCreationDuplicateUsername tests that duplicate usernames are rejected
func TestUserCreationDuplicateUsername(test *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Generate unique username for this test
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("dupuser_%d", timestamp)
	testBankName := fmt.Sprintf("Dup Bank %d", timestamp)
	// Create first user
	createUserData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
		"bankName": testBankName,
	}
	jsonData, err := json.Marshal(createUserData)
	if err != nil {
		test.Fatalf("Failed to marshal create user data: %v", err)
	}

	resp1, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		test.Fatalf("Failed to create first user: %v", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		test.Errorf("Expected status 200 for first user, got %d", resp1.StatusCode)
	}

	// Try to create duplicate user
	resp2, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		test.Fatalf("Failed to attempt duplicate user creation: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusBadRequest {
		test.Errorf("Expected status 400 for duplicate user, got %d", resp2.StatusCode)
	}

	// Cleanup: Remove test user from database
	test.Cleanup(func() {
		client, ctx, cancel := db.ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)
		collection := client.Database("ponziworld").Collection("users")
		_, err := collection.DeleteOne(ctx, bson.M{"username": testUsername})
		if err != nil {
			test.Logf("Failed to cleanup test user: %v", err)
		}
	})
}

// TestLoginNonExistentUser tests login with a user that doesn't exist
func TestLoginNonExistentUser(test *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()
	// Try to login with non-existent user
	loginData := map[string]string{
		"username": "nonexistentuser_123456789",
		"password": "somepassword",
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		test.Fatalf("Failed to marshal login data: %v", err)
	}

	resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		test.Fatalf("Failed to attempt login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		test.Errorf("Expected status 401 for non-existent user, got %d", resp.StatusCode)
	}
}
