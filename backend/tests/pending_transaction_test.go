package tests

import (
	"context"
	"fmt"
	"ponziworld/backend/models"
	"ponziworld/backend/services"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPendingTransactionService_CreateTransactions(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create test user and bank
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "testpass"
	bankName := "Test Bank"

	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	// Get the banks for the user
	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user.Id)
	if err != nil { 
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("Valid buy transaction creation", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 1000, username)
		if err != nil {
			t.Errorf("Expected no error for valid buy transaction, got: %v", err)
		}

		// Verify transactions were created (should be 2: asset purchase + cash deduction)
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions (asset + cash), got %d", len(transactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Find the asset transaction and cash transaction
		var assetTransaction, cashTransaction *models.PendingTransactionResponse
		for i := range transactions {
			if transactions[i].TargetAssetId == assetType.Id {
				assetTransaction = &transactions[i]
			} else if transactions[i].TargetAssetId == cashAssetType.Id {
				cashTransaction = &transactions[i]
			}
		}

		// Verify asset transaction
		if assetTransaction == nil {
			t.Errorf("Expected to find asset transaction for asset ID %s", assetType.Id.Hex())
		} else {
			if assetTransaction.Amount != 1000 {
				t.Errorf("Expected asset transaction amount 1000, got %d", assetTransaction.Amount)
			}
			if assetTransaction.SourceBankId != bank.Id {
				t.Errorf("Expected asset transaction SourceBankId to be %s, got %s", bank.Id.Hex(), assetTransaction.SourceBankId.Hex())
			}
		}

		// Verify cash transaction
		if cashTransaction == nil {
			t.Errorf("Expected to find cash transaction")
		} else {
			if cashTransaction.Amount != -1000 {
				t.Errorf("Expected cash transaction amount -1000, got %d", cashTransaction.Amount)
			}
			if cashTransaction.SourceBankId != bank.Id {
				t.Errorf("Expected cash transaction SourceBankId to be %s, got %s", bank.Id.Hex(), cashTransaction.SourceBankId.Hex())
			}
		}
	})

	t.Run("Valid sell transaction creation", func(t *testing.T) {
		// Clear previous transactions
		existingTransactions, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existingTransactions {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}

		err := service.CreateSellTransaction(ctx, bank.Id, assetType.Id, 500, username)
		if err != nil {
			t.Errorf("Expected no error for valid sell transaction, got: %v", err)
		}

		// Verify transactions were created (should be 2: asset sale + cash addition)
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions (asset + cash), got %d", len(transactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Find the asset transaction and cash transaction
		var assetTransaction, cashTransaction *models.PendingTransactionResponse
		for i := range transactions {
			if transactions[i].TargetAssetId == assetType.Id {
				assetTransaction = &transactions[i]
			} else if transactions[i].TargetAssetId == cashAssetType.Id {
				cashTransaction = &transactions[i]
			}
		}

		// Verify asset transaction (should be negative for sale)
		if assetTransaction == nil {
			t.Errorf("Expected to find asset transaction for asset ID %s", assetType.Id.Hex())
		} else {
			if assetTransaction.Amount != -500 {
				t.Errorf("Expected asset transaction amount -500, got %d", assetTransaction.Amount)
			}
			if assetTransaction.SourceBankId != bank.Id {
				t.Errorf("Expected asset transaction SourceBankId to be %s, got %s", bank.Id.Hex(), assetTransaction.SourceBankId.Hex())
			}
		}

		// Verify cash transaction (should be positive for sale)
		if cashTransaction == nil {
			t.Errorf("Expected to find cash transaction")
		} else {
			if cashTransaction.Amount != 500 {
				t.Errorf("Expected cash transaction amount 500, got %d", cashTransaction.Amount)
			}
			if cashTransaction.SourceBankId != bank.Id {
				t.Errorf("Expected cash transaction SourceBankId to be %s, got %s", bank.Id.Hex(), cashTransaction.SourceBankId.Hex())
			}
		}
	})

	t.Run("Zero amount rejected", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 0, username)
		if err != services.ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount for zero amount in buy transaction, got: %v", err)
		}

		err = service.CreateSellTransaction(ctx, bank.Id, assetType.Id, 0, username)
		if err != services.ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount for zero amount in sell transaction, got: %v", err)
		}
	})

	t.Run("Negative amount rejected", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, -100, username)
		if err != services.ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount for negative amount in buy transaction, got: %v", err)
		}

		err = service.CreateSellTransaction(ctx, bank.Id, assetType.Id, -100, username)
		if err != services.ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount for negative amount in sell transaction, got: %v", err)
		}
	})

	t.Run("Non-existent bank", func(t *testing.T) {
		nonExistentBankID := primitive.NewObjectID()
		err := service.CreateBuyTransaction(ctx, nonExistentBankID, assetType.Id, 1000, username)
		if err != services.ErrInvalidBankID {
			t.Errorf("Expected ErrInvalidBankID for non-existent bank, got: %v", err)
		}
		err = service.CreateSellTransaction(ctx, nonExistentBankID, assetType.Id, 1000, username)
		if err != services.ErrInvalidBankID {
			t.Errorf("Expected ErrInvalidBankID for non-existent bank in sell transaction, got: %v", err)
		}
	})

	t.Run("Non-existent asset", func(t *testing.T) {
		nonExistentAssetID := primitive.NewObjectID()
		err := service.CreateBuyTransaction(ctx, bank.Id, nonExistentAssetID, 1000, username)
		if err != services.ErrTargetAssetNotFound {
			t.Errorf("Expected ErrAssetNotFound for non-existent asset, got: %v", err)
		}
		err = service.CreateSellTransaction(ctx, bank.Id, nonExistentAssetID, 1000, username)
		if err != services.ErrTargetAssetNotFound {
			t.Errorf("Expected ErrAssetNotFound for non-existent asset in sell transaction, got: %v", err)
		}
	})

	t.Run("Self-investment", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, bank.Id, 1000, username)
		if err != services.ErrSelfInvestment {
			t.Errorf("Expected ErrSelfInvestment for self-investment, got: %v", err)
		}
		err = service.CreateSellTransaction(ctx, bank.Id, bank.Id, 1000, username)
		if err != services.ErrSelfInvestment {
			t.Errorf("Expected ErrSelfInvestment for self-investment in sell transaction, got: %v", err)
		}
	})

	t.Run("Non-existent user", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 1000, "nonexistentuser")
		if err != services.ErrPlayerNotFound {
			t.Errorf("Expected ErrInvalidBankID when user doesn't exist, got: %v", err)
		}
		err = service.CreateSellTransaction(ctx, bank.Id, assetType.Id, 1000, "nonexistentuser")
		if err != services.ErrPlayerNotFound {
			t.Errorf("Expected ErrInvalidBankID for non-existent user in sell transaction, got: %v", err)
		}
	})
}

