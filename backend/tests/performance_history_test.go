package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ponziworld/backend/models"
	"ponziworld/backend/routes"
)

func TestPerformanceHistoryEndpoint(t *testing.T) {
	// Reset game state to ensure consistent test environment
	ResetGameState()

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("perftest_%d", timestamp)
	testBankName := "Test Bank Performance"

	// Setup cleanup
	t.Cleanup(func() {
		CleanupTestData(testUsername, testBankName)
	})

	// Create player and bank
	createUserData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
		"bankName": testBankName,
	}
	jsonData, _ := json.Marshal(createUserData)

	resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 for player creation, got %d", resp.StatusCode)
	}

	// Login to get JWT
	loginData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
	}
	jsonData, _ = json.Marshal(loginData)

	resp, err = http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for login, got %d", resp.StatusCode)
	}

	var loginResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	token := loginResponse["token"]
	if token == "" {
		t.Fatal("No token received from login")
	}

	// Get bank details to get bank ID
	req, err := http.NewRequest("GET", server.URL+"/api/bank", nil)
	if err != nil {
		t.Fatalf("Failed to create bank request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get bank: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for bank retrieval, got %d", resp.StatusCode)
	}

	var bankResponse models.BankResponse
	if err := json.NewDecoder(resp.Body).Decode(&bankResponse); err != nil {
		t.Fatalf("Failed to decode bank response: %v", err)
	}

	bankId := bankResponse.Id
	if bankId == "" {
		t.Fatal("No bank ID received")
	}

	// Test performance history endpoint
	req, err = http.NewRequest("GET", server.URL+"/api/performanceHistory/ownbank/"+bankId, nil)
	if err != nil {
		t.Fatalf("Failed to create performance history request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get performance history: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for performance history, got %d", resp.StatusCode)
	}

	var historyResponse models.PerformanceHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&historyResponse); err != nil {
		t.Fatalf("Failed to decode performance history response: %v", err)
	}

	// Verify that we get 30 days of history
	if len(historyResponse.ClaimedHistory) != 30 {
		t.Fatalf("Expected 30 days of claimed history, got %d", len(historyResponse.ClaimedHistory))
	}

	// Since actual history should contain the initial entry for a new bank, we expect 1 entry
	// For a newly created bank, this should be 1 (the initial £1000 entry)
	if len(historyResponse.ActualHistory) != 1 {
		t.Fatalf("Expected 1 day of actual history for new bank, got %d", len(historyResponse.ActualHistory))
	}

	// Verify the initial actual history entry is correct
	if len(historyResponse.ActualHistory) > 0 {
		initialEntry := historyResponse.ActualHistory[0]
		if initialEntry.Day != 0 {
			t.Fatalf("Expected initial actual history entry to be for day 0, got day %d", initialEntry.Day)
		}
		if initialEntry.Value != 1000 {
			t.Fatalf("Expected initial actual history entry to have value 1000, got %d", initialEntry.Value)
		}
	}

	// Verify that all entries are properly ordered by day
	for i := 1; i < len(historyResponse.ClaimedHistory); i++ {
		if historyResponse.ClaimedHistory[i].Day <= historyResponse.ClaimedHistory[i-1].Day {
			t.Fatal("Claimed history is not properly ordered by day")
		}
	}

	// For actual history, only verify if we have entries
	for i := 1; i < len(historyResponse.ActualHistory); i++ {
		if historyResponse.ActualHistory[i].Day <= historyResponse.ActualHistory[i-1].Day {
			t.Fatal("Actual history is not properly ordered by day")
		}
	}

	// Verify that claimed and actual history have the same values for the days that exist
	// Since actual history might not exist for all days, we can't assume they're the same length
	for i := 0; i < len(historyResponse.ActualHistory); i++ {
		found := false
		for j := 0; j < len(historyResponse.ClaimedHistory); j++ {
			if historyResponse.ClaimedHistory[j].Day == historyResponse.ActualHistory[i].Day {
				if historyResponse.ClaimedHistory[j].Value != historyResponse.ActualHistory[i].Value {
					t.Fatalf("Claimed and actual history values don't match for day %d", historyResponse.ActualHistory[i].Day)
				}
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Actual history day %d not found in claimed history", historyResponse.ActualHistory[i].Day)
		}
	}
}

