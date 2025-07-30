package database

import (
	"context"
	"log"
	"os"
	"ponziworld/backend/config"
	"ponziworld/backend/services"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitializeDatabaseConnection() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	// Use a timeout context for initial connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}
	return client, nil
}

func EnsureDatabaseStructure(dbConfig *config.DatabaseConfig, assetTypeService *services.AssetTypeService) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := EnsureAllIndexes(dbConfig); err != nil {
		return err
	}

	if err := assetTypeService.EnsureAssetTypesExist(ctx); err != nil {
		return err
	}

	return nil
}