func TestPendingTransactionService_BankOwnership(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_ownership")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create two users
	user1Username := fmt.Sprintf("testuser1_%d", timestamp)
	user1Password := "testpass1"
	user1BankName := "Test Bank 1"

	_, err = CreateRegularUserForTest(container, user1Username, user1Password, user1BankName)
	if err != nil {
		t.Fatalf("Failed to create first test user: %v", err)
	}

	user2Username := fmt.Sprintf("testuser2_%d", timestamp)
	user2Password := "testpass2"
	user2BankName := "Test Bank 2"

	_, err = CreateRegularUserForTest(container, user2Username, user2Password, user2BankName)
	if err != nil {
		t.Fatalf("Failed to create second test user: %v", err)
	}

	user1, err := container.RepositoryContainer.Player.FindByUsername(ctx, user1Username)
	if err != nil {
		t.Fatalf("Failed to find first test user: %v", err)
	}

	user2, err := container.RepositoryContainer.Player.FindByUsername(ctx, user2Username)
	if err != nil {
		t.Fatalf("Failed to find second test user: %v", err)
	}

	banks1, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user1.Id)
	if err != nil { 
		t.Fatalf("Failed to find first test banks: %v", err)
	}
	bank1 := banks1[0]
	
	banks2, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user2.Id)
	if err != nil { 
		t.Fatalf("Failed to find second test banks: %v", err)
	}
	bank2 := banks2[0]

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("User owns bank", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank1.Id, assetType.Id, 1000, user1Username)
		if err != nil {
			t.Errorf("Expected no error when user uses their own bank for buy, got: %v", err)
		}
		err = service.CreateSellTransaction(ctx, bank1.Id, assetType.Id, 500, user1Username)
		if err != nil {
			t.Errorf("Expected no error when user uses their own bank for sell, got: %v", err)
		}
	})

	t.Run("User does not own bank", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank2.Id, assetType.Id, 1000, user1Username)
		if err != services.ErrUnauthorizedBank {
			t.Errorf("Expected ErrUnauthorizedBank when user tries to use another user's bank for buy, got: %v", err)
		}
		err = service.CreateSellTransaction(ctx, bank2.Id, assetType.Id, 500, user1Username)
		if err != services.ErrUnauthorizedBank {
			t.Errorf("Expected ErrUnauthorizedBank when user tries to use another user's bank for sell, got: %v", err)
		}
	})
}

