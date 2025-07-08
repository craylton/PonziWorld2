package tests

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPlayerService_CreateNewPlayer(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("player")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("createtest_%d", timestamp)

	t.Run("Valid player creation", func(t *testing.T) {
		err := container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, "testpassword123", "Test Bank Creation")
		if err != nil {
			t.Fatalf("Failed to create player: %v", err)
		}

		// Verify the player was created by trying to get their bank
		bankResponse, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, testUsername)
		if err != nil {
			t.Fatalf("Failed to get bank for created player: %v", err)
		}

		if bankResponse.BankName != "Test Bank Creation" {
			t.Errorf("Expected bank name 'Test Bank Creation', got %q", bankResponse.BankName)
		}
		if bankResponse.ClaimedCapital != 1000 {
			t.Errorf("Expected claimed capital 1000, got %d", bankResponse.ClaimedCapital)
		}
	})

	t.Run("Duplicate username", func(t *testing.T) {
		// Try to create the same player again
		err := container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, "testpassword123", "Another Bank")
		if err == nil {
			t.Error("Expected error for duplicate username, got nil")
		}
	})

	t.Run("Missing username", func(t *testing.T) {
		err := container.ServiceContainer.Player.CreateNewPlayer(ctx, "", "testpassword123", "Test Bank")
		if err == nil {
			t.Error("Expected error for missing username, got nil")
		}
	})

	t.Run("Missing password", func(t *testing.T) {
		err := container.ServiceContainer.Player.CreateNewPlayer(ctx, "testuser_missing_password", "", "Test Bank")
		if err == nil {
			t.Error("Expected error for missing password, got nil")
		}
	})

	t.Run("Missing bank name", func(t *testing.T) {
		err := container.ServiceContainer.Player.CreateNewPlayer(ctx, "testuser_missing_bank", "testpassword123", "")
		if err == nil {
			t.Error("Expected error for missing bank name, got nil")
		}
	})

	t.Run("Whitespace-only username", func(t *testing.T) {
		err := container.ServiceContainer.Player.CreateNewPlayer(ctx, "   ", "testpassword123", "Test Bank")
		if err == nil {
			t.Error("Expected error for whitespace-only username, got nil")
		}
	})

	t.Run("Whitespace-only password", func(t *testing.T) {
		err := container.ServiceContainer.Player.CreateNewPlayer(ctx, "testuser_whitespace_password", "   ", "Test Bank")
		if err == nil {
			t.Error("Expected error for whitespace-only password, got nil")
		}
	})

	t.Run("Whitespace-only bank name", func(t *testing.T) {
		err := container.ServiceContainer.Player.CreateNewPlayer(ctx, "testuser_whitespace_bank", "testpassword123", "   ")
		if err == nil {
			t.Error("Expected error for whitespace-only bank name, got nil")
		}
	})
}
