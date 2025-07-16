package tests

import (
	"context"
	"fmt"
	"ponziworld/backend/models"
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

func TestBankService_CreateBankForUsername(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("bank_create")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("bankcreatetest_%d", timestamp)
	testPassword := "testpass123"
	testBankName := fmt.Sprintf("Test Bank %d", timestamp)
	
	// Create test user and bank
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	
	// Create a new bank for the user
	newBank, err := container.ServiceContainer.Bank.CreateBankForUsername(
		ctx, 
		testUsername, 
		"Additional Bank", 
		100000,
	)
	if err != nil {
		t.Fatalf("Failed to create additional bank: %v", err)
	}
	
	if newBank.BankName != "Additional Bank" {
		t.Errorf("Expected bank name 'Additional Bank', got '%s'", newBank.BankName)
	}
	
	if newBank.ClaimedCapital != 100000 {
		t.Errorf("Expected claimed capital 100000, got %d", newBank.ClaimedCapital)
	}
	
	// Verify the bank was created in the database
	banks, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get banks: %v", err)
	}
	
	if len(banks) != 2 {
		t.Errorf("Expected 2 banks, got %d", len(banks))
	}
	
	// Sort banks to have consistent order (first bank by name, then additional bank)
	var originalBank, additionalBank *models.BankResponse
	for i := range banks {
		if banks[i].BankName == testBankName {
			originalBank = &banks[i]
		} else if banks[i].BankName == "Additional Bank" {
			additionalBank = &banks[i]
		}
	}
	
	// Validate the original bank (created during player creation)
	if originalBank == nil {
		t.Fatalf("Expected to find original bank with name '%s'", testBankName)
	}
	
	if originalBank.BankName != testBankName {
		t.Errorf("Expected original bank name '%s', got '%s'", testBankName, originalBank.BankName)
	}
	
	if originalBank.ClaimedCapital != 1000 {
		t.Errorf("Expected original bank claimed capital 1000, got %d", originalBank.ClaimedCapital)
	}
	
	if originalBank.ActualCapital != 1000 {
		t.Errorf("Expected original bank actual capital 1000, got %d", originalBank.ActualCapital)
	}
	
	if originalBank.Id == "" {
		t.Errorf("Expected original bank to have a valid ID, got empty string")
	}
	
	if len(originalBank.Assets) != 1 {
		t.Errorf("Expected original bank to have 1 asset (cash), got %d", len(originalBank.Assets))
	} else {
		// Validate the cash asset
		cashAsset := originalBank.Assets[0]
		if cashAsset.AssetType != "Cash" {
			t.Errorf("Expected original bank first asset to be 'Cash', got '%s'", cashAsset.AssetType)
		}
		if cashAsset.Amount != 1000 {
			t.Errorf("Expected original bank cash amount 1000, got %d", cashAsset.Amount)
		}
		if cashAsset.AssetTypeId == "" {
			t.Errorf("Expected original bank cash asset to have a valid AssetTypeId, got empty string")
		}
	}
	
	// Validate the additional bank (created by test)
	if additionalBank == nil {
		t.Fatalf("Expected to find additional bank with name 'Additional Bank'")
	}
	
	if additionalBank.BankName != "Additional Bank" {
		t.Errorf("Expected additional bank name 'Additional Bank', got '%s'", additionalBank.BankName)
	}
	
	if additionalBank.ClaimedCapital != 100000 {
		t.Errorf("Expected additional bank claimed capital 100000, got %d", additionalBank.ClaimedCapital)
	}
	
	if additionalBank.ActualCapital != 100000 {
		t.Errorf("Expected additional bank actual capital 100000, got %d", additionalBank.ActualCapital)
	}
	
	if additionalBank.Id == "" {
		t.Errorf("Expected additional bank to have a valid ID, got empty string")
	}
	
	if len(additionalBank.Assets) != 1 {
		t.Errorf("Expected additional bank to have 1 asset (cash), got %d", len(additionalBank.Assets))
	} else {
		// Validate the cash asset
		cashAsset := additionalBank.Assets[0]
		if cashAsset.AssetType != "Cash" {
			t.Errorf("Expected additional bank first asset to be 'Cash', got '%s'", cashAsset.AssetType)
		}
		if cashAsset.Amount != 100000 {
			t.Errorf("Expected additional bank cash amount 100000, got %d", cashAsset.Amount)
		}
		if cashAsset.AssetTypeId == "" {
			t.Errorf("Expected additional bank cash asset to have a valid AssetTypeId, got empty string")
		}
	}
	
	// Verify the two banks have different IDs
	if originalBank.Id == additionalBank.Id {
		t.Errorf("Expected banks to have different IDs, both have '%s'", originalBank.Id)
	}
}
