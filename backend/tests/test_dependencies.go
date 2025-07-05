package tests

import (
	"context"
	"fmt"
	"time"

	"ponziworld/backend/config"
	"ponziworld/backend/db"
)

// CreateTestDependencies creates handler dependencies for testing
func CreateTestDependencies(testName string) (*config.Container, error) {
	// Create a unique test database name
	timestamp := time.Now().Unix()
	testDatabaseName := fmt.Sprintf("ponziworld_test_%s_%d", testName, timestamp)

	// Connect to database
	client, err := db.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Create dependencies with test database
	container := config.NewContainer(client, testDatabaseName)
	
	// Ensure database indexes are created before running tests
	if err := db.EnsureAllIndexes(container.DatabaseConfig); err != nil {
		return nil, fmt.Errorf("failed to ensure database indexes: %w", err)
	}

	return container, nil
}

// CleanupTestDependencies properly closes test dependencies and cleans up test database
func CleanupTestDependencies(deps *config.Container) {
	if deps != nil {
		// Drop the test database with a fresh context
		ctx := context.Background()
		deps.DatabaseConfig.GetDatabase().Drop(ctx)
		// Close the connection
		deps.Close()
	}
}
