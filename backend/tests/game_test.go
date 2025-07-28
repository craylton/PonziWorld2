package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/models"
	"ponziworld/backend/routes"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGameService_NextDay(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("game")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()

	t.Run("should create initial day 0 and increment to day 1", func(t *testing.T) {
		day, err := container.ServiceContainer.Game.NextDay(ctx)
		if err != nil {
			t.Fatalf("Failed to advance to next day: %v", err)
		}

		if day != 1 {
			t.Errorf("Expected day to be 1, got %d", day)
		}
	})

	t.Run("should increment existing day", func(t *testing.T) {
		day, err := container.ServiceContainer.Game.NextDay(ctx)
		if err != nil {
			t.Fatalf("Failed to advance to next day: %v", err)
		}

		if day != 2 {
			t.Errorf("Expected day to be 2, got %d", day)
		}
	})
}

func TestNextDayEndpoint(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("game")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, container)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should reject non-admin users", func(t *testing.T) {
		// Create a regular (non-admin) user with unique username
		timestamp := time.Now().Unix()
		regularUsername := fmt.Sprintf("regularuser_%d", timestamp)
		regularToken, err := CreateRegularUserForTest(container, regularUsername, "password123", "RegularBank")
		if err != nil {
			t.Fatal("Failed to create regular user:", err)
		}

		req, err := http.NewRequest("POST", server.URL+"/api/nextDay", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+regularToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
		}
	})
}

func TestGameService_CurrentDay(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("game")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()

	t.Run("should return day 0 when no game state exists", func(t *testing.T) {
		currentDay, err := container.ServiceContainer.Game.GetCurrentDay(ctx)
		if err != nil {
			t.Fatalf("Failed to get current day: %v", err)
		}

		if currentDay != 0 {
			t.Errorf("Expected currentDay to be 0, got %d", currentDay)
		}
	})

	t.Run("should return current day when game state exists", func(t *testing.T) {
		// Advance the game to day 5 by calling NextDay service directly
		var finalDay int
		for i := range 5 {
			day, err := container.ServiceContainer.Game.NextDay(ctx)
			if err != nil {
				t.Fatalf("Failed to advance to day %d: %v", i+1, err)
			}
			finalDay = day
		}

		if finalDay != 5 {
			t.Errorf("Expected final day to be 5, got %d", finalDay)
		}

		currentDay, err := container.ServiceContainer.Game.GetCurrentDay(ctx)
		if err != nil {
			t.Fatalf("Failed to get current day: %v", err)
		}

		if currentDay != 5 {
			t.Errorf("Expected currentDay to be 5, got %d", currentDay)
		}
	})
}

