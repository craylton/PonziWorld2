package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/handlers"
	"ponziworld/backend/models"
)

func TestBankHandlerDirect(t *testing.T) {
	// Create test dependencies
	deps := CreateTestDependencies("bank_direct")
	defer CleanupTestDependencies(deps)

	// Create bank handler
	bankHandler := handlers.NewBankHandler(deps)

	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("banktest_%d", timestamp)
	testBankName := "Test Bank Direct"

	// Create a test player and bank directly using the service manager
	ctx := context.Background() // Create a fresh context for database operations
	serviceManager := deps.ServiceManager

	// Create player
	err := serviceManager.Player.CreateNewPlayer(ctx, testUsername, "testpassword123", testBankName)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}

	t.Run("Valid bank request", func(t *testing.T) {
		// Create a test request
		req := httptest.NewRequest("GET", "/api/bank", nil)
		req.Header.Set("X-Username", testUsername) // Simulate JWT middleware

		// Create a test response recorder
		w := httptest.NewRecorder()

		// Call the handler
		bankHandler.GetBank(w, req)

		// Check response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var bankResponse models.BankResponse
		if err := json.NewDecoder(w.Body).Decode(&bankResponse); err != nil {
			t.Fatalf("Failed to decode bank response: %v", err)
		}

		// Verify bank data
		if bankResponse.BankName != testBankName {
			t.Errorf("Expected bank name %q, got %q", testBankName, bankResponse.BankName)
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
				t.Errorf("Expected asset type 'Cash', got %q", asset.AssetType)
			}
			if asset.Amount != 1000 {
				t.Errorf("Expected asset amount 1000, got %d", asset.Amount)
			}
		}
	})

	t.Run("Missing username header", func(t *testing.T) {
		// Create a test request without username header
		req := httptest.NewRequest("GET", "/api/bank", nil)

		// Create a test response recorder
		w := httptest.NewRecorder()

		// Call the handler
		bankHandler.GetBank(w, req)

		// Check response
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Non-existent user", func(t *testing.T) {
		// Create a test request with non-existent username
		req := httptest.NewRequest("GET", "/api/bank", nil)
		req.Header.Set("X-Username", "nonexistent_user")

		// Create a test response recorder
		w := httptest.NewRecorder()

		// Call the handler
		bankHandler.GetBank(w, req)

		// Check response
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}