func TestPendingTransactionService_MultipleAssets(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_multiple")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create test user and bank
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "testpass"
	bankName := "Test Bank"

	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user.Id)
	if err != nil { 
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Create multiple asset types
	assetType1 := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Asset Type 1",
	}
	assetType2 := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Asset Type 2",
	}

	err = container.RepositoryContainer.AssetType.Create(ctx, assetType1)
	if err != nil {
		t.Fatalf("Failed to create asset type 1: %v", err)
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType2)
	if err != nil {
		t.Fatalf("Failed to create asset type 2: %v", err)
	}

	t.Run("Create transactions for different assets", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType1.Id, 400, username)
		if err != nil {
			t.Fatalf("Failed to create transaction for asset 1: %v", err)
		}

		err = service.CreateBuyTransaction(ctx, bank.Id, assetType2.Id, 600, username)
		if err != nil {
			t.Fatalf("Failed to create transaction for asset 2: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		// Should have 3 transactions: asset1 (+400), asset2 (+600), cash (-1000 combined)
		if len(transactions) != 3 {
			t.Errorf("Expected 3 transactions (2 assets + 1 combined cash), got %d", len(transactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Verify amounts and types are correct
		var asset1Found, asset2Found, cashFound bool
		for _, transaction := range transactions {
			if transaction.TargetAssetId == assetType1.Id {
				if transaction.Amount != 400 {
					t.Errorf("Expected amount 400 for asset 1, got %d", transaction.Amount)
				}
				asset1Found = true
			} else if transaction.TargetAssetId == assetType2.Id {
				if transaction.Amount != 600 {
					t.Errorf("Expected amount 600 for asset 2, got %d", transaction.Amount)
				}
				asset2Found = true
			} else if transaction.TargetAssetId == cashAssetType.Id {
				if transaction.Amount != -1000 {
					t.Errorf("Expected combined cash amount -1000, got %d", transaction.Amount)
				}
				cashFound = true
			}
		}

		if !asset1Found {
			t.Error("Expected to find asset 1 transaction")
		}
		if !asset2Found {
			t.Error("Expected to find asset 2 transaction")
		}
		if !cashFound {
			t.Error("Expected to find cash transaction")
		}
	})
	t.Run("Create sell transactions for different assets", func(t *testing.T) {
		// Clear previous transactions
		existing, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existing {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}
		
		// Create sell transactions for two assets
		err := service.CreateSellTransaction(ctx, bank.Id, assetType1.Id, 300, username)
		if err != nil {
			t.Fatalf("Failed to create sell transaction for asset 1: %v", err)
		}
		err = service.CreateSellTransaction(ctx, bank.Id, assetType2.Id, 400, username)
		if err != nil {
			t.Fatalf("Failed to create sell transaction for asset 2: %v", err)
		}
		
		sellTransactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get sell transactions: %v", err)
		}
		// Should have 3 transactions: asset1 (-300), asset2 (-400), cash (+700 combined)
		if len(sellTransactions) != 3 {
			t.Errorf("Expected 3 transactions (2 assets + 1 combined cash), got %d", len(sellTransactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Verify amounts and types are correct
		var asset1Found, asset2Found, cashFound bool
		for _, transaction := range sellTransactions {
			if transaction.TargetAssetId == assetType1.Id {
				if transaction.Amount != -300 {
					t.Errorf("Expected amount -300 for sell asset 1, got %d", transaction.Amount)
				}
				asset1Found = true
			} else if transaction.TargetAssetId == assetType2.Id {
				if transaction.Amount != -400 {
					t.Errorf("Expected amount -400 for sell asset 2, got %d", transaction.Amount)
				}
				asset2Found = true
			} else if transaction.TargetAssetId == cashAssetType.Id {
				if transaction.Amount != 700 {
					t.Errorf("Expected combined cash amount +700, got %d", transaction.Amount)
				}
				cashFound = true
			}
		}

		if !asset1Found {
			t.Error("Expected to find asset 1 sell transaction")
		}
		if !asset2Found {
			t.Error("Expected to find asset 2 sell transaction")
		}
		if !cashFound {
			t.Error("Expected to find cash transaction")
		}
	})
}

func TestPendingTransactionService_GetTransactionsByBankID(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_get_by_bank")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create test user and bank
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "testpass"
	bankName := "Test Bank"

	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user.Id)
	if err != nil { 
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Create second test user and bank for unauthorized access test
	username2 := fmt.Sprintf("testuser2_%d", timestamp)
	password2 := "testpass2"
	bankName2 := "Test Bank 2"

	_, err = CreateRegularUserForTest(container, username2, password2, bankName2)
	if err != nil {
		t.Fatalf("Failed to create second test user: %v", err)
	}

	user2, err := container.RepositoryContainer.Player.FindByUsername(ctx, username2)
	if err != nil {
		t.Fatalf("Failed to find second test user: %v", err)
	}

	banks2, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user2.Id)
	if err != nil { 
		t.Fatalf("Failed to find second test banks: %v", err)
	}
	bank2 := banks2[0]

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("Valid bank owner can access transactions", func(t *testing.T) {
		// Create some pending transactions
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 600, username)
		if err != nil {
			t.Fatalf("Failed to create first transaction: %v", err)
		}

		err = service.CreateBuyTransaction(ctx, bank.Id, bank2.Id, 400, username)
		if err != nil {
			t.Fatalf("Failed to create second transaction: %v", err)
		}

		// Get transactions using the new method
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Errorf("Expected no error for valid bank owner, got: %v", err)
		}

		// Should have 3 transactions: asset (+600), bank (+400), cash (-1000 combined)
		if len(transactions) != 3 {
			t.Errorf("Expected 3 transactions (2 purchases + 1 combined cash), got %d", len(transactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Verify transaction details
		var found600, found400, foundCash bool
		for _, tx := range transactions {
			if tx.Amount == 600 && tx.TargetAssetId == assetType.Id {
				found600 = true
			} else if tx.Amount == 400 && tx.TargetAssetId == bank2.Id {
				found400 = true
			} else if tx.Amount == -1000 && tx.TargetAssetId == cashAssetType.Id {
				foundCash = true
			}
		}

		if !found600 {
			t.Error("Expected to find transaction with amount 600")
		}
		if !found400 {
			t.Error("Expected to find transaction with amount 400")
		}
		if !foundCash {
			t.Error("Expected to find cash transaction with amount -1000")
		}
	})

	t.Run("Unauthorized user cannot access transactions", func(t *testing.T) {
		// Try to access bank's transactions with different user
		_, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username2)
		if err != services.ErrUnauthorizedBank {
			t.Errorf("Expected ErrUnauthorizedBank for unauthorized user, got: %v", err)
		}
	})

	t.Run("Non-existent bank", func(t *testing.T) {
		nonExistentBankID := primitive.NewObjectID()
		_, err := service.GetTransactionsByBuyerBankID(ctx, nonExistentBankID, username)
		if err != services.ErrInvalidBankID {
			t.Errorf("Expected ErrInvalidBankID for non-existent bank, got: %v", err)
		}
	})

	t.Run("Non-existent user", func(t *testing.T) {
		_, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, "nonexistentuser")
		if err != services.ErrPlayerNotFound {
			t.Errorf("Expected ErrPlayerNotFound for non-existent user, got: %v", err)
		}
	})

	t.Run("Empty transactions list", func(t *testing.T) {
		// bank2 should have no pending transactions as buyer
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank2.Id, username2)
		if err != nil {
			t.Errorf("Expected no error for valid bank owner with no transactions, got: %v", err)
		}

		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions for bank2, got %d", len(transactions))
		}
	})
}

