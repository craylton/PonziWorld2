package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"ponziworld/backend/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAssetService_GetAssetDetails_ValidScenarios(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("asset_details_valid")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("assettest_%d", timestamp)
	testBankName := "Test Bank Asset Details"
	testPassword := "testpassword123"

	// Create player and bank
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create new player: %v", err)
	}

	// Get the bank
	bankResponse, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get bank: %v", err)
	}

	bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
	if err != nil {
		t.Fatalf("Failed to parse bank ID: %v", err)
	}

	// Get asset types
	assetTypes, err := container.ServiceContainer.AssetType.GetAllAssetTypes(ctx)
	if err != nil {
		t.Fatalf("Failed to get asset types: %v", err)
	}

	// Find asset types for testing
	var cashAssetType, stocksAssetType, bondsAssetType *primitive.ObjectID
	for _, assetType := range assetTypes {
		switch assetType.Name {
		case "Cash":
			cashAssetType = &assetType.Id
		case "Stocks":
			stocksAssetType = &assetType.Id
		case "Bonds":
			bondsAssetType = &assetType.Id
		}
	}

	if cashAssetType == nil || stocksAssetType == nil || bondsAssetType == nil {
		t.Fatalf("Required asset types not found")
	}

	t.Run("Cash asset with initial investment", func(t *testing.T) {
		// Test getting asset details for Cash (should have 1000 invested, 0 pending)
		assetDetails, err := container.ServiceContainer.Asset.GetAssetDetails(ctx, testUsername, *cashAssetType, bankID)
		if err != nil {
			t.Fatalf("Failed to get asset details: %v", err)
		}

		// Verify the response
		if assetDetails.InvestedAmount != 1000 {
			t.Errorf("Expected invested amount 1000, got %d", assetDetails.InvestedAmount)
		}
		if assetDetails.PendingAmount != 0 {
			t.Errorf("Expected pending amount 0, got %d", assetDetails.PendingAmount)
		}
		if len(assetDetails.HistoricalData) == 0 {
			t.Errorf("Expected historical data, got empty array")
		}
		
		// Verify historical data has 8 days
		if len(assetDetails.HistoricalData) != 8 {
			t.Errorf("Expected 8 days of historical data, got %d", len(assetDetails.HistoricalData))
		}
	})

	t.Run("Asset with no investment", func(t *testing.T) {
		// Test getting asset details for Stocks (should have 0 invested, 0 pending)
		assetDetails, err := container.ServiceContainer.Asset.GetAssetDetails(ctx, testUsername, *stocksAssetType, bankID)
		if err != nil {
			t.Fatalf("Failed to get asset details for Stocks: %v", err)
		}

		// Verify the response
		if assetDetails.InvestedAmount != 0 {
			t.Errorf("Expected invested amount 0, got %d", assetDetails.InvestedAmount)
		}
		if assetDetails.PendingAmount != 0 {
			t.Errorf("Expected pending amount 0, got %d", assetDetails.PendingAmount)
		}
		if len(assetDetails.HistoricalData) != 8 {
			t.Errorf("Expected 8 days of historical data, got %d", len(assetDetails.HistoricalData))
		}
		
		// Verify historical data contains default values
		for _, data := range assetDetails.HistoricalData {
			if data.Value != 1000 { // DefaultPerformanceValue
				t.Errorf("Expected default performance value 1000, got %d", data.Value)
			}
		}
	})

	//todo: This should neverr happen - revisit this and figure out what to do with it
	t.Run("Asset with multiple pending transactions", func(t *testing.T) {
		// Create multiple pending buy transactions
		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(ctx, bankID, *bondsAssetType, 200, testUsername)
		if err != nil {
			t.Fatalf("Failed to create first buy transaction: %v", err)
		}

		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(ctx, bankID, *bondsAssetType, 300, testUsername)
		if err != nil {
			t.Fatalf("Failed to create second buy transaction: %v", err)
		}

		// Create a sell transaction
		err = container.ServiceContainer.PendingTransaction.CreateSellTransaction(ctx, bankID, *bondsAssetType, 100, testUsername)
		if err != nil {
			t.Fatalf("Failed to create sell transaction: %v", err)
		}

		// Test getting asset details for Bonds (should have 0 invested, 400 pending = 500 buy - 100 sell)
		assetDetails, err := container.ServiceContainer.Asset.GetAssetDetails(ctx, testUsername, *bondsAssetType, bankID)
		if err != nil {
			t.Fatalf("Failed to get asset details for Bonds: %v", err)
		}

		// Verify the response
		if assetDetails.InvestedAmount != 0 {
			t.Errorf("Expected invested amount 0, got %d", assetDetails.InvestedAmount)
		}
		if assetDetails.PendingAmount != 400 {
			t.Errorf("Expected pending amount 400 (500 buy - 100 sell), got %d", assetDetails.PendingAmount)
		}
		if len(assetDetails.HistoricalData) != 8 {
			t.Errorf("Expected 8 days of historical data, got %d", len(assetDetails.HistoricalData))
		}
	})
}

