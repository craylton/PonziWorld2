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

		// Verify transaction was created
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != 1000 {
			t.Errorf("Expected amount 1000, got %d", transactions[0].Amount)
		}
		if transactions[0].TargetAssetId != assetType.Id {
			t.Errorf("Expected TargetAssetId to be %s, got %s", assetType.Id.Hex(), transactions[0].TargetAssetId.Hex())
		}
		if transactions[0].SourceBankId != bank.Id {
			t.Errorf("Expected SourceBankId to be %s, got %s", bank.Id.Hex(), transactions[0].SourceBankId.Hex())
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

		// Verify transaction was created with negative amount (internal representation)
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != -500 {
			t.Errorf("Expected amount -500 (internal representation), got %d", transactions[0].Amount)
		}
		if transactions[0].TargetAssetId != assetType.Id {
			t.Errorf("Expected TargetAssetId to be %s, got %s", assetType.Id.Hex(), transactions[0].TargetAssetId.Hex())
		}
		if transactions[0].SourceBankId != bank.Id {
			t.Errorf("Expected SourceBankId to be %s, got %s", bank.Id.Hex(), transactions[0].SourceBankId.Hex())
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
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType1.Id, 1000, username)
		if err != nil {
			t.Fatalf("Failed to create transaction for asset 1: %v", err)
		}

		err = service.CreateBuyTransaction(ctx, bank.Id, assetType2.Id, 2000, username)
		if err != nil {
			t.Fatalf("Failed to create transaction for asset 2: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions for different assets, got %d", len(transactions))
		}

		// Verify amounts are correct
		for _, transaction := range transactions {
			if transaction.TargetAssetId == assetType1.Id && transaction.Amount != 1000 {
				t.Errorf("Expected amount 1000 for asset 1, got %d", transaction.Amount)
			}
			if transaction.TargetAssetId == assetType2.Id && transaction.Amount != 2000 {
				t.Errorf("Expected amount 2000 for asset 2, got %d", transaction.Amount)
			}
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
		if len(sellTransactions) != 2 {
			t.Errorf("Expected 2 sell transactions for different assets, got %d", len(sellTransactions))
		}
		// Verify negative amounts are correct
		for _, transaction := range sellTransactions {
			if transaction.TargetAssetId == assetType1.Id && transaction.Amount != -300 {
				t.Errorf("Expected amount -300 for sell asset 1, got %d", transaction.Amount)
			}
			if transaction.TargetAssetId == assetType2.Id && transaction.Amount != -400 {
				t.Errorf("Expected amount -400 for sell asset 2, got %d", transaction.Amount)
			}
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
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 1000, username)
		if err != nil {
			t.Fatalf("Failed to create first transaction: %v", err)
		}

		err = service.CreateBuyTransaction(ctx, bank.Id, bank2.Id, 500, username)
		if err != nil {
			t.Fatalf("Failed to create second transaction: %v", err)
		}

		// Get transactions using the new method
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Errorf("Expected no error for valid bank owner, got: %v", err)
		}

		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions, got %d", len(transactions))
		}

		// Verify transaction details
		found1000 := false
		found500 := false
		for _, tx := range transactions {
			if tx.Amount == 1000 && tx.TargetAssetId == assetType.Id {
				found1000 = true
			}
			if tx.Amount == 500 && tx.TargetAssetId == bank2.Id {
				found500 = true
			}
		}

		if !found1000 {
			t.Error("Expected to find transaction with amount 1000")
		}
		if !found500 {
			t.Error("Expected to find transaction with amount 500")
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
		err := service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 1000, username)
		if err != nil {
			t.Fatalf("Failed to create first buy transaction: %v", err)
		}

		// Create second buy transaction
		err = service.CreateBuyTransaction(ctx, bank.Id, assetType.Id, 500, username)
		if err != nil {
			t.Fatalf("Failed to create second buy transaction: %v", err)
		}

		// Verify transactions were combined
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("Expected 1 combined transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != 1500 {
			t.Errorf("Expected combined amount 1500, got %d", transactions[0].Amount)
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
		if len(transactions) != 1 {
			t.Errorf("Expected 1 combined transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != -500 {
			t.Errorf("Expected combined amount -500, got %d", transactions[0].Amount)
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

		// Verify final amount
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != 700 {
			t.Errorf("Expected amount 700 (1000 - 300), got %d", transactions[0].Amount)
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

		// Verify final amount
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id, username)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != -300 {
			t.Errorf("Expected amount -300 (-500 + 200), got %d", transactions[0].Amount)
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
