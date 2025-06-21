package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/db"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// LoginHandler handles POST /api/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	if req.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username required"})
		return
	}
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database("ponziworld").Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	json.NewEncoder(w).Encode(user)
}
