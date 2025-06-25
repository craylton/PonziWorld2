package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/auth"
	"ponziworld/backend/middleware"
	"ponziworld/backend/routes"
)

func TestJWTAuth(t *testing.T) {
	t.Run("Generate and validate token", func(t *testing.T) {
		username := "testuser"
		token, err := auth.GenerateToken(username)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		validatedUsername, err := auth.ValidateToken(token)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if validatedUsername != username {
			t.Errorf("Expected username %s, got %s", username, validatedUsername)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, err := auth.ValidateToken("invalid.token.here")
		if err == nil {
			t.Error("Expected error for invalid token")
		}
	})

	t.Run("Empty token", func(t *testing.T) {
		_, err := auth.ValidateToken("")
		if err == nil {
			t.Error("Expected error for empty token")
		}
	})
}

func TestJWTMiddleware(t *testing.T) {
	// Create a test handler that the middleware will protect
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("X-Username")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"username": username})
	})

	t.Run("Valid token", func(t *testing.T) {
		// Generate a valid token
		token, err := auth.GenerateToken("testuser")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		middleware.JWTMiddleware(testHandler)(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		var response map[string]string
		json.NewDecoder(rec.Body).Decode(&response)
		if response["username"] != "testuser" {
			t.Errorf("Expected username testuser, got %s", response["username"])
		}
	})

	t.Run("Missing Authorization header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()

		middleware.JWTMiddleware(testHandler)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})

	t.Run("Invalid Bearer format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Invalid token-format")
		rec := httptest.NewRecorder()

		middleware.JWTMiddleware(testHandler)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		rec := httptest.NewRecorder()

		middleware.JWTMiddleware(testHandler)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})

	t.Run("Empty Bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer ")
		rec := httptest.NewRecorder()

		middleware.JWTMiddleware(testHandler)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})
}

func TestLoginEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// First create a user to test login with
	createUserData := map[string]string{
		"username": "logintest_" + string(rune(time.Now().Unix())),
		"password": "testpassword123",
		"bankName": "Test Bank",
	}
	jsonData, _ := json.Marshal(createUserData)
	http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))

	t.Run("Valid login", func(t *testing.T) {
		loginData := map[string]string{
			"username": createUserData["username"],
			"password": createUserData["password"],
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		if _, ok := response["token"]; !ok {
			t.Error("Expected token in response")
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		loginData := map[string]string{
			"username": createUserData["username"],
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Missing username", func(t *testing.T) {
		loginData := map[string]string{
			"password": "testpassword123",
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Missing password", func(t *testing.T) {
		loginData := map[string]string{
			"username": createUserData["username"],
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Empty username", func(t *testing.T) {
		loginData := map[string]string{
			"username": "",
			"password": "testpassword123",
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to attempt login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer([]byte("{invalid json")))
		if err != nil {
			t.Fatalf("Failed to attempt login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}
