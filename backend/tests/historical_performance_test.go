package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHistoricalPerformanceService_GetHistoricalPerformance(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("histPerf")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("perftest_%d", timestamp)
	testBankName := "Test Bank Performance"
	testPassword := "testpassword123"

	// Create player and bank directly via service
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create new player: %v", err)
	}

	// Get bank details to get bank ID
	bankResponse, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get bank: %v", err)
	}

	bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
	if err != nil {
		t.Fatalf("Failed to convert bank ID to ObjectID: %v", err)
	}

	// Test performance history service directly
	historyResponse, err := container.ServiceContainer.HistoricalPerformance.GetHistoricalPerformance(ctx, testUsername, bankID)
	if err != nil {
		t.Fatalf("Failed to get performance history: %v", err)
	}

	// Verify that we get 30 days of history
	if len(historyResponse.ClaimedHistory) != 30 {
		t.Fatalf("Expected 30 days of claimed history, got %d", len(historyResponse.ClaimedHistory))
	}

	// Since actual history should contain the initial entry for a new bank, we expect 1 entry
	// For a newly created bank, this should be 1 (the initial £1000 entry)
	if len(historyResponse.ActualHistory) != 1 {
		t.Fatalf("Expected 1 day of actual history for new bank, got %d", len(historyResponse.ActualHistory))
	}

	// Verify the initial actual history entry is correct
	if len(historyResponse.ActualHistory) > 0 {
		initialEntry := historyResponse.ActualHistory[0]
		if initialEntry.Day != 0 {
			t.Fatalf("Expected initial actual history entry to be for day 0, got day %d", initialEntry.Day)
		}
		if initialEntry.Value != 1000 {
			t.Fatalf("Expected initial actual history entry to have value 1000, got %d", initialEntry.Value)
		}
	}

	// Verify that all entries are properly ordered by day
	for i := 1; i < len(historyResponse.ClaimedHistory); i++ {
		if historyResponse.ClaimedHistory[i].Day <= historyResponse.ClaimedHistory[i-1].Day {
			t.Fatal("Claimed history is not properly ordered by day")
		}
	}

	// For actual history, only verify if we have entries
	for i := 1; i < len(historyResponse.ActualHistory); i++ {
		if historyResponse.ActualHistory[i].Day <= historyResponse.ActualHistory[i-1].Day {
			t.Fatal("Actual history is not properly ordered by day")
		}
	}

	// Verify that claimed and actual history have the same values for the days that exist
	// Since actual history might not exist for all days, we can't assume they're the same length
	for i := 0; i < len(historyResponse.ActualHistory); i++ {
		found := false
		for j := 0; j < len(historyResponse.ClaimedHistory); j++ {
			if historyResponse.ClaimedHistory[j].Day == historyResponse.ActualHistory[i].Day {
				if historyResponse.ClaimedHistory[j].Value != historyResponse.ActualHistory[i].Value {
					t.Fatalf("Claimed and actual history values don't match for day %d", historyResponse.ActualHistory[i].Day)
				}
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Actual history day %d not found in claimed history", historyResponse.ActualHistory[i].Day)
		}
	}
}

func TestHistoricalPerformanceService_GetHistoricalPerformanceInvalidBankID(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("histPerf")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("perfinvalidtest_%d", timestamp)
	testBankName := "Test Bank Invalid"
	testPassword := "testpassword123"

	// Create player and bank
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	// Test with invalid bank ID - this should return error for bank not found
	invalidBankID := primitive.NewObjectID()
	_, err = container.ServiceContainer.HistoricalPerformance.GetHistoricalPerformance(ctx, testUsername, invalidBankID)

	// Should return error since the bank doesn't exist
	if err == nil {
		t.Fatal("Expected error for invalid bank ID, got nil")
	}
}