func TestGameService_ProcessPendingTransactions(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("game")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()

	t.Run("should process pending buy transaction and create new investment", func(t *testing.T) {
		// Setup: Create user and bank
		timestamp := time.Now().Unix()
		testUsername := fmt.Sprintf("testuser_%d", timestamp)
		testBankName := "TestBank"
		testPassword := "testpassword123"

		// Create player and bank
		err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
		if err != nil {
			t.Fatal("Failed to create test user:", err)
		}

		// Get bank
		bankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
		if err != nil {
			t.Fatal("Failed to get banks:", err)
		}
		bankResponse := bankResponses[0]

		bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
		if err != nil {
			t.Fatal("Failed to convert bank ID:", err)
		}

		// Create a pending buy transaction for an asset type
		assetTypes, err := container.RepositoryContainer.AssetType.FindAll(ctx)
		if err != nil {
			t.Fatal("Failed to get asset types:", err)
		}
		
		var targetAssetID primitive.ObjectID
		for _, assetType := range assetTypes {
			if assetType.Name != "Cash" {
				targetAssetID = assetType.Id
				break
			}
		}

		// Create buy transaction through the service
		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(
			ctx,
			bankID,
			targetAssetID,
			100,
			testUsername,
		)
		if err != nil {
			t.Fatal("Failed to create buy transaction:", err)
		}

		// Process pending transactions
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process pending transactions:", err)
		}

		// Verify the investment was created
		investment, err := container.RepositoryContainer.Investment.FindBySourceIdAndTargetId(
			ctx,
			bankID,
			targetAssetID,
		)
		if err != nil {
			t.Fatal("Investment should have been created:", err)
		}

		if investment.Amount != 100 {
			t.Errorf("Expected investment amount to be 100, got %d", investment.Amount)
		}

		// Verify pending transactions were removed
		pendingTransactions, err := container.RepositoryContainer.PendingTransaction.FindBySourceBankID(ctx, bankID)
		if err != nil {
			t.Fatal("Failed to get pending transactions:", err)
		}

		if len(pendingTransactions) != 0 {
			t.Errorf("Expected no pending transactions after processing, got %d", len(pendingTransactions))
		}
	})

	t.Run("should combine pending transaction with existing investment", func(t *testing.T) {
		// Setup: Create user and bank
		timestamp := time.Now().Unix()
		testUsername := fmt.Sprintf("testuser2_%d", timestamp)
		testBankName := "TestBank2"
		testPassword := "testpassword123"

		// Create player and bank
		err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
		if err != nil {
			t.Fatal("Failed to create test user:", err)
		}

		// Get bank
		bankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
		if err != nil {
			t.Fatal("Failed to get banks:", err)
		}
		bankResponse := bankResponses[0]

		bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
		if err != nil {
			t.Fatal("Failed to convert bank ID:", err)
		}

		// Get an asset type for testing
		assetTypes, err := container.RepositoryContainer.AssetType.FindAll(ctx)
		if err != nil {
			t.Fatal("Failed to get asset types:", err)
		}
		
		var targetAssetID primitive.ObjectID
		for _, assetType := range assetTypes {
			if assetType.Name != "Cash" {
				targetAssetID = assetType.Id
				break
			}
		}

		// Create initial investment
		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(
			ctx,
			bankID,
			targetAssetID,
			50,
			testUsername,
		)
		if err != nil {
			t.Fatal("Failed to create initial buy transaction:", err)
		}

		// Process to create initial investment
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process first transaction:", err)
		}

		// Create another buy transaction
		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(
			ctx,
			bankID,
			targetAssetID,
			30,
			testUsername,
		)
		if err != nil {
			t.Fatal("Failed to create second buy transaction:", err)
		}

		// Process pending transactions again
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process second transaction:", err)
		}

		// Verify the investment was updated (50 + 30 = 80)
		investment, err := container.RepositoryContainer.Investment.FindBySourceIdAndTargetId(
			ctx,
			bankID,
			targetAssetID,
		)
		if err != nil {
			t.Fatal("Investment should exist:", err)
		}

		if investment.Amount != 80 {
			t.Errorf("Expected investment amount to be 80, got %d", investment.Amount)
		}
	})

	t.Run("should remove investment when sell amount equals existing amount", func(t *testing.T) {
		// Setup: Create user and bank
		timestamp := time.Now().Unix()
		testUsername := fmt.Sprintf("testuser3_%d", timestamp)
		testBankName := "TestBank3"
		testPassword := "testpassword123"

		// Create player and bank
		err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
		if err != nil {
			t.Fatal("Failed to create test user:", err)
		}

		// Get bank
		bankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
		if err != nil {
			t.Fatal("Failed to get banks:", err)
		}
		bankResponse := bankResponses[0]

		bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
		if err != nil {
			t.Fatal("Failed to convert bank ID:", err)
		}

		// Get an asset type for testing
		assetTypes, err := container.RepositoryContainer.AssetType.FindAll(ctx)
		if err != nil {
			t.Fatal("Failed to get asset types:", err)
		}
		
		var targetAssetID primitive.ObjectID
		for _, assetType := range assetTypes {
			if assetType.Name != "Cash" {
				targetAssetID = assetType.Id
				break
			}
		}

		// Create initial investment
		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(
			ctx,
			bankID,
			targetAssetID,
			50,
			testUsername,
		)
		if err != nil {
			t.Fatal("Failed to create buy transaction:", err)
		}

		// Process to create initial investment
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process buy transaction:", err)
		}

		// Create sell transaction for the same amount
		err = container.ServiceContainer.PendingTransaction.CreateSellTransaction(
			ctx,
			bankID,
			targetAssetID,
			50,
			testUsername,
		)
		if err != nil {
			t.Fatal("Failed to create sell transaction:", err)
		}

		// Process pending transactions
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process sell transaction:", err)
		}

		// Verify the investment was removed
		_, err = container.RepositoryContainer.Investment.FindBySourceIdAndTargetId(
			ctx,
			bankID,
			targetAssetID,
		)
		if err == nil {
			t.Error("Investment should have been removed after selling all shares")
		}
	})

	t.Run("should reduce investment amount when selling partial amount", func(t *testing.T) {
		// Setup: Create user and bank
		timestamp := time.Now().Unix()
		testUsername := fmt.Sprintf("testuser4_%d", timestamp)
		testBankName := "TestBank4"
		testPassword := "testpassword123"

		// Create player and bank
		err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
		if err != nil {
			t.Fatal("Failed to create test user:", err)
		}

		// Get bank
		bankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
		if err != nil {
			t.Fatal("Failed to get banks:", err)
		}
		bankResponse := bankResponses[0]

		bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
		if err != nil {
			t.Fatal("Failed to convert bank ID:", err)
		}

		// Get an asset type for testing
		assetTypes, err := container.RepositoryContainer.AssetType.FindAll(ctx)
		if err != nil {
			t.Fatal("Failed to get asset types:", err)
		}
		
		var targetAssetID primitive.ObjectID
		for _, assetType := range assetTypes {
			if assetType.Name != "Cash" {
				targetAssetID = assetType.Id
				break
			}
		}

		// Create initial investment
		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(
			ctx,
			bankID,
			targetAssetID,
			100,
			testUsername,
		)
		if err != nil {
			t.Fatal("Failed to create buy transaction:", err)
		}

		// Process to create initial investment
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process buy transaction:", err)
		}

		// Create sell transaction for partial amount
		err = container.ServiceContainer.PendingTransaction.CreateSellTransaction(
			ctx,
			bankID,
			targetAssetID,
			30,
			testUsername,
		)
		if err != nil {
			t.Fatal("Failed to create sell transaction:", err)
		}

		// Process pending transactions
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process sell transaction:", err)
		}

		// Verify the investment was reduced (100 - 30 = 70)
		investment, err := container.RepositoryContainer.Investment.FindBySourceIdAndTargetId(
			ctx,
			bankID,
			targetAssetID,
		)
		if err != nil {
			t.Fatal("Investment should still exist:", err)
		}

		if investment.Amount != 70 {
			t.Errorf("Expected investment amount to be 70, got %d", investment.Amount)
		}
	})

	t.Run("should process multiple transactions for different assets", func(t *testing.T) {
		// Setup: Create user and bank
		timestamp := time.Now().Unix()
		testUsername := fmt.Sprintf("testuser5_%d", timestamp)
		testBankName := "TestBank5"
		testPassword := "testpassword123"

		// Create player and bank
		err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
		if err != nil {
			t.Fatal("Failed to create test user:", err)
		}

		// Get bank
		bankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
		if err != nil {
			t.Fatal("Failed to get banks:", err)
		}
		bankResponse := bankResponses[0]

		bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
		if err != nil {
			t.Fatal("Failed to convert bank ID:", err)
		}

		// Get multiple asset types for testing
		assetTypes, err := container.RepositoryContainer.AssetType.FindAll(ctx)
		if err != nil {
			t.Fatal("Failed to get asset types:", err)
		}
		
		var targetAssetIDs []primitive.ObjectID
		for _, assetType := range assetTypes {
			if assetType.Name != "Cash" {
				targetAssetIDs = append(targetAssetIDs, assetType.Id)
				if len(targetAssetIDs) >= 2 { // Get at least 2 non-cash asset types
					break
				}
			}
		}

		// Create multiple buy transactions
		for i, assetID := range targetAssetIDs {
			err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(
				ctx,
				bankID,
				assetID,
				int64((i+1)*25), // 25, 50, etc.
				testUsername,
			)
			if err != nil {
				t.Fatalf("Failed to create buy transaction %d: %v", i, err)
			}
		}

		// Process pending transactions
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process transactions:", err)
		}

		// Verify all investments were created correctly
		for i, assetID := range targetAssetIDs {
			investment, err := container.RepositoryContainer.Investment.FindBySourceIdAndTargetId(
				ctx,
				bankID,
				assetID,
			)
			if err != nil {
				t.Fatalf("Investment %d should have been created: %v", i, err)
			}

			expectedAmount := int64((i + 1) * 25)
			if investment.Amount != expectedAmount {
				t.Errorf("Expected investment %d amount to be %d, got %d", i, expectedAmount, investment.Amount)
			}
		}
	})
}

