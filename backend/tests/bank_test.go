package tests

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestBankService_GetBankByUsername(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("bank")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("banktest_%d", timestamp)
	testBankName := "Test Bank API"
	testPassword := "testpassword123"

	// Create player directly via service
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create new player: %v", err)
	}

	// Retrieve bank directly via service
	bankResponse, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get bank by username: %v", err)
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
	asset := bankResponse.Assets[0]
	if asset.AssetType != "Cash" {
		t.Errorf("Expected asset type 'Cash', got %q", asset.AssetType)
	}
	if asset.Amount != 1000 {
		t.Errorf("Expected asset amount 1000, got %d", asset.Amount)
	}
	if asset.AssetTypeId == "" {
		t.Error("Expected asset type ID to be present")
	}
}
