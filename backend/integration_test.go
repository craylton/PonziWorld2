package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestUserCreationAndLogin tests the complete user creation and login flow
func TestUserCreationAndLogin(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Generate unique username for this test
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("testuser_%d", timestamp)
	testBankName := fmt.Sprintf("Test Bank %d", timestamp)

	t.Run("Create User", func(t *testing.T) {
		// Test user creation
		createUserData := map[string]string{
			"username": testUsername,
			"bankName": testBankName,
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var createdUser User
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
		if createdUser.ClaimedCapital != 1000 {
			t.Errorf("Expected claimed capital 1000, got %d", createdUser.ClaimedCapital)
		}
		if createdUser.ActualCapital != 1000 {
			t.Errorf("Expected actual capital 1000, got %d", createdUser.ActualCapital)
		}

		t.Logf("Created user: %+v", createdUser)
	})

	t.Run("Login User", func(t *testing.T) {
		// Test user login
		loginData := map[string]string{
			"username": testUsername,
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

		var loggedInUser User
		if err := json.NewDecoder(resp.Body).Decode(&loggedInUser); err != nil {
			t.Fatalf("Failed to decode logged in user: %v", err)
		}

		// Verify user data from login
		if loggedInUser.Username != testUsername {
			t.Errorf("Expected username %s, got %s", testUsername, loggedInUser.Username)
		}
		if loggedInUser.BankName != testBankName {
			t.Errorf("Expected bank name %s, got %s", testBankName, loggedInUser.BankName)
		}
		if loggedInUser.ClaimedCapital != 1000 {
			t.Errorf("Expected claimed capital 1000, got %d", loggedInUser.ClaimedCapital)
		}
		if loggedInUser.ActualCapital != 1000 {
			t.Errorf("Expected actual capital 1000, got %d", loggedInUser.ActualCapital)
		}

		t.Logf("Logged in user: %+v", loggedInUser)
	})

	// Cleanup: Remove test user from database
	t.Cleanup(func() {
		client, ctx, cancel := ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)
		collection := client.Database("ponziworld").Collection("users")
		_, err := collection.DeleteOne(ctx, bson.M{"username": testUsername})
		if err != nil {
			t.Logf("Failed to cleanup test user: %v", err)
		}
	})
}

// TestExistingUserBackfill tests that existing users without capital fields get backfilled
func TestExistingUserBackfill(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Generate unique username for this test
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("olduser_%d", timestamp)
	testBankName := fmt.Sprintf("Old Bank %d", timestamp)
	// Manually insert a user without capital fields (simulating old user)
	client, ctx, cancel := ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database("ponziworld").Collection("users")

	oldUser := bson.M{
		"_id":      primitive.NewObjectID(),
		"username": testUsername,
		"bankName": testBankName,
		// Note: NOT including claimedCapital and actualCapital to simulate old user
	}
	_, err := collection.InsertOne(ctx, oldUser)
	if err != nil {
		t.Fatalf("Failed to insert old user: %v", err)
	}

	t.Run("Login Old User Should Backfill", func(t *testing.T) {
		// Test login for old user (should trigger backfill)
		loginData := map[string]string{
			"username": testUsername,
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to login old user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var loggedInUser User
		if err := json.NewDecoder(resp.Body).Decode(&loggedInUser); err != nil {
			t.Fatalf("Failed to decode logged in user: %v", err)
		}

		// Verify backfilled capital values
		if loggedInUser.ClaimedCapital != 1000 {
			t.Errorf("Expected backfilled claimed capital 1000, got %d", loggedInUser.ClaimedCapital)
		}
		if loggedInUser.ActualCapital != 1000 {
			t.Errorf("Expected backfilled actual capital 1000, got %d", loggedInUser.ActualCapital)
		}

		t.Logf("Backfilled user: %+v", loggedInUser)

		// Verify the database was actually updated
		var updatedUser User
		err = collection.FindOne(ctx, bson.M{"username": testUsername}).Decode(&updatedUser)
		if err != nil {
			t.Fatalf("Failed to find updated user: %v", err)
		}

		if updatedUser.ClaimedCapital != 1000 {
			t.Errorf("Database not updated: expected claimed capital 1000, got %d", updatedUser.ClaimedCapital)
		}
		if updatedUser.ActualCapital != 1000 {
			t.Errorf("Database not updated: expected actual capital 1000, got %d", updatedUser.ActualCapital)
		}
	})

	// Cleanup: Remove test user from database
	t.Cleanup(func() {
		_, err := collection.DeleteOne(ctx, bson.M{"username": testUsername})
		if err != nil {
			t.Logf("Failed to cleanup test user: %v", err)
		}
	})
}

// TestUserCreationDuplicateUsername tests that duplicate usernames are rejected
func TestUserCreationDuplicateUsername(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Generate unique username for this test
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("dupuser_%d", timestamp)
	testBankName := fmt.Sprintf("Dup Bank %d", timestamp)

	// Create first user
	createUserData := map[string]string{
		"username": testUsername,
		"bankName": testBankName,
	}
	jsonData, _ := json.Marshal(createUserData)

	resp1, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for first user, got %d", resp1.StatusCode)
	}

	// Try to create duplicate user
	resp2, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to attempt duplicate user creation: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for duplicate user, got %d", resp2.StatusCode)
	}

	// Cleanup: Remove test user from database
	t.Cleanup(func() {
		client, ctx, cancel := ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)
		collection := client.Database("ponziworld").Collection("users")
		_, err := collection.DeleteOne(ctx, bson.M{"username": testUsername})
		if err != nil {
			t.Logf("Failed to cleanup test user: %v", err)
		}
	})
}

// TestLoginNonExistentUser tests login with a user that doesn't exist
func TestLoginNonExistentUser(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Try to login with non-existent user
	loginData := map[string]string{
		"username": "nonexistentuser_123456789",
	}
	jsonData, _ := json.Marshal(loginData)

	resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to attempt login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for non-existent user, got %d", resp.StatusCode)
	}
}
