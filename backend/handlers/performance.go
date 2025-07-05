package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PerformanceHistoryHandler handles performance history-related requests
type PerformanceHistoryHandler struct {
	performanceHistoryService *services.PerformanceService
}

// NewPerformanceHistoryHandler creates a new PerformanceHistoryHandler
func NewPerformanceHistoryHandler(deps *config.Container) *PerformanceHistoryHandler {
	return &PerformanceHistoryHandler{
		performanceHistoryService: deps.ServiceContainer.Performance,
	}
}

// GetPerformanceHistoryHandler handles GET /api/performanceHistory/ownbank/{bankId}
func (h *PerformanceHistoryHandler) GetPerformanceHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Get username from JWT
	username := r.Header.Get("X-Username")
	if username == "" {
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank ID"})
		return
	}

	// Use the request context for proper cancellation handling
	ctx := r.Context()

	// Get performance history
	response, err := h.performanceHistoryService.GetPerformanceHistory(ctx, username, bankId)
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
		return
	}

	json.NewEncoder(w).Encode(response)
}