func TestPerformanceHistoryUnauthorized(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Test without authentication
	req, err := http.NewRequest("GET", server.URL+"/api/performanceHistory/ownbank/507f1f77bcf86cd799439011", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 for unauthorized request, got %d", resp.StatusCode)
	}
}

func TestPerformanceHistoryInvalidBankID(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("perfinvalidtest_%d", timestamp)
	testBankName := "Test Bank Invalid"

	// Setup cleanup
	t.Cleanup(func() {
		CleanupTestData(testUsername, testBankName)
	})

	// Create player and get token
	createUserData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
		"bankName": testBankName,
	}
	jsonData, _ := json.Marshal(createUserData)

	resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	defer resp.Body.Close()

	loginData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
	}
	jsonData, _ = json.Marshal(loginData)

	resp, err = http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	var loginResponse map[string]string
	json.NewDecoder(resp.Body).Decode(&loginResponse)
	token := loginResponse["token"]

	// Test with invalid bank ID
	req, err := http.NewRequest("GET", server.URL+"/api/performanceHistory/ownbank/invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid bank ID, got %d", resp.StatusCode)
	}
}

func TestPerformanceHistoryOtherPlayersBank(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	timestamp := time.Now().Unix()
	player1Username := fmt.Sprintf("perfplayer1_%d", timestamp)
	player1BankName := "Player 1 Bank"
	player2Username := fmt.Sprintf("perfplayer2_%d", timestamp)
	player2BankName := "Player 2 Bank"

	// Setup cleanup
	t.Cleanup(func() {
		CleanupMultipleTestData(map[string]string{
			player1Username: player1BankName,
			player2Username: player2BankName,
		})
	})

	// Create first player and bank
	createPlayer1Data := map[string]string{
		"username": player1Username,
		"password": "testpassword123",
		"bankName": player1BankName,
	}
	jsonData, _ := json.Marshal(createPlayer1Data)

	resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create player 1: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 for player 1 creation, got %d", resp.StatusCode)
	}

	// Create second player and bank
	createPlayer2Data := map[string]string{
		"username": player2Username,
		"password": "testpassword123",
		"bankName": player2BankName,
	}
	jsonData, _ = json.Marshal(createPlayer2Data)

	resp, err = http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create player 2: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 for player 2 creation, got %d", resp.StatusCode)
	}

	// Login as player 1 to get token
	loginData := map[string]string{
		"username": player1Username,
		"password": "testpassword123",
	}
	jsonData, _ = json.Marshal(loginData)

	resp, err = http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login player 1: %v", err)
	}
	defer resp.Body.Close()

	var loginResponse map[string]string
	json.NewDecoder(resp.Body).Decode(&loginResponse)
	player1Token := loginResponse["token"]

	// Login as player 2 to get their bank ID
	loginData = map[string]string{
		"username": player2Username,
		"password": "testpassword123",
	}
	jsonData, _ = json.Marshal(loginData)

	resp, err = http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login player 2: %v", err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&loginResponse)
	player2Token := loginResponse["token"]

	// Get player 2's bank details to get bank ID
	req, err := http.NewRequest("GET", server.URL+"/api/bank", nil)
	if err != nil {
		t.Fatalf("Failed to create bank request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+player2Token)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get player 2's bank: %v", err)
	}
	defer resp.Body.Close()

	var bankResponse models.BankResponse
	json.NewDecoder(resp.Body).Decode(&bankResponse)
	player2BankId := bankResponse.Id

	// Now, as player 1, try to access player 2's bank performance history
	req, err = http.NewRequest("GET", server.URL+"/api/performanceHistory/ownbank/"+player2BankId, nil)
	if err != nil {
		t.Fatalf("Failed to create performance history request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+player1Token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get performance history: %v", err)
	}
	defer resp.Body.Close()

	// Should now return 401 Unauthorized since player 1 doesn't own player 2's bank
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 Unauthorized for other player's bank, got %d", resp.StatusCode)
	}
}

