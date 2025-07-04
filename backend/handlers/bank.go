package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/services"
)

// BankHandler handles bank-related requests
type BankHandler struct {
	bankService *services.BankService
}

// NewBankHandler creates a new BankHandler
func NewBankHandler(container *config.Container) *BankHandler {
	return &BankHandler{
		bankService: container.ServiceContainer.Bank,
	}
}

// GetBank handles GET /api/bank
func (h *BankHandler) GetBank(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get username from the JWT (set by middleware)
	username := r.Header.Get("X-Username")
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	// Use the request context for proper cancellation handling
	ctx := r.Context()

	// Get bank by username
	bankResponse, err := h.bankService.GetBankByUsername(ctx, username)
	if err != nil {
		log.Printf("Error getting bank for username %s: %v", username, err)
		if err == services.ErrPlayerNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Player not found"})
			return
		}
		if err == services.ErrBankNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bank not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	json.NewEncoder(w).Encode(bankResponse)
}
