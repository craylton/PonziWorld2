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
	"ponziworld/backend/routes"
)

func TestUserCreation(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("createtest_%d", timestamp)

	t.Run("Valid user creation", func(t *testing.T) {
		createUserData := map[string]string{
			"username": testUsername,
			"password": "testpassword123",
			"bankName": "Test Bank Creation",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	// Cleanup
	t.Cleanup(func() {
		client, ctx, cancel := db.ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)
		collection := client.Database("ponziworld").Collection("users")
		collection.DeleteOne(ctx, bson.M{"username": testUsername})
	})

	t.Run("Duplicate username", func(t *testing.T) {
		// Try to create the same user again
		createUserData := map[string]string{
			"username": testUsername,
			"password": "testpassword123",
			"bankName": "Another Bank",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt duplicate user creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Missing username", func(t *testing.T) {
		createUserData := map[string]string{
			"password": "testpassword123",
			"bankName": "Test Bank",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt user creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Missing password", func(t *testing.T) {
		createUserData := map[string]string{
			"username": "testuser_missing_password",
			"bankName": "Test Bank",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt user creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Missing bank name", func(t *testing.T) {
		createUserData := map[string]string{
			"username": "testuser_missing_bank",
			"password": "testpassword123",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt user creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Empty username", func(t *testing.T) {
		createUserData := map[string]string{
			"username": "",
			"password": "testpassword123",
			"bankName": "Test Bank",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt user creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Empty password", func(t *testing.T) {
		createUserData := map[string]string{
			"username": "testuser_empty_password",
			"password": "",
			"bankName": "Test Bank",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt user creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Whitespace-only username", func(t *testing.T) {
		createUserData := map[string]string{
			"username": "   ",
			"password": "testpassword123",
			"bankName": "Test Bank",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt user creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer([]byte("{invalid json")))
		if err != nil {
			t.Fatalf("Failed to attempt user creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}
