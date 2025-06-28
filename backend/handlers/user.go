package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ponziworld/backend/db"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CreateNewPlayerHandler handles POST /api/newPlayer
func CreateNewPlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		BankName string `json:"bankName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	// Trim whitespace and validate
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.BankName = strings.TrimSpace(req.BankName)
	
	if req.Username == "" || req.Password == "" || req.BankName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username, password, and bank name required"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to hash password"})
		return
	}
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)
	
	// Create the player first
	playersCollection := client.Database("ponziworld").Collection("players")
	player := models.Player{
		Id:       primitive.NewObjectID(),
		Username: req.Username,
		Password: string(hashedPassword),
	}
	_, err = playersCollection.InsertOne(ctx, player)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Username already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create player"})
		return
	}

	// Create the bank for this player
	banksCollection := client.Database("ponziworld").Collection("banks")
	bank := models.Bank{
		Id:             primitive.NewObjectID(),
		PlayerId:         player.Id,
		BankName:       req.BankName,
		ClaimedCapital: 1000,
	}
	_, err = banksCollection.InsertOne(ctx, bank)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create bank"})
		return
	}

	// Create initial cash asset
	assetsCollection := client.Database("ponziworld").Collection("assets")
	asset := models.Asset{
		ID:        primitive.NewObjectID(),
		BankID:    bank.Id,
		Amount:    1000,
		AssetType: "Cash",
	}
	_, err = assetsCollection.InsertOne(ctx, asset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create initial asset"})
		return
	}

	// Create initial performance history (30 days of dummy claimed data + today's actual data)
	err = CreateInitialPerformanceHistory(client, ctx, bank.Id, 0, asset.Amount) // Using day 0 as current day
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create performance history"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