func TestNextDay_WithPendingTransactions(t *testing.T) {
	// Create test dependencies
	container, err := CreateTestDependencies("game")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()

	t.Run("should process pending transactions when advancing day", func(t *testing.T) {
		// Setup: Create user and bank
		timestamp := time.Now().Unix()
		testUsername := fmt.Sprintf("daytest_%d", timestamp)
		testBankName := "DayTestBank"
		testPassword := "testpassword123"

		// Create player and bank
		err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
		if err != nil {
			t.Fatal("Failed to create test user:", err)
		}

		// Get bank
		bankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
		if err != nil {
			t.Fatal("Failed to get banks:", err)
		}
		bankResponse := bankResponses[0]

		bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
		if err != nil {
			t.Fatal("Failed to convert bank ID:", err)
		}

		// Get an asset type for testing
		assetTypes, err := container.RepositoryContainer.AssetType.FindAll(ctx)
		if err != nil {
			t.Fatal("Failed to get asset types:", err)
		}
		
		var targetAssetID primitive.ObjectID
		for _, assetType := range assetTypes {
			if assetType.Name != "Cash" {
				targetAssetID = assetType.Id
				break
			}
		}

		// Create a pending buy transaction
		err = container.ServiceContainer.PendingTransaction.CreateBuyTransaction(
			ctx,
			bankID,
			targetAssetID,
			75,
			testUsername,
		)
		if err != nil {
			t.Fatal("Failed to create buy transaction:", err)
		}

		// Verify pending transaction exists
		pendingTransactions, err := container.RepositoryContainer.PendingTransaction.FindBySourceBankID(ctx, bankID)
		if err != nil {
			t.Fatal("Failed to get pending transactions:", err)
		}

		// Should have 2 pending transactions (asset + cash)
		if len(pendingTransactions) != 2 {
			t.Errorf("Expected 2 pending transactions before advancing day, got %d", len(pendingTransactions))
		}

		// Advance to next day (this should process pending transactions)
		day, err := container.ServiceContainer.Game.NextDay(ctx)
		if err != nil {
			t.Fatal("Failed to advance to next day:", err)
		}

		if day != 1 {
			t.Errorf("Expected day to be 1, got %d", day)
		}

		// Verify pending transactions were processed and removed
		pendingTransactions, err = container.RepositoryContainer.PendingTransaction.FindBySourceBankID(ctx, bankID)
		if err != nil {
			t.Fatal("Failed to get pending transactions after day advance:", err)
		}

		if len(pendingTransactions) != 0 {
			t.Errorf("Expected no pending transactions after advancing day, got %d", len(pendingTransactions))
		}

		// Verify investment was created
		investment, err := container.RepositoryContainer.Investment.FindBySourceIdAndTargetId(
			ctx,
			bankID,
			targetAssetID,
		)
		if err != nil {
			t.Fatal("Investment should have been created:", err)
		}

		if investment.Amount != 75 {
			t.Errorf("Expected investment amount to be 75, got %d", investment.Amount)
		}
	})

	t.Run("should allow cash transactions to result in negative balances", func(t *testing.T) {
		// Setup: Create user and bank
		timestamp := time.Now().Unix()
		testUsername := fmt.Sprintf("cashtest_%d", timestamp)
		testBankName := "CashTestBank"
		testPassword := "testpassword123"

		// Create player and bank
		err = container.ServiceContainer.Player.CreateNewPlayer(ctx, testUsername, testPassword, testBankName)
		if err != nil {
			t.Fatal("Failed to create test user:", err)
		}

		// Get bank
		bankResponses, err := container.ServiceContainer.Bank.GetAllBanksByUsername(ctx, testUsername)
		if err != nil {
			t.Fatal("Failed to get banks:", err)
		}
		bankResponse := bankResponses[0]

		bankID, err := primitive.ObjectIDFromHex(bankResponse.Id)
		if err != nil {
			t.Fatal("Failed to convert bank ID:", err)
		}

		// Get the cash asset type
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatal("Failed to get cash asset type:", err)
		}

		// Create a manual pending transaction that would result in negative cash
		// This simulates a scenario where pending transactions were created but conditions changed
		pendingTransaction := &models.PendingTransactionResponse{
			Id:            primitive.NewObjectID(),
			SourceBankId:  bankID,
			TargetAssetId: cashAssetType.Id,
			Amount:        -2000, // More than the initial 1000 cash
		}
		
		err = container.RepositoryContainer.PendingTransaction.Create(ctx, pendingTransaction)
		if err != nil {
			t.Fatal("Failed to create pending transaction:", err)
		}

		// Process pending transactions - should allow negative cash
		err = container.ServiceContainer.Game.ProcessPendingTransactions(ctx)
		if err != nil {
			t.Fatal("Failed to process pending transactions:", err)
		}

		// Verify the cash investment is negative
		cashInvestment, err := container.RepositoryContainer.Investment.FindBySourceIdAndTargetId(
			ctx,
			bankID,
			cashAssetType.Id,
		)
		if err != nil {
			t.Fatal("Cash investment should exist:", err)
		}

		// Initial cash was 1000, we spent 2000, so should be -1000
		expectedAmount := int64(-1000)
		if cashInvestment.Amount != expectedAmount {
			t.Errorf("Expected cash amount to be %d, got %d", expectedAmount, cashInvestment.Amount)
		}

		// Verify pending transaction was removed
		pendingTransactions, err := container.RepositoryContainer.PendingTransaction.FindBySourceBankID(ctx, bankID)
		if err != nil {
			t.Fatal("Failed to get pending transactions:", err)
		}

		if len(pendingTransactions) != 0 {
			t.Errorf("Expected no pending transactions after processing, got %d", len(pendingTransactions))
		}
	})
}