func TestAssetService_GetAssetDetails_ErrorCases(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("asset_details_errors")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("assettest_%d", timestamp)
	testBankName := "Test Bank Asset Details Errors"
	testPassword := "testpassword123"

	// Create player and bank
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create new player: %v", err)
	}

	// Get the bank
	bankResponse, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get bank: %v", err)
	}

	bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
	if err != nil {
		t.Fatalf("Failed to parse bank ID: %v", err)
	}

	// Get asset types
	assetTypes, err := container.ServiceContainer.AssetType.GetAllAssetTypes(ctx)
	if err != nil {
		t.Fatalf("Failed to get asset types: %v", err)
	}

	var cashAssetType *primitive.ObjectID
	for _, assetType := range assetTypes {
		if assetType.Name == "Cash" {
			cashAssetType = &assetType.Id
			break
		}
	}

	if cashAssetType == nil {
		t.Fatalf("Cash asset type not found")
	}

	t.Run("Invalid asset ID", func(t *testing.T) {
		// Test with invalid asset ID
		invalidAssetID := primitive.NewObjectID()
		_, err = container.ServiceContainer.Asset.GetAssetDetails(ctx, testUsername, invalidAssetID, bankID)
		if err == nil {
			t.Error("Expected error for invalid asset ID, got nil")
		}
		if err != services.ErrAssetNotFound {
			t.Errorf("Expected ErrAssetNotFound, got %v", err)
		}
	})

	t.Run("Invalid bank ID", func(t *testing.T) {
		// Test with invalid bank ID
		invalidBankID := primitive.NewObjectID()
		_, err = container.ServiceContainer.Asset.GetAssetDetails(ctx, testUsername, *cashAssetType, invalidBankID)
		if err == nil {
			t.Error("Expected error for invalid bank ID, got nil")
		}
		if err != services.ErrBankNotFound {
			t.Errorf("Expected ErrBankNotFound, got %v", err)
		}
	})

	t.Run("Non-existent user", func(t *testing.T) {
		// Test with non-existent user
		_, err = container.ServiceContainer.Asset.GetAssetDetails(ctx, "nonexistentuser", *cashAssetType, bankID)
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
		if err != services.ErrPlayerNotFound {
			t.Errorf("Expected ErrPlayerNotFound, got %v", err)
		}
	})

	t.Run("Unauthorized bank access", func(t *testing.T) {
		// Create another user
		otherUsername := fmt.Sprintf("otheruser_%d", timestamp)
		err = container.ServiceContainer.Player.CreateNewPlayer(ctx, otherUsername, testPassword, "Other Bank")
		if err != nil {
			t.Fatalf("Failed to create other player: %v", err)
		}

		// Try to access first bank's asset details with second user's credentials
		_, err = container.ServiceContainer.Asset.GetAssetDetails(ctx, otherUsername, *cashAssetType, bankID)
		if err == nil {
			t.Error("Expected error for unauthorized access, got nil")
		}
		if err != services.ErrUnauthorized {
			t.Errorf("Expected ErrUnauthorized, got %v", err)
		}
	})
}