func TestPendingTransactionService_CreateBuyTransaction(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_buy")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create test user and bank
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "testpass"
	bankName := "Test Bank"

	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user.Id)
	if err != nil { 
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("Multiple buy transactions combine", func(t *testing.T) {
		// Clear any existing transactions
		existingTransactions, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existingTransactions {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}

		// Create first buy transaction
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 600, username)
		if err != nil {
			t.Fatalf("Failed to create first buy transaction: %v", err)
		}

		// Create second buy transaction
		err = service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 300, username)
		if err != nil {
			t.Fatalf("Failed to create second buy transaction: %v", err)
		}

		// Verify transactions were combined
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		// Should have 2 transactions: combined asset (+900) and combined cash (-900)
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions (combined asset + combined cash), got %d", len(transactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Find and verify the combined transactions
		var assetTransaction, cashTransaction *models.PendingTransactionResponse
		for i := range transactions {
			if transactions[i].TargetAssetId == assetType.Id {
				assetTransaction = &transactions[i]
			} else if transactions[i].TargetAssetId == cashAssetType.Id {
				cashTransaction = &transactions[i]
			}
		}

		if assetTransaction == nil {
			t.Error("Expected to find combined asset transaction")
		} else if assetTransaction.Amount != 900 {
			t.Errorf("Expected combined asset amount 900, got %d", assetTransaction.Amount)
		}

		if cashTransaction == nil {
			t.Error("Expected to find combined cash transaction")
		} else if cashTransaction.Amount != -900 {
			t.Errorf("Expected combined cash amount -900, got %d", cashTransaction.Amount)
		}
	})
}

