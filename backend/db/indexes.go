package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureAllIndexes() error {
	client, ctx, cancel := ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)
	
	if err := EnsurePlayerIndexes(client); err != nil {
		return err
	}
	
	if err := EnsureBankIndexes(client); err != nil {
		return err
	}
	
	if err := EnsureAssetIndexes(client); err != nil {
		return err
	}
	
	if err := EnsurePerformanceHistoryIndexes(client); err != nil {
		return err
	}
	
	if err := EnsureGameIndexes(client); err != nil {
		return err
	}
	
	return nil
}

func ensureIndex(client *mongo.Client, collectionName, indexName string, model mongo.IndexModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := client.Database("ponziworld").Collection(collectionName)
	view := coll.Indexes()
	specs, err := view.ListSpecifications(ctx)
	if err != nil {
		return err
	}
	// Skip if already present
	for _, spec := range specs {
		if spec.Name == indexName {
			return nil
		}
	}
	// Create the index
	_, err = view.CreateOne(ctx, model)
	return err
}

func EnsurePlayerIndexes(client *mongo.Client) error {
	return ensureIndex(client, "players", "username_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("username_idx"),
	})
}

func EnsureBankIndexes(client *mongo.Client) error {
	return ensureIndex(client, "banks", "banks_playerId_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "playerId", Value: 1}},
		Options: options.Index().SetName("banks_playerId_idx"),
	})
}

func EnsureAssetIndexes(client *mongo.Client) error {
	return ensureIndex(client, "assets", "assets_bankId_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "bankId", Value: 1}},
		Options: options.Index().SetName("assets_bankId_idx"),
	})
}

func EnsurePerformanceHistoryIndexes(client *mongo.Client) error {
	// Index for efficient queries by bankId, isClaimed, and day
	err := ensureIndex(client, "historicalPerformance", "performance_bankId_isClaimed_day_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "bankId", Value: 1}, {Key: "isClaimed", Value: 1}, {Key: "day", Value: 1}},
		Options: options.Index().SetName("performance_bankId_isClaimed_day_idx"),
	})
	if err != nil {
		return err
	}

	// Unique index to prevent duplicate entries for the same bank, day, and claimed status
	return ensureIndex(client, "historicalPerformance", "performance_unique_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "bankId", Value: 1}, {Key: "day", Value: 1}, {Key: "isClaimed", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("performance_unique_idx"),
	})
}

func EnsureGameIndexes(client *mongo.Client) error {
	// For the game collection, we don't need specific indexes as there's only one document
	// But we can ensure the collection exists by creating a basic index on _id
	return ensureIndex(client, "game", "game_id_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "_id", Value: 1}},
		Options: options.Index().SetName("game_id_idx"),
	})
}