func TestHistoricalPerformanceService_GetHistoricalPerformanceDataPersistence(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("histPerf")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("perfpersist_%d", timestamp)
	testBankName := "Test Bank Persistence"
	testPassword := "testpassword123"

	// Create player and bank
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	// Get bank details to get bank ID
	bankResponse, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get bank: %v", err)
	}

	bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
	if err != nil {
		t.Fatalf("Failed to convert bank ID to ObjectID: %v", err)
	}

	// First call to performance history service
	firstResponse, err := container.ServiceContainer.HistoricalPerformance.GetHistoricalPerformance(ctx, testUsername, bankID)
	if err != nil {
		t.Fatalf("Failed to get performance history first time: %v", err)
	}

	// Second call to performance history service (should return identical data)
	secondResponse, err := container.ServiceContainer.HistoricalPerformance.GetHistoricalPerformance(ctx, testUsername, bankID)
	if err != nil {
		t.Fatalf("Failed to get performance history second time: %v", err)
	}

	// Verify that both responses are identical (data persisted in database)
	if len(firstResponse.ClaimedHistory) != len(secondResponse.ClaimedHistory) {
		t.Fatalf("Claimed history length differs between calls: %d vs %d",
			len(firstResponse.ClaimedHistory), len(secondResponse.ClaimedHistory))
	}

	if len(firstResponse.ActualHistory) != len(secondResponse.ActualHistory) {
		t.Fatalf("Actual history length differs between calls: %d vs %d",
			len(firstResponse.ActualHistory), len(secondResponse.ActualHistory))
	}

	for i := 0; i < len(firstResponse.ClaimedHistory); i++ {
		if firstResponse.ClaimedHistory[i].Day != secondResponse.ClaimedHistory[i].Day ||
			firstResponse.ClaimedHistory[i].Value != secondResponse.ClaimedHistory[i].Value {
			t.Fatalf("Claimed history differs at index %d: first=(%d,%d), second=(%d,%d)",
				i, firstResponse.ClaimedHistory[i].Day, firstResponse.ClaimedHistory[i].Value,
				secondResponse.ClaimedHistory[i].Day, secondResponse.ClaimedHistory[i].Value)
		}
	}

	for i := 0; i < len(firstResponse.ActualHistory); i++ {
		if firstResponse.ActualHistory[i].Day != secondResponse.ActualHistory[i].Day ||
			firstResponse.ActualHistory[i].Value != secondResponse.ActualHistory[i].Value {
			t.Fatalf("Actual history differs at index %d: first=(%d,%d), second=(%d,%d)",
				i, firstResponse.ActualHistory[i].Day, firstResponse.ActualHistory[i].Value,
				secondResponse.ActualHistory[i].Day, secondResponse.ActualHistory[i].Value)
		}
	}
}

func TestHistoricalPerformanceService_GetHistoricalPerformanceUnauthorized(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("histPerf")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()

	// Create first player and bank (owner)
	ownerUsername := fmt.Sprintf("perfowner_%d", timestamp)
	ownerBankName := "Owner Bank"
	ownerPassword := "testpassword123"

	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, ownerUsername, ownerPassword, ownerBankName)
	if err != nil {
		t.Fatalf("Failed to create owner player: %v", err)
	}

	// Create second player (unauthorized user)
	unauthorizedUsername := fmt.Sprintf("perfunauth_%d", timestamp)
	unauthorizedBankName := "Unauthorized Bank"
	unauthorizedPassword := "testpassword123"

	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, unauthorizedUsername, unauthorizedPassword, unauthorizedBankName)
	if err != nil {
		t.Fatalf("Failed to create unauthorized player: %v", err)
	}

	// Get the owner's bank details to get bank ID
	ownerBankResponse, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, ownerUsername)
	if err != nil {
		t.Fatalf("Failed to get owner bank: %v", err)
	}

	ownerBankID, err := primitive.ObjectIDFromHex(ownerBankResponse.Id)
	if err != nil {
		t.Fatalf("Failed to convert owner bank ID to ObjectID: %v", err)
	}

	// Try to get historical performance for owner's bank using unauthorized user's credentials
	_, err = container.ServiceContainer.HistoricalPerformance.GetHistoricalPerformance(ctx, unauthorizedUsername, ownerBankID)

	// Should return error since the user doesn't own the bank
	if err == nil {
		t.Fatal("Expected error for unauthorized access to bank, got nil")
	}

	// The error should be related to unauthorized access
	// Based on the bank service code, it should return ErrUnauthorized
	expectedError := "unauthorized access"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
