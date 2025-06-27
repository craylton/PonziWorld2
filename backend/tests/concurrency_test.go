package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"ponziworld/backend/routes"
)

func TestConcurrentUserCreation(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("Multiple users created simultaneously", func(t *testing.T) {
		const numUsers = 10
		var wg sync.WaitGroup
		results := make(chan int, numUsers)
		timestamp := time.Now().Unix()

		for i := range numUsers {
			wg.Add(1)
			go func(userNum int) {
				defer wg.Done()
				
				createUserData := map[string]string{
					"username": fmt.Sprintf("concurrent_%d_%d", timestamp, userNum),
					"password": "testpassword123",
					"bankName": fmt.Sprintf("Bank %d", userNum),
				}
				jsonData, _ := json.Marshal(createUserData)

				resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					results <- 500
					return
				}
				resp.Body.Close()
				results <- resp.StatusCode
			}(i)
		}

		wg.Wait()
		close(results)

		successCount := 0
		for statusCode := range results {
			if statusCode == http.StatusCreated {
				successCount++
			}
		}

		if successCount != numUsers {
			t.Errorf("Expected %d successful user creations, got %d", numUsers, successCount)
		}

		// Cleanup
		usersAndBanks := make(map[string]string)
		for i := 0; i < numUsers; i++ {
			username := fmt.Sprintf("concurrent_%d_%d", timestamp, i)
			bankName := fmt.Sprintf("Bank %d", i)
			usersAndBanks[username] = bankName
		}
		CleanupMultipleTestData(usersAndBanks)
	})

	t.Run("Duplicate username creation race condition", func(t *testing.T) {
		const numAttempts = 5
		var wg sync.WaitGroup
		results := make(chan int, numAttempts)
		timestamp := time.Now().Unix()
		duplicateUsername := fmt.Sprintf("racetest_%d", timestamp)

		// Try to create the same username multiple times simultaneously
		for i := 0; i < numAttempts; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				createUserData := map[string]string{
					"username": duplicateUsername,
					"password": "testpassword123",
					"bankName": "Race Test Bank",
				}
				jsonData, _ := json.Marshal(createUserData)

				resp, err := http.Post(server.URL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					results <- 500
					return
				}
				resp.Body.Close()
				results <- resp.StatusCode
			}()
		}

		wg.Wait()
		close(results)

		successCount := 0
		errorCount := 0
		for statusCode := range results {
			if statusCode == http.StatusCreated {
				successCount++
			} else if statusCode == http.StatusBadRequest {
				errorCount++
			}
		}

		// Only one should succeed, the rest should fail with duplicate username error
		if successCount != 1 {
			t.Errorf("Expected exactly 1 successful creation, got %d", successCount)
		}
		if errorCount != numAttempts-1 {
			t.Errorf("Expected %d duplicate username errors, got %d", numAttempts-1, errorCount)
		}

		// Cleanup
		CleanupTestData(duplicateUsername, "Race Test Bank")
	})
}
