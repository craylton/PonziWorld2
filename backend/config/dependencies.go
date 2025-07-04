package config

import (
	"context"
	"ponziworld/backend/services"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// DatabaseConfig holds database connection and configuration
type DatabaseConfig struct {
	DatabaseName string
	Client       *mongo.Client
	// Note: We don't store the connection context here as it's meant for initial connection only
	// Individual handlers should create their own contexts as needed
	connectionCancel context.CancelFunc // Keep this for cleanup only
}

// HandlerDependencies holds all dependencies needed by handlers
type HandlerDependencies struct {
	ServiceManager *services.ServiceManager
	DatabaseConfig *DatabaseConfig
}

// NewHandlerDependencies creates a new HandlerDependencies instance
func NewHandlerDependencies(client *mongo.Client, ctx context.Context, cancel context.CancelFunc, databaseName string) *HandlerDependencies {
	dbConfig := &DatabaseConfig{
		DatabaseName:     databaseName,
		Client:           client,
		connectionCancel: cancel,
	}

	serviceManager := services.NewServiceManager(client.Database(databaseName))

	return &HandlerDependencies{
		ServiceManager: serviceManager,
		DatabaseConfig: dbConfig,
	}
}

// Close properly closes the database connection
func (d *HandlerDependencies) Close() {
	if d.DatabaseConfig.connectionCancel != nil {
		d.DatabaseConfig.connectionCancel()
	}
	if d.DatabaseConfig.Client != nil {
		// Create a context for disconnection
		ctx := context.Background()
		d.DatabaseConfig.Client.Disconnect(ctx)
	}
}
