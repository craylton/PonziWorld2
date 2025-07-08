package tests

import (
	"context"
	"testing"
)

func TestAssetTypeService_GetAllAssetTypes(t *testing.T) {
	container, err := CreateTestDependencies("asset_types")
	if err != nil {
		t.Fatalf("Failed to create test dependencies: %v", err)
	}
	defer CleanupTestDependencies(container)

	ctx := context.Background()
	assetTypes, err := container.ServiceContainer.AssetType.GetAllAssetTypes(ctx)
	if err != nil {
		t.Fatalf("Failed to get asset types: %v", err)
	}

	expected := []string{"Cash", "HYSA", "Bonds", "Stocks", "Crypto"}
	if len(assetTypes) != len(expected) {
		t.Errorf("Expected %d asset types, got %d", len(expected), len(assetTypes))
	}

	names := make(map[string]bool)
	for _, at := range assetTypes {
		names[at.Name] = true
		if at.Id.IsZero() {
			t.Errorf("Asset type '%s' has empty ID", at.Name)
		}
		if at.Name == "" {
			t.Errorf("Asset type with ID '%s' has empty name", at.Id.Hex())
		}
	}

	for _, name := range expected {
		if !names[name] {
			t.Errorf("Expected asset type '%s' not found", name)
		}
	}
}
