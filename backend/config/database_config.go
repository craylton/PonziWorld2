package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DatabaseConfig struct {
	DatabaseName string
	Client       *mongo.Client
}

func (d *DatabaseConfig) GetDatabase() *mongo.Database {
	if d.Client == nil {
		return nil
	}
	return d.Client.Database(d.DatabaseName)
}

// Close disconnects the MongoDB client
func (d *DatabaseConfig) Close() {
	if d.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		d.Client.Disconnect(ctx)
	}
}