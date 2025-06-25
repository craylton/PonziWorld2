package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureUserIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("ponziworld").Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func EnsureBankIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("ponziworld").Collection("banks")
	// Index on userId for fast lookup of assets by user
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func EnsureAssetIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("ponziworld").Collection("assets")
	
	// Index on bankId for fast lookup of assets by bank
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "bankId", Value: 1}},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
