package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ponziworld/backend/db"
	"ponziworld/backend/services"
)

// CreateNewPlayerHandler handles POST /api/newPlayer
func CreateNewPlayerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {		
		w.Header().Set("Allow", http.MethodPost)
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

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

	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	// Create service manager
	serviceManager := services.NewServiceManager(client.Database("ponziworld"))

	// Create new player with bank and initial assets
	err := serviceManager.Player.CreateNewPlayer(ctx, req.Username, req.Password, req.BankName)
	if err != nil {
		if err == services.ErrUsernameExists {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Username already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create player"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
