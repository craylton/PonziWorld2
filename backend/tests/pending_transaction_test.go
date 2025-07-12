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

func TestPendingTransactionService_CreateTransaction(t *testing.T) {
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

	bank, err := container.RepositoryContainer.Bank.FindByPlayerID(ctx, user.Id)
	if err != nil {
		t.Fatalf("Failed to find test bank: %v", err)
	}

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("Valid transaction creation", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, assetType.Id, 1000, username)
		if err != nil {
			t.Errorf("Expected no error for valid transaction, got: %v", err)
		}

		// Verify transaction was created
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != 1000 {
			t.Errorf("Expected amount 1000, got %d", transactions[0].Amount)
		}
	})

	t.Run("Zero amount", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, assetType.Id, 0, username)
		if err != services.ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount for zero amount, got: %v", err)
		}
	})

	t.Run("Non-existent bank", func(t *testing.T) {
		nonExistentBankID := primitive.NewObjectID()
		err := service.CreateTransaction(ctx, nonExistentBankID, assetType.Id, 1000, username)
		if err != services.ErrInvalidBankID {
			t.Errorf("Expected ErrInvalidBankID for non-existent bank, got: %v", err)
		}
	})

	t.Run("Non-existent asset", func(t *testing.T) {
		nonExistentAssetID := primitive.NewObjectID()
		err := service.CreateTransaction(ctx, bank.Id, nonExistentAssetID, 1000, username)
		if err != services.ErrAssetNotFound {
			t.Errorf("Expected ErrAssetNotFound for non-existent asset, got: %v", err)
		}
	})

	t.Run("Self-investment", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, bank.Id, 1000, username)
		if err != services.ErrSelfInvestment {
			t.Errorf("Expected ErrSelfInvestment for self-investment, got: %v", err)
		}
	})

	t.Run("Non-existent user", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, assetType.Id, 1000, "nonexistentuser")
		if err != services.ErrInvalidBankID {
			t.Errorf("Expected ErrInvalidBankID when user doesn't exist, got: %v", err)
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

	bank1, err := container.RepositoryContainer.Bank.FindByPlayerID(ctx, user1.Id)
	if err != nil {
		t.Fatalf("Failed to find first test bank: %v", err)
	}

	bank2, err := container.RepositoryContainer.Bank.FindByPlayerID(ctx, user2.Id)
	if err != nil {
		t.Fatalf("Failed to find second test bank: %v", err)
	}

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
		err := service.CreateTransaction(ctx, bank1.Id, assetType.Id, 1000, user1Username)
		if err != nil {
			t.Errorf("Expected no error when user uses their own bank, got: %v", err)
		}
	})

	t.Run("User does not own bank", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank2.Id, assetType.Id, 1000, user1Username)
		if err != services.ErrUnauthorizedBank {
			t.Errorf("Expected ErrUnauthorizedBank when user tries to use another user's bank, got: %v", err)
		}
	})
}

func TestPendingTransactionService_TransactionCombining(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_combining")
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

	bank, err := container.RepositoryContainer.Bank.FindByPlayerID(ctx, user.Id)
	if err != nil {
		t.Fatalf("Failed to find test bank: %v", err)
	}

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("Initial buy transaction", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, assetType.Id, 1000, username)
		if err != nil {
			t.Fatalf("Failed to create initial transaction: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Fatalf("Expected 1 transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != 1000 {
			t.Errorf("Expected amount 1000, got %d", transactions[0].Amount)
		}
	})

	t.Run("Additional buy transaction combines", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, assetType.Id, 500, username)
		if err != nil {
			t.Fatalf("Failed to create second transaction: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id)
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

	t.Run("Sell transaction reduces amount", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, assetType.Id, -800, username)
		if err != nil {
			t.Fatalf("Failed to create sell transaction: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction after sell, got %d", len(transactions))
		}
		if transactions[0].Amount != 700 {
			t.Errorf("Expected reduced amount 700, got %d", transactions[0].Amount)
		}
	})

	t.Run("Sell all deletes transaction", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, assetType.Id, -700, username)
		if err != nil {
			t.Fatalf("Failed to create final sell transaction: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions after selling all, got %d", len(transactions))
		}
	})
}

