package tests

import (
	"context"
	"fmt"
	"time"

	"ponziworld/backend/config"
	"ponziworld/backend/db"
)

// CreateTestDependencies creates handler dependencies for testing
func CreateTestDependencies(testName string) *config.HandlerDependencies {
	// Create a unique test database name
	timestamp := time.Now().Unix()
	testDatabaseName := fmt.Sprintf("ponziworld_test_%s_%d", testName, timestamp)
	
	// Connect to database
	client, ctx, cancel := db.ConnectDB()
	
	// Create dependencies with test database
	return config.NewHandlerDependencies(client, ctx, cancel, testDatabaseName)
}

// CleanupTestDependencies properly closes test dependencies and cleans up test database
func CleanupTestDependencies(deps *config.HandlerDependencies) {
	if deps != nil {
		// Drop the test database with a fresh context
		ctx := context.Background()
		deps.DatabaseConfig.Client.Database(deps.DatabaseConfig.DatabaseName).Drop(ctx)
		// Close the connection
		deps.Close()
	}
}
