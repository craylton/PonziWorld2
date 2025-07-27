package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"ponziworld/backend/config"
	"ponziworld/backend/models"
	"ponziworld/backend/requestcontext"
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

// Transaction type for differentiating buy vs sell operations
type transactionType string

const (
	transactionTypeBuy  transactionType = "buy"
	transactionTypeSell transactionType = "sell"
)

func (h *PendingTransactionHandler) BuyAsset(w http.ResponseWriter, r *http.Request) {
	h.handleTransaction(w, r, transactionTypeBuy)
}

func (h *PendingTransactionHandler) SellAsset(w http.ResponseWriter, r *http.Request) {
	h.handleTransaction(w, r, transactionTypeSell)
}

func (h *PendingTransactionHandler) GetPendingTransactions(w http.ResponseWriter, r *http.Request) {
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

	// Extract bank ID from URL path parameter
	bankIdStr := r.PathValue("bankId")
	if bankIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bank ID required"})
		return
	}

	bankId, err := primitive.ObjectIDFromHex(bankIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank ID format"})
		return
	}

	// Get pending transactions for the bank
	transactions, err := h.pendingTransactionService.GetTransactionsByBuyerBankID(ctx, bankId, username)
	if err != nil {
		log.Printf("Error fetching pending transactions: %v", err)

		switch err {
		case services.ErrInvalidBankID:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bank not found"})
		case services.ErrPlayerNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Player not found"})
		case services.ErrUnauthorizedBank:
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You do not own this bank"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch pending transactions"})
		}
		return
	}

	if transactions == nil {
		transactions = []models.PendingTransactionResponse{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

func (h *PendingTransactionHandler) handleTransaction(
	w http.ResponseWriter,
	r *http.Request,
	transactionType transactionType,
) {
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

	// Parse the request body
	var req models.PendingTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Convert sourceBankId string to ObjectID
	sourceBankObjectID, err := primitive.ObjectIDFromHex(req.SourceBankId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid source bank ID format"})
		return
	}

	// Convert assetId string to ObjectID
	targetAssetObjectID, err := primitive.ObjectIDFromHex(req.TargetAssetId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid target asset ID format"})
		return
	}

	if transactionType == transactionTypeBuy {
		err = h.pendingTransactionService.CreateBuyTransaction(
			ctx,
			sourceBankObjectID,
			targetAssetObjectID,
			req.Amount,
			username,
		)
		if err != nil {
			log.Printf("Error creating buy transaction: %v", err)
			h.handleTransactionError(w, err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "buy transaction created successfully",
		})
	}
	if transactionType == transactionTypeSell {
		err = h.pendingTransactionService.CreateSellTransaction(
			ctx,
			sourceBankObjectID,
			targetAssetObjectID,
			req.Amount,
			username,
		)
		if err != nil {
			log.Printf("Error creating sell transaction: %v", err)
			h.handleTransactionError(w, err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "sell transaction created successfully",
		})
	}
}

func (h *PendingTransactionHandler) handleTransactionError(w http.ResponseWriter, err error) {
	switch err {
	case services.ErrInvalidAmount:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Amount must be positive"})
	case services.ErrTargetAssetNotFound:
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
	case services.ErrCashNotTradable:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cash cannot be bought or sold"})
	case services.ErrInsufficientFunds:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Insufficient cash balance for this purchase"})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create transaction"})
	}
}
