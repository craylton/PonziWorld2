package tests

import (
	"context"
	"fmt"
	"time"

	"ponziworld/backend/config"
	"ponziworld/backend/database"
	"ponziworld/backend/logging"
)

// CreateTestDependencies creates handler dependencies for testing
func CreateTestDependencies(testName string) (*config.Container, error) {
	// Create a unique test database name
	timestamp := time.Now().Unix()
	testDatabaseName := fmt.Sprintf("ponziworld_test_%s_%d", testName, timestamp)

	client, err := database.InitializeDatabaseConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Create a logger for testing
	logger := logging.NewLogger()

	container := config.NewContainer(client, testDatabaseName, logger)

	err = database.EnsureDatabaseStructure(container.DatabaseConfig, container.ServiceContainer.AssetType)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure database structure: %w", err)
	}

	return container, nil
}

// CleanupTestDependencies properly closes test dependencies and cleans up test database
func CleanupTestDependencies(container *config.Container) {
	if container != nil {
		// Drop the test database with a fresh context
		ctx := context.Background()
		container.DatabaseConfig.GetDatabase().Drop(ctx)
		// Close the connection
		container.Close()
	}
}
