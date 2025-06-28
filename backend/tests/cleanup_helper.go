package tests

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"ponziworld/backend/db"
	"ponziworld/backend/models"
)

// CleanupTestData removes test data from the database
// This function should be called in t.Cleanup() for each test
func CleanupTestData(username, bankName string) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	// Delete player
	playersCollection := client.Database("ponziworld").Collection("players")
	playersCollection.DeleteOne(ctx, bson.M{"username": username})

	// Delete bank and associated assets
	cleanupBankAndAssets(ctx, client, bankName)
}

// CleanupMultipleTestData removes multiple test players and their data
func CleanupMultipleTestData(playersAndBanks map[string]string) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	playersCollection := client.Database("ponziworld").Collection("players")

	for username, bankName := range playersAndBanks {
		// Delete player
		playersCollection.DeleteOne(ctx, bson.M{"username": username})
		
		// Delete bank and associated assets
		cleanupBankAndAssets(ctx, client, bankName)
	}
}

// cleanupBankAndAssets is a helper function to clean up banks and their associated assets
func cleanupBankAndAssets(ctx context.Context, client *mongo.Client, bankName string) {
	banksCollection := client.Database("ponziworld").Collection("banks")
	assetsCollection := client.Database("ponziworld").Collection("assets")
	historyCollection := client.Database("ponziworld").Collection("historicalPerformance")

	// Find all banks with the given name and delete their assets and performance history
	cursor, err := banksCollection.Find(ctx, bson.M{"bankName": bankName})
	if err == nil {
		for cursor.Next(ctx) {
			var bank models.Bank
			cursor.Decode(&bank)
			// Delete associated assets
			assetsCollection.DeleteMany(ctx, bson.M{"bankId": bank.Id})
			// Delete associated performance history
			historyCollection.DeleteMany(ctx, bson.M{"bankId": bank.Id})
		}
		cursor.Close(ctx)
	}

	// Delete the banks
	banksCollection.DeleteMany(ctx, bson.M{"bankName": bankName})
}
