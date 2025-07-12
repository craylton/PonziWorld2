package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"ponziworld/backend/config"
	"ponziworld/backend/models"
	"ponziworld/backend/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PendingTransactionHandler struct {
	pendingTransactionService *services.PendingTransactionService
	bankService               *services.BankService
}

func NewPendingTransactionHandler(container *config.Container) *PendingTransactionHandler {
	return &PendingTransactionHandler{
		pendingTransactionService: container.ServiceContainer.PendingTransaction,
		bankService:               container.ServiceContainer.Bank,
	}
}

func (h *PendingTransactionHandler) BuyAsset(w http.ResponseWriter, r *http.Request) {
	h.handleTransaction(w, r)
}

func (h *PendingTransactionHandler) SellAsset(w http.ResponseWriter, r *http.Request) {
	h.handleTransaction(w, r)
}

func (h *PendingTransactionHandler) handleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get username from JWT (set by middleware) - used for authorization only
	username := r.Header.Get("X-Username")
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	ctx := r.Context()

	// Parse the request body
	var req models.PendingTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Convert buyerBankId string to ObjectID
	buyerBankObjectID, err := primitive.ObjectIDFromHex(req.BuyerBankId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid buyer bank ID format"})
		return
	}

	// Convert assetId string to ObjectID
	assetObjectID, err := primitive.ObjectIDFromHex(req.AssetId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid asset ID format"})
		return
	}

	// Create the pending transaction
	err = h.pendingTransactionService.CreateTransaction(ctx, buyerBankObjectID, assetObjectID, req.Amount, username)
	if err != nil {
		log.Printf("Error creating transaction: %v", err)
		
		switch err {
		case services.ErrInvalidAmount:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Amount must not be zero"})
		case services.ErrAssetNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Asset not found"})
		case services.ErrInvalidBankID:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank ID"})
		case services.ErrSelfInvestment:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bank cannot invest in itself"})
		case services.ErrUnauthorizedBank:
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You do not own this bank"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create transaction"})
		}
		return
	}

	// Determine transaction type for response message
	transactionType := "buy"
	if req.Amount < 0 {
		transactionType = "sell"
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": transactionType + " transaction created successfully",
	})
}
