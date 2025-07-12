package tests

import (
	"context"
	"ponziworld/backend/models"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Unit tests for PendingTransactionRepository (repository layer only)
func TestPendingTransactionRepository_FindByBuyerBankIDAndAssetID(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	// Create test data
	buyerBankID := primitive.NewObjectID()
	assetID1 := primitive.NewObjectID()
	assetID2 := primitive.NewObjectID()

	// Create transactions
	transaction1 := &models.PendingTransaction{
		BuyerBankId: buyerBankID,
		AssetId:     assetID1,
		Amount:      1000,
	}
	transaction2 := &models.PendingTransaction{
		BuyerBankId: buyerBankID,
		AssetId:     assetID2,
		Amount:      2000,
	}
	transaction3 := &models.PendingTransaction{
		BuyerBankId: primitive.NewObjectID(), // Different buyer
		AssetId:     assetID1,
		Amount:      3000,
	}

	repo := deps.Container.RepositoryContainer.PendingTransaction
	ctx := context.Background()

	// Insert transactions
	err := repo.Create(ctx, transaction1)
	if err != nil {
		t.Fatalf("Failed to create transaction 1: %v", err)
	}
	err = repo.Create(ctx, transaction2)
	if err != nil {
		t.Fatalf("Failed to create transaction 2: %v", err)
	}
	err = repo.Create(ctx, transaction3)
	if err != nil {
		t.Fatalf("Failed to create transaction 3: %v", err)
	}

	// Test 1: Find by exact bank and asset combination
	results, err := repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, assetID1)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 transaction for buyer+asset1, got %d", len(results))
	}
	if len(results) > 0 && results[0].Amount != 1000 {
		t.Errorf("Expected amount 1000, got %d", results[0].Amount)
	}

	// Test 2: Find by different asset for same buyer
	results, err = repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, assetID2)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 transaction for buyer+asset2, got %d", len(results))
	}
	if len(results) > 0 && results[0].Amount != 2000 {
		t.Errorf("Expected amount 2000, got %d", results[0].Amount)
	}

	// Test 3: Find non-existent combination
	nonExistentAssetID := primitive.NewObjectID()
	results, err = repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, nonExistentAssetID)
	if err != nil {
		t.Fatalf("Failed to find transactions: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 transactions for non-existent combination, got %d", len(results))
	}
}

func TestPendingTransactionRepository_UpdateAmount(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	// Create a transaction
	transaction := &models.PendingTransaction{
		BuyerBankId: primitive.NewObjectID(),
		AssetId:     primitive.NewObjectID(),
		Amount:      1000,
	}

	repo := deps.Container.RepositoryContainer.PendingTransaction
	ctx := context.Background()

	// Insert transaction
	err := repo.Create(ctx, transaction)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Test 1: Update to positive amount
	err = repo.UpdateAmount(ctx, transaction.Id, 1500)
	if err != nil {
		t.Fatalf("Failed to update amount: %v", err)
	}

	// Verify update
	transactions, err := repo.FindByBuyerBankID(ctx, transaction.BuyerBankId)
	if err != nil {
		t.Fatalf("Failed to find updated transaction: %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(transactions))
	}
	if transactions[0].Amount != 1500 {
		t.Errorf("Expected updated amount 1500, got %d", transactions[0].Amount)
	}

	// Test 2: Update to negative amount
	err = repo.UpdateAmount(ctx, transaction.Id, -500)
	if err != nil {
		t.Fatalf("Failed to update to negative amount: %v", err)
	}

	// Verify negative update
	transactions, err = repo.FindByBuyerBankID(ctx, transaction.BuyerBankId)
	if err != nil {
		t.Fatalf("Failed to find updated transaction: %v", err)
	}
	if transactions[0].Amount != -500 {
		t.Errorf("Expected updated amount -500, got %d", transactions[0].Amount)
	}

	// Test 3: Update to zero
	err = repo.UpdateAmount(ctx, transaction.Id, 0)
	if err != nil {
		t.Fatalf("Failed to update to zero: %v", err)
	}

	// Verify zero update
	transactions, err = repo.FindByBuyerBankID(ctx, transaction.BuyerBankId)
	if err != nil {
		t.Fatalf("Failed to find updated transaction: %v", err)
	}
	if transactions[0].Amount != 0 {
		t.Errorf("Expected updated amount 0, got %d", transactions[0].Amount)
	}

	// Test 4: Update non-existent transaction
	nonExistentID := primitive.NewObjectID()
	err = repo.UpdateAmount(ctx, nonExistentID, 1000)
	if err == nil {
		t.Error("Expected error when updating non-existent transaction")
	}
}

func TestPendingTransactionRepository_CompleteWorkflow(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	repo := deps.Container.RepositoryContainer.PendingTransaction
	ctx := context.Background()

	buyerBankID := primitive.NewObjectID()
	assetID := primitive.NewObjectID()

	// Step 1: Create initial transaction
	transaction := &models.PendingTransaction{
		BuyerBankId: buyerBankID,
		AssetId:     assetID,
		Amount:      1000,
	}
	err := repo.Create(ctx, transaction)
	if err != nil {
		t.Fatalf("Failed to create initial transaction: %v", err)
	}

	// Step 2: Find and verify
	results, err := repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, assetID)
	if err != nil {
		t.Fatalf("Failed to find transaction: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(results))
	}
	transactionID := results[0].Id

	// Step 3: Update amount multiple times
	err = repo.UpdateAmount(ctx, transactionID, 1500)
	if err != nil {
		t.Fatalf("Failed to update amount to 1500: %v", err)
	}

	err = repo.UpdateAmount(ctx, transactionID, 800)
	if err != nil {
		t.Fatalf("Failed to update amount to 800: %v", err)
	}

	// Verify final amount
	results, err = repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, assetID)
	if err != nil {
		t.Fatalf("Failed to find updated transaction: %v", err)
	}
	if results[0].Amount != 800 {
		t.Errorf("Expected final amount 800, got %d", results[0].Amount)
	}

	// Step 4: Delete transaction
	err = repo.Delete(ctx, transactionID)
	if err != nil {
		t.Fatalf("Failed to delete transaction: %v", err)
	}

	// Step 5: Verify deletion
	results, err = repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, assetID)
	if err != nil {
		t.Fatalf("Failed to find after deletion: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 transactions after deletion, got %d", len(results))
	}
}