func TestAssetService_GetAssetDetails_HistoricalDataGeneration(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("asset_details_historical")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("assettest_%d", timestamp)
	testBankName := "Test Bank Historical Data"
	testPassword := "testpassword123"

	// Create player and bank
	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
	if err != nil {
		t.Fatalf("Failed to create new player: %v", err)
	}

	// Get the bank
	bankResponse, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, testUsername)
	if err != nil {
		t.Fatalf("Failed to get bank: %v", err)
	}

	bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
	if err != nil {
		t.Fatalf("Failed to parse bank ID: %v", err)
	}

	// Get asset types
	assetTypes, err := container.ServiceContainer.AssetType.GetAllAssetTypes(ctx)
	if err != nil {
		t.Fatalf("Failed to get asset types: %v", err)
	}

	var stocksAssetType *primitive.ObjectID
	for _, assetType := range assetTypes {
		if assetType.Name == "Stocks" {
			stocksAssetType = &assetType.Id
			break
		}
	}

	if stocksAssetType == nil {
		t.Fatalf("Stocks asset type not found")
	}

	t.Run("Historical data generation for new asset", func(t *testing.T) {
		// Get asset details for an asset with no existing historical data
		assetDetails, err := container.ServiceContainer.Asset.GetAssetDetails(ctx, testUsername, *stocksAssetType, bankID)
		if err != nil {
			t.Fatalf("Failed to get asset details: %v", err)
		}

		// Verify historical data was generated
		if len(assetDetails.HistoricalData) != 8 {
			t.Errorf("Expected 8 days of historical data, got %d", len(assetDetails.HistoricalData))
		}

		// Verify all historical data has default values
		for i, data := range assetDetails.HistoricalData {
			if data.Value != 1000 {
				t.Errorf("Historical data[%d]: expected value 1000, got %d", i, data.Value)
			}
			// Days can be negative (representing days before current day) or zero (current day)
			// Just verify we have actual day values, not checking specific values
		}

		// Verify days are in ascending order
		for i := 1; i < len(assetDetails.HistoricalData); i++ {
			if assetDetails.HistoricalData[i].Day <= assetDetails.HistoricalData[i-1].Day {
				t.Errorf("Historical data days not in ascending order: %d <= %d", 
					assetDetails.HistoricalData[i].Day, assetDetails.HistoricalData[i-1].Day)
			}
		}
	})

	t.Run("Historical data persistence", func(t *testing.T) {
		// Call the endpoint twice to ensure data is persisted
		assetDetails1, err := container.ServiceContainer.Asset.GetAssetDetails(ctx, testUsername, *stocksAssetType, bankID)
		if err != nil {
			t.Fatalf("Failed to get asset details first time: %v", err)
		}

		assetDetails2, err := container.ServiceContainer.Asset.GetAssetDetails(ctx, testUsername, *stocksAssetType, bankID)
		if err != nil {
			t.Fatalf("Failed to get asset details second time: %v", err)
		}

		// Verify both calls return the same historical data
		if len(assetDetails1.HistoricalData) != len(assetDetails2.HistoricalData) {
			t.Errorf("Historical data length mismatch: %d vs %d", 
				len(assetDetails1.HistoricalData), len(assetDetails2.HistoricalData))
		}

		for i := 0; i < len(assetDetails1.HistoricalData); i++ {
			if assetDetails1.HistoricalData[i].Day != assetDetails2.HistoricalData[i].Day {
				t.Errorf("Historical data day mismatch at index %d: %d vs %d", 
					i, assetDetails1.HistoricalData[i].Day, assetDetails2.HistoricalData[i].Day)
			}
			if assetDetails1.HistoricalData[i].Value != assetDetails2.HistoricalData[i].Value {
				t.Errorf("Historical data value mismatch at index %d: %d vs %d", 
					i, assetDetails1.HistoricalData[i].Value, assetDetails2.HistoricalData[i].Value)
			}
		}
	})
}

