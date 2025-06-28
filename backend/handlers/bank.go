package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/db"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// GetBankHandler handles GET /api/bank
func GetBankHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get username from the JWT token (set by middleware)
	username := r.Header.Get("X-Username")
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	// First, find the player to get their ID
	playersCollection := client.Database("ponziworld").Collection("players")
	var player models.Player
	err := playersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&player)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Player not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Find the bank for this player
	banksCollection := client.Database("ponziworld").Collection("banks")
	var bank models.Bank
	err = banksCollection.FindOne(ctx, bson.M{"playerId": player.Id}).Decode(&bank)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bank not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Find all assets for this bank
	assetsCollection := client.Database("ponziworld").Collection("assets")
	cursor, err := assetsCollection.Find(ctx, bson.M{"bankId": bank.Id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer cursor.Close(ctx)

	var assets []models.Asset
	if err = cursor.All(ctx, &assets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Calculate actual capital by summing all asset amounts
	var actualCapital int64 = 0
	for _, asset := range assets {
		actualCapital += asset.Amount
	}

	// Create response
	response := models.BankResponse{
		ID:             bank.Id.Hex(),
		BankName:       bank.BankName,
		ClaimedCapital: bank.ClaimedCapital,
		ActualCapital:  actualCapital,
		Assets:         assets,
	}

	json.NewEncoder(w).Encode(response)
}