func TestPendingTransactionService_CreateSellTransaction(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_sell")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create test user and bank
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "testpass"
	bankName := "Test Bank"

	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user.Id)
	if err != nil { 
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("Multiple sell transactions combine", func(t *testing.T) {
		// Clear any existing transactions
		existingTransactions, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existingTransactions {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}

		// Create first sell transaction
		err := service.CreateSellTransaction(ctx, bank.Id, assetType.Id, 300, username)
		if err != nil {
			t.Fatalf("Failed to create first sell transaction: %v", err)
		}

		// Create second sell transaction
		err = service.CreateSellTransaction(ctx, bank.Id, assetType.Id, 200, username)
		if err != nil {
			t.Fatalf("Failed to create second sell transaction: %v", err)
		}

		// Verify transactions were combined
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		// Should have 2 transactions: combined asset (-500) and combined cash (+500)
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions (combined asset + combined cash), got %d", len(transactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Find and verify the combined transactions
		var assetTransaction, cashTransaction *models.PendingTransactionResponse
		for i := range transactions {
			if transactions[i].TargetAssetId == assetType.Id {
				assetTransaction = &transactions[i]
			} else if transactions[i].TargetAssetId == cashAssetType.Id {
				cashTransaction = &transactions[i]
			}
		}

		if assetTransaction == nil {
			t.Error("Expected to find combined asset transaction")
		} else if assetTransaction.Amount != -500 {
			t.Errorf("Expected combined asset amount -500, got %d", assetTransaction.Amount)
		}

		if cashTransaction == nil {
			t.Error("Expected to find combined cash transaction")
		} else if cashTransaction.Amount != 500 {
			t.Errorf("Expected combined cash amount +500, got %d", cashTransaction.Amount)
		}
	})
}

