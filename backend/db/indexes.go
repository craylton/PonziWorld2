package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ensureIndex installs a named index on the given collection, skipping creation if it already exists.
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

func EnsureUserIndexes(client *mongo.Client) error {
	return ensureIndex(client, "users", "username_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("username_idx"),
	})
}

func EnsureBankIndexes(client *mongo.Client) error {
	return ensureIndex(client, "banks", "banks_userId_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}},
		Options: options.Index().SetName("banks_userId_idx"),
	})
}

func EnsureAssetIndexes(client *mongo.Client) error {
	return ensureIndex(client, "assets", "assets_bankId_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "bankId", Value: 1}},
		Options: options.Index().SetName("assets_bankId_idx"),
	})
}