func TestPerformanceHistoryDataPersistence(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("perfpersist_%d", timestamp)
	testBankName := "Test Bank Persistence"

	// Setup cleanup
	t.Cleanup(func() {
		CleanupTestData(testUsername, testBankName)
	})

	// Create player and bank
	createUserData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
		"bankName": testBankName,
	}
	jsonData, _ := json.Marshal(createUserData)

	resp, err := http.Post(server.URL+"/api/newPlayer", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	defer resp.Body.Close()

	// Login and get bank ID
	loginData := map[string]string{
		"username": testUsername,
		"password": "testpassword123",
	}
	jsonData, _ = json.Marshal(loginData)

	resp, err = http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	var loginResponse map[string]string
	json.NewDecoder(resp.Body).Decode(&loginResponse)
	token := loginResponse["token"]

	req, err := http.NewRequest("GET", server.URL+"/api/bank", nil)
	if err != nil {
		t.Fatalf("Failed to create bank request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get bank: %v", err)
	}
	defer resp.Body.Close()

	var bankResponse models.BankResponse
	json.NewDecoder(resp.Body).Decode(&bankResponse)
	bankId := bankResponse.Id

	// First call to performance history endpoint
	req, err = http.NewRequest("GET", server.URL+"/api/performanceHistory/ownbank/"+bankId, nil)
	if err != nil {
		t.Fatalf("Failed to create performance history request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get performance history: %v", err)
	}
	defer resp.Body.Close()

	var firstResponse models.PerformanceHistoryResponse
	json.NewDecoder(resp.Body).Decode(&firstResponse)

	// Second call to performance history endpoint (should return identical data)
	req, err = http.NewRequest("GET", server.URL+"/api/performanceHistory/ownbank/"+bankId, nil)
	if err != nil {
		t.Fatalf("Failed to create second performance history request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get performance history second time: %v", err)
	}
	defer resp.Body.Close()

	var secondResponse models.PerformanceHistoryResponse
	json.NewDecoder(resp.Body).Decode(&secondResponse)

	// Verify that both responses are identical (data persisted in database)
	if len(firstResponse.ClaimedHistory) != len(secondResponse.ClaimedHistory) {
		t.Fatalf("Claimed history length differs between calls: %d vs %d",
			len(firstResponse.ClaimedHistory), len(secondResponse.ClaimedHistory))
	}

	if len(firstResponse.ActualHistory) != len(secondResponse.ActualHistory) {
		t.Fatalf("Actual history length differs between calls: %d vs %d",
			len(firstResponse.ActualHistory), len(secondResponse.ActualHistory))
	}

	for i := 0; i < len(firstResponse.ClaimedHistory); i++ {
		if firstResponse.ClaimedHistory[i].Day != secondResponse.ClaimedHistory[i].Day ||
			firstResponse.ClaimedHistory[i].Value != secondResponse.ClaimedHistory[i].Value {
			t.Fatalf("Claimed history differs at index %d: first=(%d,%d), second=(%d,%d)",
				i, firstResponse.ClaimedHistory[i].Day, firstResponse.ClaimedHistory[i].Value,
				secondResponse.ClaimedHistory[i].Day, secondResponse.ClaimedHistory[i].Value)
		}
	}

	for i := 0; i < len(firstResponse.ActualHistory); i++ {
		if firstResponse.ActualHistory[i].Day != secondResponse.ActualHistory[i].Day ||
			firstResponse.ActualHistory[i].Value != secondResponse.ActualHistory[i].Value {
			t.Fatalf("Actual history differs at index %d: first=(%d,%d), second=(%d,%d)",
				i, firstResponse.ActualHistory[i].Day, firstResponse.ActualHistory[i].Value,
				secondResponse.ActualHistory[i].Day, secondResponse.ActualHistory[i].Value)
		}
	}
}