func TestPendingTransactionService_BuyAndSellCombination(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_buy_sell")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create test user and bank
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "testpass"
	bankName := "Test Bank"

	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user.Id)
	if err != nil { 
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("Buy then sell reduces amount", func(t *testing.T) {
		// Create buy transaction
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 1000, username)
		if err != nil {
			t.Fatalf("Failed to create buy transaction: %v", err)
		}

		// Create sell transaction
		err = service.CreateSellTransaction(ctx, bank.Id, assetType.Id, 300, username)
		if err != nil {
			t.Fatalf("Failed to create sell transaction: %v", err)
		}

		// Verify final amounts
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		// Should have 2 transactions: asset (+700) and cash (-700)
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions (asset + cash), got %d", len(transactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Find and verify the final transactions
		var assetTransaction, cashTransaction *models.PendingTransactionResponse
		for i := range transactions {
			if transactions[i].TargetAssetId == assetType.Id {
				assetTransaction = &transactions[i]
			} else if transactions[i].TargetAssetId == cashAssetType.Id {
				cashTransaction = &transactions[i]
			}
		}

		if assetTransaction == nil {
			t.Error("Expected to find asset transaction")
		} else if assetTransaction.Amount != 700 {
			t.Errorf("Expected asset amount 700 (1000 - 300), got %d", assetTransaction.Amount)
		}

		if cashTransaction == nil {
			t.Error("Expected to find cash transaction")
		} else if cashTransaction.Amount != -700 {
			t.Errorf("Expected cash amount -700, got %d", cashTransaction.Amount)
		}
	})

	t.Run("Sell then buy reduces sell amount", func(t *testing.T) {
		// Clear previous transactions
		existingTransactions, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existingTransactions {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}

		// Create sell transaction first
		err := service.CreateSellTransaction(ctx, bank.Id, assetType.Id, 500, username)
		if err != nil {
			t.Fatalf("Failed to create sell transaction: %v", err)
		}

		// Create buy transaction
		err = service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 200, username)
		if err != nil {
			t.Fatalf("Failed to create buy transaction: %v", err)
		}

		// Verify final amounts
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		// Should have 2 transactions: asset (-300) and cash (+300)
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions (asset + cash), got %d", len(transactions))
		}

		// Get cash asset type for verification
		cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("Failed to get cash asset type: %v", err)
		}

		// Find and verify the final transactions
		var assetTransaction, cashTransaction *models.PendingTransactionResponse
		for i := range transactions {
			if transactions[i].TargetAssetId == assetType.Id {
				assetTransaction = &transactions[i]
			} else if transactions[i].TargetAssetId == cashAssetType.Id {
				cashTransaction = &transactions[i]
			}
		}

		if assetTransaction == nil {
			t.Error("Expected to find asset transaction")
		} else if assetTransaction.Amount != -300 {
			t.Errorf("Expected asset amount -300 (-500 + 200), got %d", assetTransaction.Amount)
		}

		if cashTransaction == nil {
			t.Error("Expected to find cash transaction")
		} else if cashTransaction.Amount != 300 {
			t.Errorf("Expected cash amount +300, got %d", cashTransaction.Amount)
		}
	})

	t.Run("Equal buy and sell cancel out", func(t *testing.T) {
		// Clear previous transactions
		existingTransactions, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existingTransactions {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}

		// Create buy transaction
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 400, username)
		if err != nil {
			t.Fatalf("Failed to create buy transaction: %v", err)
		}

		// Create equal sell transaction
		err = service.CreateSellTransaction(ctx, bank.Id, assetType.Id, 400, username)
		if err != nil {
			t.Fatalf("Failed to create sell transaction: %v", err)
		}

		// Verify transaction was deleted
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions after cancellation, got %d", len(transactions))
		}
	})
}

func TestPendingTransactionService_CashRestriction(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_cash")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create test user and bank
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "testpass"
	bankName := "Test Bank"

	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	// Get the banks for the user
	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user.Id)
	if err != nil {
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Get the Cash asset type
	cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
	if err != nil {
		t.Fatalf("Failed to find Cash asset type: %v", err)
	}

	t.Run("Cannot buy cash", func(t *testing.T) {
		// Try to create a buy transaction for cash
		err := service.CreateBuyTransaction(ctx, bank.Id, cashAssetType.Id, 100, username)
		if err == nil {
			t.Error("Expected error when trying to buy cash, got nil")
		}
		if err != services.ErrCashNotTradable {
			t.Errorf("Expected ErrCashNotTradable, got %v", err)
		}
	})

	t.Run("Cannot sell cash", func(t *testing.T) {
		// Try to create a sell transaction for cash
		err := service.CreateSellTransaction(ctx, bank.Id, cashAssetType.Id, 100, username)
		if err == nil {
			t.Error("Expected error when trying to sell cash, got nil")
		}
		if err != services.ErrCashNotTradable {
			t.Errorf("Expected ErrCashNotTradable, got %v", err)
		}
	})
}

func TestPendingTransactionService_InsufficientFunds(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_funds")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create test user and bank
	username := fmt.Sprintf("testuser_%d", timestamp)
	password := "testpass"
	bankName := "Test Bank"

	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	// Get the banks for the user
	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, user.Id)
	if err != nil {
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	// Get the user's initial cash balance (should be 1000 from CreateRegularUserForTest)
	_, err = container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
	if err != nil {
		t.Fatalf("Failed to find Cash asset type: %v", err)
	}

	// Try to buy more than the available cash balance (1000)
	t.Run("Cannot buy more than available cash", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 1500, username)
		if err == nil {
			t.Error("Expected error when trying to buy more than available cash, got nil")
		}
		if err != services.ErrInsufficientFunds {
			t.Errorf("Expected ErrInsufficientFunds, got %v", err)
		}
	})

	// Should be able to buy exactly the available cash balance
	t.Run("Can buy exactly available cash amount", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 1000, username)
		if err != nil {
			t.Errorf("Should be able to buy with exact cash balance, got error: %v", err)
		}
	})

	// After using all cash, should not be able to buy more
	t.Run("Cannot buy after using all cash", func(t *testing.T) {
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 1, username)
		if err == nil {
			t.Error("Expected error when trying to buy after using all cash, got nil")
		}
		if err != services.ErrInsufficientFunds {
			t.Errorf("Expected ErrInsufficientFunds, got %v", err)
		}
	})
}

