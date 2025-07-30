package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/requestcontext"
	"ponziworld/backend/services"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvestmentHandler struct {
	investmentService *services.InvestmentService
	logger            zerolog.Logger
}

func NewInvestmentHandler(container *config.Container) *InvestmentHandler {
	return &InvestmentHandler{
		investmentService: container.ServiceContainer.Investment,
		logger:            container.Logger,
	}
}

// GetInvestmentDetails handles GET /api/investment/{targetAssetId}/{sourceBankId}
func (h *InvestmentHandler) GetInvestmentDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		h.logger.Error().Msg("Invalid method for GetInvestmentDetails")
		return
	}

	// Get username from context
	username, ok := requestcontext.UsernameFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		h.logger.Error().Msg("Username not found in context")
		return
	}

	// Get path parameters
	targetAssetIdStr := r.PathValue("targetAssetId")
	sourceBankIdStr := r.PathValue("sourceBankId")

	// Validate and convert targetAssetId
	targetAssetId, err := primitive.ObjectIDFromHex(targetAssetIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid asset ID format"})
		h.logger.Error().Err(err).Msg("Failed to convert targetAssetId")
		return
	}

	// Validate and convert sourceBankId
	sourceBankId, err := primitive.ObjectIDFromHex(sourceBankIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank ID format"})
		h.logger.Error().Err(err).Msg("Failed to convert sourceBankId")
		return
	}

	// Get investment details
	response, err := h.investmentService.GetInvestmentDetails(r.Context(), username, targetAssetId, sourceBankId)
	if err != nil {
		switch err {
		case services.ErrUnauthorized:
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized access"})
		case services.ErrTargetAssetNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Target asset not found"})
		case services.ErrBankNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bank not found"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		}
		h.logger.Error().Err(err).Msg("Failed to get investment details")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
