package tests

import (
	"context"
	"fmt"
	"ponziworld/backend/models"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBankService_GetAllBanksByUsername(t *testing.T) {
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

	// Retrieve banks directly via service
	bankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get banks by username: %v", err)
	}

	// Should have exactly one bank
	if len(bankResponses) != 1 {
		t.Fatalf("Expected 1 bank, got %d", len(bankResponses))
	}

	bankResponse := bankResponses[0]

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
	if len(bankResponse.AvailableAssets) != 5 {
		t.Errorf("Expected 5 available assets, got %d", len(bankResponse.AvailableAssets))
	}
	
	// Find the cash asset
	var cashAsset *models.AvailableAssetResponse
	for _, asset := range bankResponse.AvailableAssets {
		if asset.AssetName == "Cash" {
			cashAsset = &asset
			break
		}
	}
	
	if cashAsset == nil {
		t.Error("Expected to find Cash asset in available assets")
	} else {
		if !cashAsset.IsInvestedOrPending {
			t.Error("Expected Cash asset to be invested (bank has cash)")
		}
		if cashAsset.AssetTypeId == "" {
			t.Error("Expected asset type ID to be present")
		}
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
	
	if len(originalBank.AvailableAssets) != 5 {
		t.Errorf("Expected original bank to have 5 available assets, got %d", len(originalBank.AvailableAssets))
	} else {
		// Find the cash asset
		var cashAsset *models.AvailableAssetResponse
		for _, asset := range originalBank.AvailableAssets {
			if asset.AssetName == "Cash" {
				cashAsset = &asset
				break
			}
		}
		
		if cashAsset == nil {
			t.Error("Expected to find Cash asset in available assets")
		} else {
			if !cashAsset.IsInvestedOrPending {
				t.Error("Expected Cash asset to be invested (bank has cash)")
			}
			if cashAsset.AssetTypeId == "" {
				t.Error("Expected original bank cash asset to have a valid AssetTypeId")
			}
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
	
	if len(additionalBank.AvailableAssets) != 5 {
		t.Errorf("Expected additional bank to have 5 available assets, got %d", len(additionalBank.AvailableAssets))
	} else {
		// Find the cash asset
		var cashAsset *models.AvailableAssetResponse
		for _, asset := range additionalBank.AvailableAssets {
			if asset.AssetName == "Cash" {
				cashAsset = &asset
				break
			}
		}
		
		if cashAsset == nil {
			t.Error("Expected to find Cash asset in available assets")
		} else {
			if !cashAsset.IsInvestedOrPending {
				t.Error("Expected Cash asset to be invested (bank has cash)")
			}
			if cashAsset.AssetTypeId == "" {
				t.Error("Expected additional bank cash asset to have a valid AssetTypeId")
			}
		}
	}
	
	// Verify the two banks have different IDs
	if originalBank.Id == additionalBank.Id {
		t.Errorf("Expected banks to have different IDs, both have '%s'", originalBank.Id)
	}
}

func TestBankService_IsInvestedOrPendingFlag(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("bank_invested_pending")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("investedpendingtest_%d", timestamp)
	testPassword := "testpass123"
	testBankName := fmt.Sprintf("Test Bank %d", timestamp)
	
	// Create test user and bank
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}

	// Get asset types for testing
	assetTypes, err := container.ServiceContainer.AssetType.GetAllAssetTypes(ctx)
	if err != nil {
		t.Fatalf("Failed to get asset types: %v", err)
	}

	// Find specific asset types
	var cashAssetType, stocksAssetType, bondsAssetType *models.AssetType
	for _, assetType := range assetTypes {
		switch assetType.Name {
		case "Cash":
			cashAssetType = &assetType
		case "Stocks":
			stocksAssetType = &assetType
		case "Bonds":
			bondsAssetType = &assetType
		}
	}

	if cashAssetType == nil || stocksAssetType == nil || bondsAssetType == nil {
		t.Fatalf("Failed to find required asset types")
	}

	// Test initial state: Cash should be invested (bank starts with cash), others should not
	initialBankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get initial bank response: %v", err)
	}
	if len(initialBankResponses) == 0 {
		t.Fatalf("Expected at least one bank for user")
	}
	initialBankResponse := initialBankResponses[0]

	assetMap := make(map[string]*models.AvailableAssetResponse)
	for _, asset := range initialBankResponse.AvailableAssets {
		assetMap[asset.AssetName] = &asset
	}

	// Cash should be invested (bank starts with cash)
	if cashAsset, exists := assetMap["Cash"]; !exists {
		t.Error("Cash asset should be present in available assets")
	} else if !cashAsset.IsInvestedOrPending {
		t.Error("Cash asset should be invested initially (bank starts with cash)")
	}

	// Stocks should not be invested or pending initially
	if stocksAsset, exists := assetMap["Stocks"]; !exists {
		t.Error("Stocks asset should be present in available assets")
	} else if stocksAsset.IsInvestedOrPending {
		t.Error("Stocks asset should not be invested or pending initially")
	}

	// Bonds should not be invested or pending initially
	if bondsAsset, exists := assetMap["Bonds"]; !exists {
		t.Error("Bonds asset should be present in available assets")
	} else if bondsAsset.IsInvestedOrPending {
		t.Error("Bonds asset should not be invested or pending initially")
	}

	// Create a pending transaction for Stocks
	bankId, err := primitive.ObjectIDFromHex(initialBankResponse.Id)
	if err != nil {
		t.Fatalf("Failed to convert bank ID to ObjectID: %v", err)
	}
	
	err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(
		ctx, 
		bankId, 
		stocksAssetType.Id, 
		100, 
		testUsername,
	)
	if err != nil {
		t.Fatalf("Failed to create pending transaction: %v", err)
	}

	// Test after pending transaction: Stocks should now be pending
	afterPendingResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get bank response after pending transaction: %v", err)
	}
	if len(afterPendingResponses) == 0 {
		t.Fatalf("Expected at least one bank for user")
	}
	afterPendingResponse := afterPendingResponses[0]

	afterPendingAssetMap := make(map[string]*models.AvailableAssetResponse)
	for _, asset := range afterPendingResponse.AvailableAssets {
		afterPendingAssetMap[asset.AssetName] = &asset
	}

	// Cash should still be invested
	if cashAsset, exists := afterPendingAssetMap["Cash"]; !exists {
		t.Error("Cash asset should be present in available assets after pending transaction")
	} else if !cashAsset.IsInvestedOrPending {
		t.Error("Cash asset should still be invested after pending transaction")
	}

	// Stocks should now be pending
	if stocksAsset, exists := afterPendingAssetMap["Stocks"]; !exists {
		t.Error("Stocks asset should be present in available assets after pending transaction")
	} else if !stocksAsset.IsInvestedOrPending {
		t.Error("Stocks asset should be pending after creating pending transaction")
	}

	// Bonds should still not be invested or pending
	if bondsAsset, exists := afterPendingAssetMap["Bonds"]; !exists {
		t.Error("Bonds asset should be present in available assets after pending transaction")
	} else if bondsAsset.IsInvestedOrPending {
		t.Error("Bonds asset should still not be invested or pending after pending transaction")
	}

	// Create an actual investment in Bonds by directly adding an asset
	bondsAsset := &models.Investment{
		Id:          primitive.NewObjectID(),
		SourceBankId:      bankId,
		Amount:      500,
		TargetAssetId: bondsAssetType.Id,
	}
	err = container.RepositoryContainer.Investment.Create(ctx, bondsAsset)
	if err != nil {
		t.Fatalf("Failed to create bonds asset: %v", err)
	}

	// Test after actual investment: Bonds should now be invested
	afterInvestmentResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get bank response after investment: %v", err)
	}
	if len(afterInvestmentResponses) == 0 {
		t.Fatalf("Expected at least one bank for user")
	}
	afterInvestmentResponse := afterInvestmentResponses[0]

	afterInvestmentAssetMap := make(map[string]*models.AvailableAssetResponse)
	for _, asset := range afterInvestmentResponse.AvailableAssets {
		afterInvestmentAssetMap[asset.AssetName] = &asset
	}

	// Cash should still be invested
	if cashAsset, exists := afterInvestmentAssetMap["Cash"]; !exists {
		t.Error("Cash asset should be present in available assets after investment")
	} else if !cashAsset.IsInvestedOrPending {
		t.Error("Cash asset should still be invested after investment")
	}

	// Stocks should still be pending
	if stocksAsset, exists := afterInvestmentAssetMap["Stocks"]; !exists {
		t.Error("Stocks asset should be present in available assets after investment")
	} else if !stocksAsset.IsInvestedOrPending {
		t.Error("Stocks asset should still be pending after investment")
	}

	// Bonds should now be invested
	if bondsAsset, exists := afterInvestmentAssetMap["Bonds"]; !exists {
		t.Error("Bonds asset should be present in available assets after investment")
	} else if !bondsAsset.IsInvestedOrPending {
		t.Error("Bonds asset should be invested after creating investment")
	}
}