func TestPendingTransactionRepository_ConcurrentOperations(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	repo := deps.Container.RepositoryContainer.PendingTransaction
	ctx := context.Background()

	buyerBankID := primitive.NewObjectID()
	
	// Create multiple transactions with same buyer but different assets
	numTransactions := 10
	assetIDs := make([]primitive.ObjectID, numTransactions)
	
	for i := 0; i < numTransactions; i++ {
		assetIDs[i] = primitive.NewObjectID()
		transaction := &models.PendingTransaction{
			BuyerBankId: buyerBankID,
			AssetId:     assetIDs[i],
			Amount:      int64(1000 + i*100),
		}
		err := repo.Create(ctx, transaction)
		if err != nil {
			t.Fatalf("Failed to create transaction %d: %v", i, err)
		}
	}

	// Verify all transactions exist
	allTransactions, err := repo.FindByBuyerBankID(ctx, buyerBankID)
	if err != nil {
		t.Fatalf("Failed to find all transactions: %v", err)
	}
	if len(allTransactions) != numTransactions {
		t.Errorf("Expected %d transactions, got %d", numTransactions, len(allTransactions))
	}

	// Test finding by specific asset
	for i, assetID := range assetIDs {
		results, err := repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, assetID)
		if err != nil {
			t.Fatalf("Failed to find transaction for asset %d: %v", i, err)
		}
		if len(results) != 1 {
			t.Errorf("Expected 1 transaction for asset %d, got %d", i, len(results))
		}
		if len(results) > 0 && results[0].Amount != int64(1000+i*100) {
			t.Errorf("Expected amount %d for asset %d, got %d", 1000+i*100, i, results[0].Amount)
		}
	}

	// Update amounts for all transactions
	for i, assetID := range assetIDs {
		results, _ := repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, assetID)
		if len(results) > 0 {
			newAmount := int64(2000 + i*200)
			err := repo.UpdateAmount(ctx, results[0].Id, newAmount)
			if err != nil {
				t.Fatalf("Failed to update amount for transaction %d: %v", i, err)
			}
		}
	}

	// Verify all updates
	for i, assetID := range assetIDs {
		results, err := repo.FindByBuyerBankIDAndAssetID(ctx, buyerBankID, assetID)
		if err != nil {
			t.Fatalf("Failed to find updated transaction for asset %d: %v", i, err)
		}
		if len(results) > 0 && results[0].Amount != int64(2000+i*200) {
			t.Errorf("Expected updated amount %d for asset %d, got %d", 2000+i*200, i, results[0].Amount)
		}
	}
}
