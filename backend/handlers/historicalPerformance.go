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

type HistoricalPerformanceHandler struct {
	historicalPerformanceService *services.HistoricalPerformanceService
	logger                       zerolog.Logger
}

func NewHistoricalPerformanceHandler(container *config.Container) *HistoricalPerformanceHandler {
	return &HistoricalPerformanceHandler{
		historicalPerformanceService: container.ServiceContainer.HistoricalPerformance,
		logger:                       container.Logger,
	}
}

// GetHistoricalPerformanceHandler handles GET /api/historicalperformance/ownbank/{bankId}
func (h *HistoricalPerformanceHandler) GetHistoricalPerformance(w http.ResponseWriter, r *http.Request) {
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
		h.logger.Error().Msg("Username not found in context")
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank ID"})
		h.logger.Error().Err(err).Str("bankIdStr", bankIdStr).Msg("Invalid bank ID format")
		return
	}

	// Get performance history
	response, err := h.historicalPerformanceService.GetOwnBankHistoricalPerformance(ctx, username, bankId)
	if err != nil {
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
		if err == services.ErrUnauthorized {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: You can only view your own bank's performance history"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		h.logger.Error().Err(err).Str("username", username).Str("bankId", bankId.Hex()).Msg("Failed to get historical performance")
		return
	}

	json.NewEncoder(w).Encode(response)
}

// GetAssetHistoricalPerformanceHandler handles GET /api/historicalPerformance/asset/{targetAssetId}/{sourceBankId}
func (h *HistoricalPerformanceHandler) GetAssetHistoricalPerformance(w http.ResponseWriter, r *http.Request) {
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
		h.logger.Error().Msg("Username not found in context for GetAssetHistoricalPerformance")
		return
	}

	// Extract target asset ID from URL path parameter
	targetAssetIdStr := r.PathValue("targetAssetId")
	if targetAssetIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Target asset ID required"})
		return
	}

	targetAssetId, err := primitive.ObjectIDFromHex(targetAssetIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid target asset ID"})
		h.logger.Error().Err(err).Str("targetAssetIdStr", targetAssetIdStr).Msg("Invalid target asset ID format")
		return
	}

	// Extract source bank ID from URL path parameter
	sourceBankIdStr := r.PathValue("sourceBankId")
	if sourceBankIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Source bank ID required"})
		return
	}

	sourceBankId, err := primitive.ObjectIDFromHex(sourceBankIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid source bank ID"})
		h.logger.Error().Err(err).Str("sourceBankIdStr", sourceBankIdStr).Msg("Invalid source bank ID format")
		return
	}

	// Get asset performance history
	response, err := h.historicalPerformanceService.GetAssetHistoricalPerformance(
		ctx,
		username,
		targetAssetId,
		sourceBankId,
		30,
	)
	if err != nil {
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
		if err == services.ErrUnauthorized {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: You can only view performance history for your own banks"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		h.logger.Error().Err(err).Str("username", username).Str("targetAssetId", targetAssetId.Hex()).Str("sourceBankId", sourceBankId.Hex()).Msg("Failed to get asset historical performance")
		return
	}

	json.NewEncoder(w).Encode(response)
}
