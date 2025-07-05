package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/routes"
)

func TestPlayerCreation(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("player")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, container)
	server := httptest.NewServer(mux)
	defer server.Close()

	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("createtest_%d", timestamp)

	t.Run("Valid player creation", func(t *testing.T) {
		createUserData := map[string]string{
			"username": testUsername,
			"password": "testpassword123",
			"bankName": "Test Bank Creation",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create player: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("Duplicate username", func(t *testing.T) {
		// Try to create the same player again
		createUserData := map[string]string{
			"username": testUsername,
			"password": "testpassword123",
			"bankName": "Another Bank",
		}
		jsonData, _ := json.Marshal(createUserData)

		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt duplicate player creation: %v", err)
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

		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt player creation: %v", err)
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

		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt player creation: %v", err)
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

		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt player creation: %v", err)
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

		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt player creation: %v", err)
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

		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt player creation: %v", err)
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

		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt player creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer([]byte("{invalid json")))
		if err != nil {
			t.Fatalf("Failed to attempt player creation: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}
