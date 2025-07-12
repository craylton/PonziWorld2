package tests

import (
	"context"
	"ponziworld/backend/models"
	"ponziworld/backend/services"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Unit tests for PendingTransactionService (service layer only)
func TestPendingTransactionService_CreateTransaction_OwnershipValidation(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	// Create two users
	user1, _ := createTestUser(t, deps)
	user2Username := "testuser2"
	user2Password := "testpass2"
	user2BankName := "Test Bank 2"
	
	_, err := CreateRegularUserForTest(deps.Container, user2Username, user2Password, user2BankName)
	if err != nil {
		t.Fatalf("Failed to create second test user: %v", err)
	}

	user2, err := deps.Container.RepositoryContainer.Player.FindByUsername(context.Background(), user2Username)
	if err != nil {
		t.Fatalf("Failed to find second test user: %v", err)
	}

	// Get banks for both users
	bank1 := createTestBank(t, deps, user1.Id)
	bank2 := createTestBank(t, deps, user2.Id)

	// Create an asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err = deps.Container.RepositoryContainer.AssetType.Create(context.Background(), assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	// Test 1: User1 tries to create transaction with their own bank (should succeed)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank1.Id,
		assetType.Id,
		1000,
		user1.Username,
	)
	if err != nil {
		t.Errorf("Expected no error when user uses their own bank, got: %v", err)
	}

	// Test 2: User1 tries to create transaction with User2's bank (should fail)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank2.Id,
		assetType.Id,
		1000,
		user1.Username,
	)
	if err != services.ErrUnauthorizedBank {
		t.Errorf("Expected ErrUnauthorizedBank when user tries to use another user's bank, got: %v", err)
	}

	// Test 3: Non-existent user tries to create transaction (should fail)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank1.Id,
		assetType.Id,
		1000,
		"nonexistentuser",
	)
	if err != services.ErrInvalidBankID {
		t.Errorf("Expected ErrInvalidBankID when user doesn't exist, got: %v", err)
	}
}

func TestPendingTransactionService_CreateTransaction_TransactionCombining(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	user, _ := createTestUser(t, deps)
	bank := createTestBank(t, deps, user.Id)

	// Create an asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err := deps.Container.RepositoryContainer.AssetType.Create(context.Background(), assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	// Test 1: Create initial buy transaction
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType.Id,
		1000,
		user.Username,
	)
	if err != nil {
		t.Fatalf("Failed to create initial transaction: %v", err)
	}

	// Verify initial transaction exists
	transactions, err := deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank.Id)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(transactions))
	}
	if transactions[0].Amount != 1000 {
		t.Errorf("Expected amount 1000, got %d", transactions[0].Amount)
	}

	// Test 2: Add another buy transaction (should combine)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType.Id,
		500,
		user.Username,
	)
	if err != nil {
		t.Fatalf("Failed to create second transaction: %v", err)
	}

	// Verify transactions are combined
	transactions, err = deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank.Id)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
	}
	if len(transactions) != 1 {
		t.Errorf("Expected 1 combined transaction, got %d", len(transactions))
	}
	if transactions[0].Amount != 1500 {
		t.Errorf("Expected combined amount 1500, got %d", transactions[0].Amount)
	}

	// Test 3: Add sell transaction (should reduce amount)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType.Id,
		-800,
		user.Username,
	)
	if err != nil {
		t.Fatalf("Failed to create sell transaction: %v", err)
	}

	// Verify amount is reduced
	transactions, err = deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank.Id)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
	}
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction after sell, got %d", len(transactions))
	}
	if transactions[0].Amount != 700 {
		t.Errorf("Expected reduced amount 700, got %d", transactions[0].Amount)
	}

	// Test 4: Sell exact remaining amount (should delete transaction)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType.Id,
		-700,
		user.Username,
	)
	if err != nil {
		t.Fatalf("Failed to create final sell transaction: %v", err)
	}

	// Verify transaction is deleted
	transactions, err = deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank.Id)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
	}
	if len(transactions) != 0 {
		t.Errorf("Expected 0 transactions after selling all, got %d", len(transactions))
	}
}