func TestAssetService_GetAssetDetails_BankAsAsset(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("asset_details_bank_as_asset")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	timestamp := time.Now().Unix()
	
	// Create two players and their banks
	investor1Username := fmt.Sprintf("investor1_%d", timestamp)
	investor2Username := fmt.Sprintf("investor2_%d", timestamp)
	testPassword := "testpassword123"

	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, investor1Username, testPassword, "Investor 1 Bank")
	if err != nil {
		t.Fatalf("Failed to create investor 1: %v", err)
	}

	err = container.ServiceContainer.Player.CreateNewPlayer(ctx, investor2Username, testPassword, "Investor 2 Bank")
	if err != nil {
		t.Fatalf("Failed to create investor 2: %v", err)
	}

	// Get the banks
	bank1Response, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, investor1Username)
	if err != nil {
		t.Fatalf("Failed to get bank 1: %v", err)
	}

	bank2Response, err := container.ServiceContainer.Bank.GetBankByUsername(ctx, investor2Username)
	if err != nil {
		t.Fatalf("Failed to get bank 2: %v", err)
	}

	bank1ID, err := primitive.ObjectIDFromHex(bank1Response.Id)
	if err != nil {
		t.Fatalf("Failed to parse bank 1 ID: %v", err)
	}

	bank2ID, err := primitive.ObjectIDFromHex(bank2Response.Id)
	if err != nil {
		t.Fatalf("Failed to parse bank 2 ID: %v", err)
	}

	t.Run("Bank as asset - no investment", func(t *testing.T) {
		// Test investor 1 getting details about bank 2 as an asset (no investment yet)
		assetDetails, err := container.ServiceContainer.Asset.GetAssetDetails(ctx, investor1Username, bank2ID, bank1ID)
		if err != nil {
			t.Fatalf("Failed to get asset details for bank as asset: %v", err)
		}

		// Verify the response
		if assetDetails.InvestedAmount != 0 {
			t.Errorf("Expected invested amount 0, got %d", assetDetails.InvestedAmount)
		}
		if assetDetails.PendingAmount != 0 {
			t.Errorf("Expected pending amount 0, got %d", assetDetails.PendingAmount)
		}
		if len(assetDetails.HistoricalData) != 8 {
			t.Errorf("Expected 8 days of historical data, got %d", len(assetDetails.HistoricalData))
		}
	})

	t.Run("Bank as asset - with pending investment", func(t *testing.T) {
		// Create pending transaction for investor 1 to invest in bank 2
		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(ctx, bank1ID, bank2ID, 500, investor1Username)
		if err != nil {
			t.Fatalf("Failed to create buy transaction for bank investment: %v", err)
		}

		// Test investor 1 getting details about bank 2 as an asset (with pending investment)
		assetDetails, err := container.ServiceContainer.Asset.GetAssetDetails(ctx, investor1Username, bank2ID, bank1ID)
		if err != nil {
			t.Fatalf("Failed to get asset details for bank as asset with pending: %v", err)
		}

		// Verify the response
		if assetDetails.InvestedAmount != 0 {
			t.Errorf("Expected invested amount 0, got %d", assetDetails.InvestedAmount)
		}
		if assetDetails.PendingAmount != 500 {
			t.Errorf("Expected pending amount 500, got %d", assetDetails.PendingAmount)
		}
		if len(assetDetails.HistoricalData) != 8 {
			t.Errorf("Expected 8 days of historical data, got %d", len(assetDetails.HistoricalData))
		}
		
		// Verify historical data contains bank 2's performance data
		// Historical data should have default values for bank 2
		for i, data := range assetDetails.HistoricalData {
			if data.Value != 1000 {
				t.Errorf("Historical data[%d]: expected value 1000, got %d", i, data.Value)
			}
		}
	})
}
