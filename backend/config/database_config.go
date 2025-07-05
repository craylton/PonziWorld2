package config

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DatabaseConfig struct {
	DatabaseName string
	Client       *mongo.Client
	// Note: We don't store the connection context here as it's meant for initial connection only
	// Individual handlers should create their own contexts as needed
	connectionCancel context.CancelFunc // Keep this for cleanup only
}

func (d *DatabaseConfig) GetDatabase() *mongo.Database {
	if d.Client == nil {
		return nil
	}
	return d.Client.Database(d.DatabaseName)
}

func (d *DatabaseConfig) Close() {
	if d.connectionCancel != nil {
		d.connectionCancel()
	}
	if d.Client != nil {
		ctx := context.Background()
		d.Client.Disconnect(ctx)
	}
}