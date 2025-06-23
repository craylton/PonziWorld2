package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	// Set up timeout context for connection and ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		cancel()
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	return client, ctx, cancel
}
