package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/auth"
	"ponziworld/backend/middleware"
	"ponziworld/backend/requestcontext"
	"ponziworld/backend/services"
)

func TestJwtAuth(t *testing.T) {
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

func TestJwtMiddleware(t *testing.T) {
	// Create a test container with mock services
	container, err := CreateTestDependencies("jwt_middleware")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	// Create a test handler that the middleware will protect
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _ := requestcontext.UsernameFromContext(r.Context())
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

		middleware.JwtMiddleware(testHandler, container.Logger)(rec, req)

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

		middleware.JwtMiddleware(testHandler, container.Logger)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})

	t.Run("Invalid Bearer format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Invalid token-format")
		rec := httptest.NewRecorder()

		middleware.JwtMiddleware(testHandler, container.Logger)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		rec := httptest.NewRecorder()

		middleware.JwtMiddleware(testHandler, container.Logger)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})

	t.Run("Empty Bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer ")
		rec := httptest.NewRecorder()

		middleware.JwtMiddleware(testHandler, container.Logger)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})
}

// Unit test for AuthService.Login without HTTP
func TestAuthService_Login(t *testing.T) {
	container, err := CreateTestDependencies("auth")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("authtest_%d", timestamp)
	password := "testpassword123"

	// Create a new player via PlayerService
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, username, password, "Test Bank")
	if err != nil {
		t.Fatalf("Failed to create new player: %v", err)
	}

	t.Run("Valid login", func(t *testing.T) {
		player, err := container.ServiceContainer.Auth.Login(ctx, username, password)
		if err != nil {
			t.Fatalf("Expected successful login, got error: %v", err)
		}
		if player.Username != username {
			t.Errorf("Expected username %s, got %s", username, player.Username)
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		_, err := container.ServiceContainer.Auth.Login(ctx, username, "wrongpassword")
		if !errors.Is(err, services.ErrInvalidCredentials) {
			t.Errorf("Expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("Missing username", func(t *testing.T) {
		_, err := container.ServiceContainer.Auth.Login(ctx, "", password)
		if err == nil {
			t.Error("Expected error for missing username, got nil")
		}
	})

	t.Run("Missing password", func(t *testing.T) {
		_, err := container.ServiceContainer.Auth.Login(ctx, username, "")
		if err == nil {
			t.Error("Expected error for missing password, got nil")
		}
	})
}
