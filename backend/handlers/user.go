package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ponziworld/backend/db"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserHandler handles POST /api/user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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
	collection := client.Database("ponziworld").Collection("users")
	user := models.User{
		ID:             primitive.NewObjectID(),
		Username:       req.Username,
		Password:       string(hashedPassword),
		BankName:       req.BankName,
		ClaimedCapital: 1000,
		ActualCapital:  1000,
	}
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Username already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetUserHandler handles GET /api/user
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
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
	collection := client.Database("ponziworld").Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	json.NewEncoder(w).Encode(user)
}
