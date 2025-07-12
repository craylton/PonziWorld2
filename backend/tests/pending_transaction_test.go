package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ponziworld/backend/config"
	"ponziworld/backend/handlers"
	"ponziworld/backend/models"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestDependencies struct {
	Container *config.Container
}

func (d *TestDependencies) Cleanup() {
	CleanupTestDependencies(d.Container)
}

func setupTestDependencies(t *testing.T) *TestDependencies {
	container, err := CreateTestDependencies("pending_transaction_test")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	return &TestDependencies{Container: container}
}

func createTestUser(t *testing.T, deps *TestDependencies) (*models.Player, string) {
	username := "testuser"
	password := "testpass"
	bankName := "Test Bank"

	token, err := CreateRegularUserForTest(deps.Container, username, password, bankName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Retrieve the user to get the ID
	user, err := deps.Container.RepositoryContainer.Player.FindByUsername(context.Background(), username)
	if err != nil {
		t.Fatalf("Failed to find test user: %v", err)
	}

	return user, token
}

func createTestBank(t *testing.T, deps *TestDependencies, playerID primitive.ObjectID) *models.Bank {
	bank, err := deps.Container.RepositoryContainer.Bank.FindByPlayerID(context.Background(), playerID)
	if err != nil {
		t.Fatalf("Failed to find test bank: %v", err)
	}
	return bank
}

func TestBuyAsset(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	// Create a test user and bank
	user, _ := createTestUser(t, deps)
	bank := createTestBank(t, deps, user.Id)

	// Create a test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err := deps.Container.RepositoryContainer.AssetType.Create(context.Background(), assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	// Create buy request
	buyRequest := models.PendingTransactionRequest{
		BuyerBankId: bank.Id.Hex(),
		AssetId:     assetType.Id.Hex(),
		Amount:      1000, // Positive amount = buy
	}

	requestBody, _ := json.Marshal(buyRequest)

	// Create test request
	req := httptest.NewRequest(http.MethodPost, "/api/buy", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Username", user.Username) // Simulate JWT middleware

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create handler
	handler := handlers.NewPendingTransactionHandler(deps.Container)
	handler.BuyAsset(rr, req)

	// Check response
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		t.Logf("Response body: %s", rr.Body.String())
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "buy transaction created successfully" {
		t.Errorf("Expected success message, got: %s", response["message"])
	}

	// Verify transaction was created in database
	transactions, err := deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank.Id)
	if err != nil {
		t.Fatalf("Failed to find pending transactions: %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("Expected 1 pending transaction, got %d", len(transactions))
	}

	if transactions[0].Amount != 1000 {
		t.Errorf("Expected amount 1000, got %d", transactions[0].Amount)
	}

	if transactions[0].Amount <= 0 {
		t.Errorf("Expected positive amount for buy transaction, got %d", transactions[0].Amount)
	}
}

func TestSelfInvestmentPrevention(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	// Create a test user and bank
	user, _ := createTestUser(t, deps)
	bank := createTestBank(t, deps, user.Id)

	// Create request where bank tries to invest in itself
	buyRequest := models.PendingTransactionRequest{
		BuyerBankId: bank.Id.Hex(),
		AssetId:     bank.Id.Hex(), // Same as buyer bank ID
		Amount:      1000,
	}

	requestBody, _ := json.Marshal(buyRequest)

	// Create test request
	req := httptest.NewRequest(http.MethodPost, "/api/buy", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Username", user.Username) // Simulate JWT middleware

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create handler
	handler := handlers.NewPendingTransactionHandler(deps.Container)
	handler.BuyAsset(rr, req)

	// Check response
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		t.Logf("Response body: %s", rr.Body.String())
	}

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Bank cannot invest in itself" {
		t.Errorf("Expected self-investment error, got: %s", response["error"])
	}
}

func TestSellAssetTransaction(t *testing.T) {
	deps := setupTestDependencies(t)
	defer deps.Cleanup()

	// Create a test user and bank
	user, _ := createTestUser(t, deps)
	bank := createTestBank(t, deps, user.Id)

	// Create a test asset type
	assetType := &models.AssetType{
		Id:   primitive.NewObjectID(),
		Name: "Test Asset",
	}
	err := deps.Container.RepositoryContainer.AssetType.Create(context.Background(), assetType)
	if err != nil {
		t.Fatalf("Failed to create test asset type: %v", err)
	}

	// Create sell request (negative amount)
	sellRequest := models.PendingTransactionRequest{
		BuyerBankId: bank.Id.Hex(),
		AssetId:     assetType.Id.Hex(),
		Amount:      -500, // Negative amount = sell
	}

	requestBody, _ := json.Marshal(sellRequest)

	// Create test request
	req := httptest.NewRequest(http.MethodPost, "/api/sell", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Username", user.Username) // Simulate JWT middleware

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create handler
	handler := handlers.NewPendingTransactionHandler(deps.Container)
	handler.SellAsset(rr, req)

	// Check response
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		t.Logf("Response body: %s", rr.Body.String())
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "sell transaction created successfully" {
		t.Errorf("Expected success message, got: %s", response["message"])
	}

	// Verify transaction was created in database
	transactions, err := deps.Container.RepositoryContainer.PendingTransaction.FindByBuyerBankID(context.Background(), bank.Id)
	if err != nil {
		t.Fatalf("Failed to find pending transactions: %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("Expected 1 pending transaction, got %d", len(transactions))
	}

	if transactions[0].Amount != -500 {
		t.Errorf("Expected amount -500, got %d", transactions[0].Amount)
	}

	if transactions[0].Amount >= 0 {
		t.Errorf("Expected negative amount for sell transaction, got %d", transactions[0].Amount)
	}
}
