package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/requestcontext"
	"ponziworld/backend/services"
)

type BankHandler struct {
	bankService *services.BankService
}

func NewBankHandler(container *config.Container) *BankHandler {
	return &BankHandler{
		bankService: container.ServiceContainer.Bank,
	}
}

// GetBanks handles GET /api/banks
func (h *BankHandler) GetBanks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get username from context (set by JwtMiddleware)
	ctx := r.Context()
	username, ok := requestcontext.UsernameFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	// Get all banks by username
	bankResponses, err := h.bankService.GetAllBanksByUsername(ctx, username)
	if err != nil {
		log.Printf("Error getting banks for username %s: %v", username, err)
		if err == services.ErrPlayerNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Player not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	json.NewEncoder(w).Encode(bankResponses)
}

// HandleBanks handles both GET and POST for /api/banks
func (h *BankHandler) HandleBanks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetBanks(w, r)
	case http.MethodPost:
		h.CreateBanks(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// CreateBanks handles POST /api/banks
func (h *BankHandler) CreateBanks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get username from context (set by JwtMiddleware)
	ctx := r.Context()
	username, ok := requestcontext.UsernameFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	// Parse request body
	var request struct {
		BankName       string `json:"bankName"`
		ClaimedCapital int64  `json:"claimedCapital"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate required fields
	if request.BankName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bank name is required"})
		return
	}

	if request.ClaimedCapital <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Claimed capital must be greater than 0"})
		return
	}

	// Create bank
	bank, err := h.bankService.CreateBankForUsername(ctx, username, request.BankName, request.ClaimedCapital)
	if err != nil {
		log.Printf("Error creating bank for username %s: %v", username, err)
		if err == services.ErrPlayerNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Player not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Return the created bank
	response := map[string]interface{}{
		"id":             bank.Id.Hex(),
		"bankName":       bank.BankName,
		"claimedCapital": bank.ClaimedCapital,
		"message":        "Bank created successfully",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
