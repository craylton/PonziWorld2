package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/requestcontext"
	"ponziworld/backend/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvestmentHandler struct {
	investmentService *services.InvestmentService
}

func NewInvestmentHandler(container *config.Container) *InvestmentHandler {
	return &InvestmentHandler{
		investmentService: container.ServiceContainer.Investment,
	}
}

// GetInvestmentDetails handles GET /api/investment/{targetAssetId}/{sourceBankId}
func (h *InvestmentHandler) GetInvestmentDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get username from context
	username, ok := requestcontext.UsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "username not found in context", http.StatusUnauthorized)
		return
	}

	// Get path parameters
	targetAssetIdStr := r.PathValue("targetAssetId")
	sourceBankIdStr := r.PathValue("sourceBankId")

	// Validate and convert targetAssetId
	targetAssetId, err := primitive.ObjectIDFromHex(targetAssetIdStr)
	if err != nil {
		http.Error(w, "invalid asset ID", http.StatusBadRequest)
		return
	}

	// Validate and convert sourceBankId
	sourceBankId, err := primitive.ObjectIDFromHex(sourceBankIdStr)
	if err != nil {
		http.Error(w, "invalid bank ID", http.StatusBadRequest)
		return
	}

	// Get investment details
	response, err := h.investmentService.GetInvestmentDetails(r.Context(), username, targetAssetId, sourceBankId)
	if err != nil {
		switch err {
		case services.ErrUnauthorized:
			http.Error(w, "unauthorized access", http.StatusUnauthorized)
		case services.ErrTargetAssetNotFound:
			http.Error(w, "target asset not found", http.StatusNotFound)
		case services.ErrBankNotFound:
			http.Error(w, "bank not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