func TestPendingTransactionService_BankAsAsset(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_bank_asset")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	service := container.ServiceContainer.PendingTransaction
	timestamp := time.Now().Unix()

	// Create two users with banks
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

	bank1, err := container.RepositoryContainer.Bank.FindByPlayerID(ctx, user1.Id)
	if err != nil {
		t.Fatalf("Failed to find first test bank: %v", err)
	}

	bank2, err := container.RepositoryContainer.Bank.FindByPlayerID(ctx, user2.Id)
	if err != nil {
		t.Fatalf("Failed to find second test bank: %v", err)
	}

	t.Run("Invest in another bank", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank1.Id, bank2.Id, 1000, user1Username)
		if err != nil {
			t.Errorf("Expected no error when investing in another bank, got: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank1.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}
		if transactions[0].AssetId != bank2.Id {
			t.Errorf("Expected AssetId to be bank2 ID, got different ID")
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

	bank, err := container.RepositoryContainer.Bank.FindByPlayerID(ctx, user.Id)
	if err != nil {
		t.Fatalf("Failed to find test bank: %v", err)
	}

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
		err := service.CreateTransaction(ctx, bank.Id, assetType1.Id, 1000, username)
		if err != nil {
			t.Fatalf("Failed to create transaction for asset 1: %v", err)
		}

		err = service.CreateTransaction(ctx, bank.Id, assetType2.Id, 2000, username)
		if err != nil {
			t.Fatalf("Failed to create transaction for asset 2: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions for different assets, got %d", len(transactions))
		}

		// Verify amounts are correct
		for _, txn := range transactions {
			if txn.AssetId == assetType1.Id && txn.Amount != 1000 {
				t.Errorf("Expected amount 1000 for asset 1, got %d", txn.Amount)
			}
			if txn.AssetId == assetType2.Id && txn.Amount != 2000 {
				t.Errorf("Expected amount 2000 for asset 2, got %d", txn.Amount)
			}
		}
	})

	t.Run("Add to existing asset combines", func(t *testing.T) {
		err := service.CreateTransaction(ctx, bank.Id, assetType1.Id, 500, username)
		if err != nil {
			t.Fatalf("Failed to add to transaction for asset 1: %v", err)
		}

		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions after combining, got %d", len(transactions))
		}

		for _, txn := range transactions {
			if txn.AssetId == assetType1.Id && txn.Amount != 1500 {
				t.Errorf("Expected combined amount 1500 for asset 1, got %d", txn.Amount)
			}
			if txn.AssetId == assetType2.Id && txn.Amount != 2000 {
				t.Errorf("Expected unchanged amount 2000 for asset 2, got %d", txn.Amount)
			}
		}
	})
}

func TestPendingTransactionService_GetTransactions(t *testing.T) {
	container, err := CreateTestDependencies("pending_transaction_get")
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

	bank, err := container.RepositoryContainer.Bank.FindByPlayerID(ctx, user.Id)
	if err != nil {
		t.Fatalf("Failed to find test bank: %v", err)
	}

	// Create test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = container.RepositoryContainer.AssetType.Create(ctx, assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	t.Run("No transactions initially", func(t *testing.T) {
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, bank.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions initially, got %d", len(transactions))
		}

		transactions, err = service.GetTransactionsByAssetID(ctx, assetType.Id)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions initially, got %d", len(transactions))
		}
	})

	t.Run("Non-existent IDs return empty", func(t *testing.T) {
		nonExistentBankID := primitive.NewObjectID()
		transactions, err := service.GetTransactionsByBuyerBankID(ctx, nonExistentBankID)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions for non-existent bank, got %d", len(transactions))
		}

		nonExistentAssetID := primitive.NewObjectID()
		transactions, err = service.GetTransactionsByAssetID(ctx, nonExistentAssetID)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions for non-existent asset, got %d", len(transactions))
		}
	})
}