func TestPendingTransactionService_CreateTransaction_EdgeCases(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	user, _ := createTestUser(t, deps)
	bank := createTestBank(t, deps, user.Id)

	// Create an asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err := deps.Container.RepositoryContainer.AssetType.Create(context.Background(), assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	// Test 1: Zero amount (should fail)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType.Id,
		0,
		user.Username,
	)
	if err != services.ErrInvalidAmount {
		t.Errorf("Expected ErrInvalidAmount for zero amount, got: %v", err)
	}

	// Test 2: Non-existent bank (should fail)
	nonExistentBankID := primitive.NewObjectID()
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		nonExistentBankID,
		assetType.Id,
		1000,
		user.Username,
	)
	if err != services.ErrInvalidBankID {
		t.Errorf("Expected ErrInvalidBankID for non-existent bank, got: %v", err)
	}

	// Test 3: Non-existent asset (should fail)
	nonExistentAssetID := primitive.NewObjectID()
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		nonExistentAssetID,
		1000,
		user.Username,
	)
	if err != services.ErrAssetNotFound {
		t.Errorf("Expected ErrAssetNotFound for non-existent asset, got: %v", err)
	}

	// Test 4: Self-investment (should fail)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		bank.Id,
		1000,
		user.Username,
	)
	if err != services.ErrSelfInvestment {
		t.Errorf("Expected ErrSelfInvestment for self-investment, got: %v", err)
	}

	// Test 5: Very large amounts (should work)
	largeAmount := int64(9999999999)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType.Id,
		largeAmount,
		user.Username,
	)
	if err != nil {
		t.Errorf("Expected no error for large amount, got: %v", err)
	}

	// Test 6: Very large negative amounts (should work)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType.Id,
		-largeAmount,
		user.Username,
	)
	if err != nil {
		t.Errorf("Expected no error for large negative amount, got: %v", err)
	}
}

func TestPendingTransactionService_CreateTransaction_BankAsAsset(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	// Create two users with banks
	user1, _ := createTestUser(t, deps)
	user2Username := "testuser2"
	user2Password := "testpass2"
	user2BankName := "Test Bank 2"
	
	_, err := CreateRegularUserForTest(deps.Container, user2Username, user2Password, user2BankName)
	if err != nil {
		t.Fatalf("Failed to create second test user: %v", err)
	}

	user2, err := deps.Container.RepositoryContainer.Player.FindByUsername(context.Background(), user2Username)
	if err != nil {
		t.Fatalf("Failed to find second test user: %v", err)
	}

	bank1 := createTestBank(t, deps, user1.Id)
	bank2 := createTestBank(t, deps, user2.Id)

	// Test 1: User1 invests in User2's bank (should work since banks are also assets)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank1.Id,
		bank2.Id,
		1000,
		user1.Username,
	)
	if err != nil {
		t.Errorf("Expected no error when investing in another bank, got: %v", err)
	}

	// Verify transaction was created
	transactions, err := deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank1.Id)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
	}
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}
	if transactions[0].AssetId != bank2.Id {
		t.Errorf("Expected AssetId to be bank2 ID, got different ID")
	}
}

func TestPendingTransactionService_CreateTransaction_MultipleAssets(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	user, _ := createTestUser(t, deps)
	bank := createTestBank(t, deps, user.Id)

	// Create multiple asset types
	assetType1 := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Asset Type 1",
	}
	assetType2 := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Asset Type 2",
	}

	err := deps.Container.RepositoryContainer.AssetType.Create(context.Background(), assetType1)
	if err != nil {
		t.Fatalf("Failed to create asset type 1: %v", err)
	}
	err = deps.Container.RepositoryContainer.AssetType.Create(context.Background(), assetType2)
	if err != nil {
		t.Fatalf("Failed to create asset type 2: %v", err)
	}

	// Create transactions for different assets
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType1.Id,
		1000,
		user.Username,
	)
	if err != nil {
		t.Fatalf("Failed to create transaction for asset 1: %v", err)
	}

	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType2.Id,
		2000,
		user.Username,
	)
	if err != nil {
		t.Fatalf("Failed to create transaction for asset 2: %v", err)
	}

	// Verify both transactions exist separately
	transactions, err := deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank.Id)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
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

	// Add more to asset 1 (should combine)
	err = deps.Container.ServiceContainer.PendingTransaction.CreateTransaction(
		context.Background(),
		bank.Id,
		assetType1.Id,
		500,
		user.Username,
	)
	if err != nil {
		t.Fatalf("Failed to add to transaction for asset 1: %v", err)
	}

	// Verify asset 1 amount is combined but asset 2 remains unchanged
	transactions, err = deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank.Id)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
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
}
