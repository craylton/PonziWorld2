package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ponziworld/backend/models"
	"ponziworld/backend/routes"
)

func TestAssetTypesEndpoint(t *testing.T) {
	container, err := CreateTestDependencies("asset_types")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	// Create test server
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, container)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should return all asset types", func(t *testing.T) {
		// Create a test user and get authentication token
		token, err := CreateRegularUserForTest(container, "testuser", "testpassword123", "Test Bank")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Make authenticated request to get asset types
		req, err := http.NewRequest("GET", server.URL+"/api/assetTypes", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		// Parse response
		var assetTypes []models.AssetType
		if err := json.NewDecoder(resp.Body).Decode(&assetTypes); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Check that we have the expected asset types
		expectedAssetTypes := []string{"Cash", "HYSA", "Bonds", "Stocks", "Crypto"}
		if len(assetTypes) != len(expectedAssetTypes) {
			t.Errorf("Expected %d asset types, got %d", len(expectedAssetTypes), len(assetTypes))
		}

		// Check that all expected asset types are present
		assetTypeNames := make(map[string]bool)
		for _, assetType := range assetTypes {
			assetTypeNames[assetType.Name] = true
		}

		for _, expectedName := range expectedAssetTypes {
			if !assetTypeNames[expectedName] {
				t.Errorf("Expected asset type '%s' not found in response", expectedName)
			}
		}

		// Check that each asset type has a valid ID and name
		for _, assetType := range assetTypes {
			if assetType.Id.IsZero() {
				t.Errorf("Asset type '%s' has empty ID", assetType.Name)
			}
			if assetType.Name == "" {
				t.Errorf("Asset type with ID '%s' has empty name", assetType.Id.Hex())
			}
		}
	})
}
