package database

import (
	"context"
	"ponziworld/backend/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureAllIndexes(dbConfig *config.DatabaseConfig) error {
	if err := EnsurePlayerIndexes(dbConfig); err != nil {
		return err
	}

	if err := EnsureBankIndexes(dbConfig); err != nil {
		return err
	}

	if err := EnsureAssetIndexes(dbConfig); err != nil {
		return err
	}

	if err := EnsureAssetTypeIndexes(dbConfig); err != nil {
		return err
	}

	if err := EnsurePerformanceHistoryIndexes(dbConfig); err != nil {
		return err
	}

	if err := EnsureGameIndexes(dbConfig); err != nil {
		return err
	}

	if err := EnsurePendingTransactionIndexes(dbConfig); err != nil {
		return err
	}

	return nil
}

func ensureIndex(dbConfig *config.DatabaseConfig, collectionName, indexName string, model mongo.IndexModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := dbConfig.GetDatabase().Collection(collectionName)
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

func EnsurePlayerIndexes(dbConfig *config.DatabaseConfig) error {
	return ensureIndex(dbConfig, "players", "username_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("username_idx"),
	})
}

func EnsureBankIndexes(dbConfig *config.DatabaseConfig) error {
	return ensureIndex(dbConfig, "banks", "banks_playerId_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "playerId", Value: 1}},
		Options: options.Index().SetName("banks_playerId_idx"),
	})
}

func EnsureAssetIndexes(dbConfig *config.DatabaseConfig) error {
	return ensureIndex(dbConfig, "assets", "assets_bankId_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "bankId", Value: 1}},
		Options: options.Index().SetName("assets_bankId_idx"),
	})
}

func EnsureAssetTypeIndexes(dbConfig *config.DatabaseConfig) error {
	return ensureIndex(dbConfig, "assetTypes", "assetTypes_name_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("assetTypes_name_idx"),
	})
}

func EnsurePerformanceHistoryIndexes(dbConfig *config.DatabaseConfig) error {
	// Index for efficient queries by bankId, isClaimed, and day
	err := ensureIndex(dbConfig, "historicalPerformance", "performance_bankId_isClaimed_day_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "bankId", Value: 1}, {Key: "isClaimed", Value: 1}, {Key: "day", Value: 1}},
		Options: options.Index().SetName("performance_bankId_isClaimed_day_idx"),
	})
	if err != nil {
		return err
	}

	// Unique index to prevent duplicate entries for the same bank, day, and claimed status
	return ensureIndex(dbConfig, "historicalPerformance", "performance_unique_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "bankId", Value: 1}, {Key: "day", Value: 1}, {Key: "isClaimed", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("performance_unique_idx"),
	})
}

func EnsureGameIndexes(dbConfig *config.DatabaseConfig) error {
	// For the game collection, we don't need specific indexes as there's only one document
	// But we can ensure the collection exists by creating a basic index on _id
	return ensureIndex(dbConfig, "game", "game_id_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "_id", Value: 1}},
		Options: options.Index().SetName("game_id_idx"),
	})
}

func EnsurePendingTransactionIndexes(dbConfig *config.DatabaseConfig) error {
	// Index on buyerBankId for efficient lookups by buyer
	if err := ensureIndex(dbConfig, "pendingTransactions", "pending_buyerBankId_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "buyerBankId", Value: 1}},
		Options: options.Index().SetName("pending_buyerBankId_idx"),
	}); err != nil {
		return err
	}

	// Compound index on buyerBankId and assetId for efficient duplicate checking
	return ensureIndex(dbConfig, "pendingTransactions", "pending_buyer_asset_idx", mongo.IndexModel{
		Keys:    bson.D{{Key: "buyerBankId", Value: 1}, {Key: "assetId", Value: 1}},
		Options: options.Index().SetName("pending_buyer_asset_idx"),
	})
}