func TestPendingTransactionService_DualTransactionSystem(t *testing.T) {
	container, err := CreateTestDependencies("dual_transaction_system")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	username := "dualtestuser"
	password := "password"
	bankName := "Dual Test Bank"

	// Create test user with initial cash
	_, err = CreateRegularUserForTest(container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the bank for the user
	player, err := container.RepositoryContainer.Player.FindByUsername(ctx, username)
	if err != nil {
		t.Fatalf("Failed to find test player: %v", err)
	}

	banks, err := container.RepositoryContainer.Bank.FindAllByPlayerID(ctx, player.Id)
	if err != nil {
		t.Fatalf("Failed to find test banks: %v", err)
	}
	bank := banks[0]

	// Create test asset types
	assetType1 := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Stock A",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType1)
	if err != nil {
		t.Fatalf("Failed to create asset type 1: %v", err)
	}

	assetType2 := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Stock B",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType2)
	if err != nil {
		t.Fatalf("Failed to create asset type 2: %v", err)
	}

	// Get cash asset type
	cashAssetType, err := container.RepositoryContainer.AssetType.FindByName(ctx, "Cash")
	if err != nil {
		t.Fatalf("Failed to get cash asset type: %v", err)
	}
	
	t.Run("Multiple transactions combine cash correctly", func(t *testing.T) {
		// Clear any existing transactions
		existing, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existing {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}

		// Create multiple buy transactions for different assets
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType1.Id, 200, username)
		if err != nil {
			t.Fatalf("Failed to create first buy transaction: %v", err)
		}

		err = service.CreateBuyTransaction(ctx, bank.Id, assetType2.Id, 300, username)
		if err != nil {
			t.Fatalf("Failed to create second buy transaction: %v", err)
		}

		// Should have 3 transactions: asset1 (+200), asset2 (+300), cash (-500 combined)
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 3 {
			t.Errorf("Expected 3 transactions, got %d", len(transactions))
		}

		// Verify the cash transaction is properly combined
		var cashTx *models.PendingTransactionResponse
		for i := range transactions {
			if transactions[i].TargetAssetId == cashAssetType.Id {
				cashTx = &transactions[i]
				break
			}
		}

		if cashTx == nil {
			t.Error("Expected to find combined cash transaction")
		} else {
			if cashTx.Amount != -500 {
				t.Errorf("Expected combined cash transaction amount -500, got %d", cashTx.Amount)
			}
		}
	})

	t.Run("Buy and sell transactions net out correctly", func(t *testing.T) {
		// Clear any existing transactions
		existing, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existing {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}

		// Buy and then sell the same amount should result in no transactions
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType1.Id, 400, username)
		if err != nil {
			t.Fatalf("Failed to create buy transaction: %v", err)
		}

		err = service.CreateSellTransaction(ctx, bank.Id, assetType1.Id, 400, username)
		if err != nil {
			t.Fatalf("Failed to create sell transaction: %v", err)
		}

		// Should have no transactions left (they should cancel out)
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions (they should cancel out), got %d", len(transactions))
		}
	})

	t.Run("Cannot directly buy or sell cash", func(t *testing.T) {
		// Clear any existing transactions
		existing, _ := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		for _, tx := range existing {
			container.RepositoryContainer.PendingTransaction.Delete(ctx, tx.Id)
		}

		// Attempting to buy cash should fail
		err := service.CreateBuyTransaction(ctx, bank.Id, cashAssetType.Id, 100, username)
		if err == nil {
			t.Error("Expected error when trying to buy cash, got nil")
		}

		// Attempting to sell cash should fail
		err = service.CreateSellTransaction(ctx, bank.Id, cashAssetType.Id, 100, username)
		if err == nil {
			t.Error("Expected error when trying to sell cash, got nil")
		}

		// Should have no transactions created
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions, got %d", len(transactions))
		}
	})
}
